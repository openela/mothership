package srpm_import

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	storage2 "github.com/go-git/go-git/v5/storage"
	"github.com/openela/mothership/base/storage"
	"github.com/pkg/errors"
	srpmprocpb "github.com/rocky-linux/srpmproc/pb"
	"github.com/rocky-linux/srpmproc/pkg/data"
	"github.com/rocky-linux/srpmproc/pkg/directives"
	"github.com/sassoftware/go-rpmutils"
	"golang.org/x/crypto/openpgp"
	"google.golang.org/protobuf/encoding/prototext"
)

var (
	elDistRegex  = regexp.MustCompile(`el\d+`)
	releaseRegex = regexp.MustCompile(`.*release (\d+\.\d+).*`)
	branchRegex  = regexp.MustCompile(`^(el-\d+\.\d+)`)
)

type State struct {
	// tempDir is the temporary directory where the SRPM is extracted to.
	tempDir string

	// rpm is the SRPM.
	rpm *rpmutils.Rpm

	// authorName is the name of the author of the commit.
	authorName string

	// authorEmail is the email of the author of the commit.
	authorEmail string

	// lookasideBlobs is a map of blob names to their SHA256 hashes.
	lookasideBlobs map[string]string

	// rolling determines how the branch is named.
	// if true, the branch is named "elX" where X is the major release
	// if false, the branch is named "el-X.Y" where X.Y is the full release
	rolling bool

	// tag is the tag name
	tag string
}

type ImportOutput struct {
	// Commit is the commit object
	Commit *object.Commit

	// Branch is the branch name
	Branch string

	// Tag is the tag name
	Tag string
}

// copyFromOS copies specified file from OS filesystem to target filesystem.
func copyFromOS(targetFS billy.Filesystem, path string, targetPath string) error {
	// Open file from OS filesystem.
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to get file info")
	}

	// Create file in target filesystem.
	targetFile, err := targetFS.OpenFile(targetPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer targetFile.Close()

	// Copy contents of file from OS filesystem to target filesystem.
	_, err = io.Copy(targetFile, f)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}

	return nil
}

// FromFile creates a new State from an SRPM file.
// The SRPM file is extracted to a temporary directory.
func FromFile(path string, rolling bool, keys ...*openpgp.Entity) (*State, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	// If keys is not empty, then verify the RPM signature.
	if len(keys) > 0 {
		_, _, err := rpmutils.Verify(f, keys)
		if err != nil {
			return nil, errors.Wrap(err, "failed to verify RPM")
		}

		// After verifying the RPM, seek back to the start of the file.
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, errors.Wrap(err, "failed to seek to start of file")
		}
	}

	rpm, err := rpmutils.ReadRpm(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read RPM")
	}

	state := &State{
		rpm:            rpm,
		authorName:     "Mship Bot",
		authorEmail:    "no-reply+mshipbot@openela.org",
		lookasideBlobs: make(map[string]string),
		rolling:        rolling,
	}

	// Create a temporary directory.
	state.tempDir, err = os.MkdirTemp("", "srpm_import-*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary directory")
	}

	// Extract the SRPM.
	err = rpm.ExpandPayload(state.tempDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract SRPM")
	}

	return state, nil
}

func (s *State) Close() error {
	return os.RemoveAll(s.tempDir)
}

func (s *State) GetDir() string {
	return s.tempDir
}

func (s *State) SetAuthor(name, email string) {
	s.authorName = name
	s.authorEmail = email
}

// determineLookasideBlobs determines which blobs need to be uploaded to the
// lookaside cache.
// Currently, the rule is that if a file is larger than 5MB, and is binary,
// then it should be uploaded to the lookaside cache.
// If the file name contains ".tar", then it is assumed to be a tarball, and
// is ALWAYS uploaded to the lookaside cache.
func (s *State) determineLookasideBlobs() error {
	// Read all files in tempDir, except for the SPEC file
	// For each file, if it is larger than 5MB, and is binary, then add it to
	// the lookasideBlobs map.
	// If the file is not binary, then skip it.
	ls, err := os.ReadDir(s.tempDir)
	if err != nil {
		return errors.Wrap(err, "failed to read directory")
	}

	for _, f := range ls {
		// If file ends with ".spec", then skip it.
		if f.IsDir() || strings.HasSuffix(f.Name(), ".spec") {
			continue
		}

		// If file is larger than 5MB, then add it to the lookasideBlobs map.
		info, err := f.Info()
		if err != nil {
			return errors.Wrap(err, "failed to get file info")
		}

		if info.Size() > 5*1024*1024 || strings.Contains(f.Name(), ".tar") {
			sum, err := func() (string, error) {
				hash := sha256.New()
				file, err := os.Open(filepath.Join(s.tempDir, f.Name()))
				if err != nil {
					return "", errors.Wrap(err, "failed to open file")
				}
				defer file.Close()

				_, err = io.Copy(hash, file)
				if err != nil {
					return "", errors.Wrap(err, "failed to copy file")
				}

				return hex.EncodeToString(hash.Sum(nil)), nil
			}()
			if err != nil {
				return err
			}

			s.lookasideBlobs[f.Name()] = sum
		}
	}

	return nil
}

// uploadLookasideBlobs uploads all blobs in the lookasideBlobs map to the
// lookaside cache.
func (s *State) uploadLookasideBlobs(lookaside storage.Storage) error {
	// The object name is the SHA256 hash of the file.
	for path, hash := range s.lookasideBlobs {
		// First check if they exist, since it's a waste of time to upload
		// something that already exists.
		// They are uploaded by hash, so if the hash already exists, then the
		// file already exists.
		exists, err := lookaside.Exists(hash)
		if err != nil {
			return errors.Wrap(err, "failed to check if blob exists")
		}

		if exists {
			continue
		}

		_, err = lookaside.Put(hash, filepath.Join(s.tempDir, path))
		if err != nil {
			return errors.Wrap(err, "failed to upload file")
		}
	}

	return nil
}

// writeMetadata file writes the metadata map file.
// The metadata file contains lines of the format:
//
//	<path to download> <blob hash>
//
// For example:
//
//	1234567890abcdef SOURCES/bar
func (s *State) writeMetadataFile(targetFS billy.Filesystem) error {
	// Open metadata file for writing.
	name, err := s.rpm.Header.GetStrings(rpmutils.NAME)
	if err != nil {
		return errors.Wrap(err, "failed to get RPM name")
	}

	metadataFile := fmt.Sprintf(".%s.metadata", name[0])

	// Delete the file if it exists
	_ = targetFS.Remove(metadataFile)

	f, err := targetFS.Create(metadataFile)
	if err != nil {
		return errors.Wrap(err, "failed to open metadata file")
	}
	defer f.Close()

	// Write each line to the metadata file.
	for path, hash := range s.lookasideBlobs {
		// RPM sources MUST be in SOURCES/ directory
		_, err = f.Write([]byte(hash + " " + filepath.Join("SOURCES", path) + "\n"))
		if err != nil {
			return errors.Wrap(err, "failed to write line to metadata file")
		}
	}

	// Each file in metadata needs to be added to gitignore
	// Overwrite the gitignore file
	gitignoreFile := ".gitignore"
	f, err = targetFS.OpenFile(gitignoreFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open gitignore file")
	}

	// Write each line to the gitignore file.
	for path, _ := range s.lookasideBlobs {
		_, err = f.Write([]byte(filepath.Join("SOURCES", path) + "\n"))
		if err != nil {
			return errors.Wrap(err, "failed to write line to gitignore file")
		}
	}

	return nil
}

// ExpandLayout expands the layout of the SRPM into the target filesystem.
// Moves all sources into SOURCES/ directory.
// Spec file is moved to SPECS/ directory.
func (s *State) ExpandLayout(targetFS billy.Filesystem) error {
	// Create SOURCES/ directory.
	err := targetFS.MkdirAll("SOURCES", 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create SOURCES directory")
	}

	// Copy all files from OS filesystem to target filesystem.
	ls, err := os.ReadDir(s.tempDir)
	if err != nil {
		return errors.Wrap(err, "failed to read directory")
	}

	for _, f := range ls {
		baseName := filepath.Base(f.Name())
		// If file ends with ".spec", then copy to SPECS/ directory.
		if strings.HasSuffix(f.Name(), ".spec") {
			err := copyFromOS(targetFS, filepath.Join(s.tempDir, f.Name()), filepath.Join("SPECS", baseName))
			if err != nil {
				return errors.Wrap(err, "failed to copy spec file")
			}
		} else {
			// Copy all other files to SOURCES/ directory.
			// Only if they are not present in lookasideBlobs
			_, ok := s.lookasideBlobs[f.Name()]
			if ok {
				continue
			}
			err := copyFromOS(targetFS, filepath.Join(s.tempDir, f.Name()), filepath.Join("SOURCES", baseName))
			if err != nil {
				return errors.Wrap(err, "failed to copy file")
			}
		}
	}

	return nil
}

// getStreamSuffix adds a "-stream-X" suffix if the given RPM is a module component.
// This is determined using Modularitylabel (5096). If the label is present, then
// the RPM is a module component. Label format is MODULE_NAME:STREAM:VERSION:CONTEXT.
// This function returns an empty string if the RPM is not a module component.
func (s *State) getStreamSuffix() (string, error) {
	// Check the modularity label
	label, err := s.rpm.Header.GetString(5096)
	if err != nil {
		// If it's not present at all, it will fail with "No such entry 5096"
		return "", nil
	}

	// If the label is empty, then the RPM is not a module component
	if label == "" {
		return "", nil
	}

	// Split the label
	parts := strings.Split(label, ":")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid modularity label")
	}

	// Return the stream
	return fmt.Sprintf("-stream-%s", parts[1]), nil
}

// getRepo returns the target repository for the SRPM.
// This is where the payload is uploaded to.
func (s *State) getRepo(opts *git.CloneOptions, storer storage2.Storer, targetFS billy.Filesystem, osRelease string) (*git.Repository, string, error) {
	// Determine branch
	// If the OS release is not specified, then we use the dist tag
	var branch string
	if osRelease == "" {
		// Determine dist tag
		nevra, err := s.rpm.Header.GetNEVRA()
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to get NEVRA")
		}

		// The dist tag will be used as the branch
		dist := elDistRegex.FindString(nevra.Release)
		if dist == "" {
			return nil, "", errors.Wrap(err, "failed to determine dist tag")
		}

		if s.rolling {
			branch = dist
		} else {
			branch = "el-" + dist[2:]
		}
	} else {
		// Determine branch from OS release
		if !releaseRegex.MatchString(osRelease) {
			return nil, "", fmt.Errorf("invalid OS release %s", osRelease)
		}
		ver := releaseRegex.FindStringSubmatch(osRelease)[1]

		if s.rolling {
			dist := elDistRegex.FindString("el" + ver)
			if dist == "" {
				return nil, "", errors.New("failed to determine dist tag")
			}
			branch = dist
		} else {
			branch = "el-" + ver
		}
	}

	// Check if module component
	streamSuffix, err := s.getStreamSuffix()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get stream suffix")
	}
	branch += streamSuffix

	// Set branch to dist tag
	opts.ReferenceName = plumbing.NewBranchReferenceName(branch)
	opts.SingleBranch = true

	// Clone the repository, to the target filesystem.
	// We do an init, then a fetch, then a checkout
	// If the repo doesn't exist, then we init only
	repo, err := git.Init(storer, targetFS)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil, "", errors.Wrap(err, "failed to init repo")
		}
		repo, err = git.Open(storer, targetFS)
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to open repo")
		}
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get worktree")
	}

	// Create a new remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{opts.URL},
		Fetch: []config.RefSpec{
			"refs/heads/*:refs/heads/*",
		},
	})
	if err != nil {
		if !errors.Is(err, git.ErrRemoteExists) {
			return nil, "", errors.Wrap(err, "failed to create remote")
		}
	}

	// Fetch the remote
	err = repo.Fetch(&git.FetchOptions{
		Auth:       opts.Auth,
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			"refs/heads/*:refs/heads/*",
		},
	})
	if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil, "", errors.Wrap(err, "failed to fetch remote")
	}

	// Checkout the branch
	refName := plumbing.NewBranchReferenceName(branch)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err != nil {
		// Sort branches by descending and find the closest branch (usually one point release less)
		// and check if that branch exists. If it exists then copy "PATCHES" directory from that branch
		// to the current branch
		// This is so we can carry over patches from previous point releases
		// This only applies to non-rolling releases as rolling uses the same branch
		if !s.rolling && osRelease != "" {
			tempFs := memfs.New()

			var branchesToCheck []string
			branchSubmatch := branchRegex.FindStringSubmatch(branch)
			if len(branchSubmatch) == 2 {
				lastDigit, err := strconv.Atoi(branchSubmatch[1][len(branchSubmatch[1])-1:])
				if err != nil {
					return nil, "", errors.Wrap(err, "failed to convert branch number")
				}
				branchWithoutPrefix := strings.TrimPrefix(branch, branchSubmatch[0])
				elWithoutLastDigit := branchSubmatch[0][:len(branchSubmatch[0])-1]

				// From digit to 0, they are all potential branches
				for i := lastDigit; i >= 0; i-- {
					newPrefix := fmt.Sprintf("%s%d%s", elWithoutLastDigit, i, branchWithoutPrefix)
					branchesToCheck = append(branchesToCheck, newPrefix)
				}
			}

			// Check if any of the branches exist
			var closestBranch string
			for _, branchToCheck := range branchesToCheck {
				err = wt.Checkout(&git.CheckoutOptions{
					Branch: plumbing.NewBranchReferenceName(branchToCheck),
					Force:  true,
				})
				if err == nil {
					closestBranch = branchToCheck
					break
				}
			}

			// If a branch exists, then copy the PATCHES directory
			if closestBranch != "" {
				// Copy the PATCHES directory
				err = data.CopyFromFs(wt.Filesystem, tempFs, "PATCHES")
				if err != nil {
					return nil, "", errors.Wrap(err, "failed to copy PATCHES directory")
				}

				// Checkout the original branch
				h := plumbing.NewSymbolicReference(plumbing.HEAD, refName)
				if err := repo.Storer.CheckAndSetReference(h, nil); err != nil {
					return nil, "", errors.Wrap(err, "failed to checkout branch")
				}

				// Copy the PATCHES directory
				err = data.CopyFromFs(tempFs, wt.Filesystem, "PATCHES")
				if err != nil {
					return nil, "", errors.Wrap(err, "failed to copy PATCHES directory")
				}
			}

			h := plumbing.NewSymbolicReference(plumbing.HEAD, refName)
			if err := repo.Storer.CheckAndSetReference(h, nil); err != nil {
				return nil, "", errors.Wrap(err, "failed to checkout branch")
			}
		} else {
			h := plumbing.NewSymbolicReference(plumbing.HEAD, refName)
			if err := repo.Storer.CheckAndSetReference(h, nil); err != nil {
				return nil, "", errors.Wrap(err, "failed to checkout branch")
			}
		}
	}

	return repo, branch, nil
}

// cleanTargetRepo deletes all files in the target repository.
func (s *State) cleanTargetRepo(wt *git.Worktree, root string) error {
	// Delete all files in the target repository.
	ls, err := wt.Filesystem.ReadDir(root)
	if err != nil {
		return errors.Wrap(err, "failed to read directory")
	}

	for _, f := range ls {
		// Don't delete the PATCHES directory
		if f.Name() == "PATCHES" && f.IsDir() {
			continue
		}

		// If it's a directory, then recurse into it.
		if f.IsDir() {
			err := s.cleanTargetRepo(wt, filepath.Join(root, f.Name()))
			if err != nil {
				return errors.Wrap(err, "failed to clean target repo")
			}
		} else {
			// Otherwise, delete the file.
			_, err := wt.Remove(filepath.Join(root, f.Name()))
			if err != nil {
				return errors.Wrap(err, "failed to remove file")
			}
		}
	}

	return nil
}

// populateTargetRepo runs the following steps:
// 1. Clean the target repository.
// 2. Determine which blobs need to be uploaded to the lookaside cache.
// 3. Upload blobs to the lookaside cache.
// 4. Write the metadata file.
// 5. Expand the layout of the SRPM.
// 6. Commit the changes to the target repository.
func (s *State) populateTargetRepo(repo *git.Repository, targetFS billy.Filesystem, lookaside storage.Storage, branch string) error {
	// Clean the target repository.
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	err = s.cleanTargetRepo(wt, ".")
	if err != nil {
		return errors.Wrap(err, "failed to clean target repo")
	}

	// Determine which blobs need to be uploaded to the lookaside cache.
	err = s.determineLookasideBlobs()
	if err != nil {
		return errors.Wrap(err, "failed to determine lookaside blobs")
	}

	// Upload blobs to the lookaside cache.
	err = s.uploadLookasideBlobs(lookaside)
	if err != nil {
		return errors.Wrap(err, "failed to upload lookaside blobs")
	}

	// Write the metadata file.
	err = s.writeMetadataFile(targetFS)
	if err != nil {
		return errors.Wrap(err, "failed to write metadata file")
	}

	// Expand the layout of the SRPM.
	err = s.ExpandLayout(targetFS)
	if err != nil {
		return errors.Wrap(err, "failed to expand layout")
	}

	// If the target FS has patches, apply the directives
	err = s.patchTargetRepo(repo, lookaside)
	if err != nil {
		return errors.Wrap(err, "failed to patch target repo")
	}

	// Commit the changes to the target repository.
	_, err = wt.Add(".")
	if err != nil {
		return errors.Wrap(err, "failed to add files")
	}

	nevra, err := s.rpm.Header.GetNEVRA()
	if err != nil {
		return errors.Wrap(err, "failed to get NEVRA")
	}
	importStr := fmt.Sprintf("import %s-%s-%s", nevra.Name, nevra.Version, nevra.Release)
	hash, err := wt.Commit(importStr, &git.CommitOptions{
		Author: &object.Signature{
			Name:  s.authorName,
			Email: s.authorEmail,
			When:  time.Now(),
		},
		AllowEmptyCommits: true,
	})
	if err != nil {
		return errors.Wrap(err, "failed to commit changes")
	}

	// Create a tag
	// The tag should follow the following format:
	//   imports/<branch>/<nvra>
	tag := fmt.Sprintf("imports/%s/%s-%s-%s", branch, nevra.Name, nevra.Version, nevra.Release)
	// Escape ^ in tag
	tag = strings.ReplaceAll(tag, "^", "_")
	// Escape ~ in tag
	tag = strings.ReplaceAll(tag, "~", "_")
	// Replace % with _
	tag = strings.ReplaceAll(tag, "%", "_")

	s.tag = tag

	_, err = repo.CreateTag(tag, hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  s.authorName,
			Email: s.authorEmail,
			When:  time.Now(),
		},
		Message: tag,
	})
	if err != nil {
		if errors.Is(err, git.ErrTagExists) {
			err = repo.DeleteTag(tag)
			if err != nil {
				return errors.Wrap(err, "failed to delete tag")
			}
			_, err = repo.CreateTag(tag, hash, &git.CreateTagOptions{
				Tagger: &object.Signature{
					Name:  s.authorName,
					Email: s.authorEmail,
					When:  time.Now(),
				},
				Message: tag,
			})
		}
		if err != nil {
			return errors.Wrap(err, "failed to create tag")
		}
	}

	return nil
}

// pushTargetRepo pushes the target repository to the upstream repository.
func (s *State) pushTargetRepo(repo *git.Repository, opts *git.PushOptions) error {
	// Push the target repository to the upstream repository.
	err := repo.Push(opts)
	if err != nil {
		return errors.Wrap(err, "failed to push repo")
	}

	return nil
}

func (s *State) patchTargetRepo(repo *git.Repository, lookaside storage.Storage) error {
	// We can re-use srpmproc as we should stay compatible with it
	// Instead of OpenPatch, we'll look for patches in the targetFS
	// todo(mustafa): RESF still uses OpenPatch, so we'll need to change that
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	nevra, err := s.rpm.Header.GetNEVRA()
	if err != nil {
		return errors.Wrap(err, "failed to get NEVRA")
	}

	dist := elDistRegex.FindString(nevra.Release)
	if dist == "" {
		return errors.Wrap(err, "failed to determine dist tag")
	}
	distNum, err := strconv.Atoi(dist[2:])
	if err != nil {
		return errors.Wrap(err, "failed to parse dist tag")
	}

	pd := &data.ProcessData{
		ImportBranchPrefix: "el",
		Version:            distNum,
		BlobStorage:        &srpmprocBlobCompat{lookaside},
		Importer:           &srpmprocImportModeCompat{},
		Log:                log.New(os.Stderr, "", 0),
	}
	md := &data.ModeData{
		SourcesToIgnore: []*data.IgnoredSource{},
	}

	// Look in the PATCHES/ directory for any .cfg files
	patchesLs, err := wt.Filesystem.ReadDir("PATCHES")
	if err != nil {
		return errors.Wrap(err, "failed to read PATCHES directory")
	}

	for _, f := range patchesLs {
		// Skip directories
		if f.IsDir() {
			continue
		}

		// Skip non-cfg files
		if !strings.HasSuffix(f.Name(), ".cfg") {
			continue
		}

		// Open the file
		file, err := wt.Filesystem.Open(filepath.Join("PATCHES", f.Name()))
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}

		// Process the file
		directivesBytes, err := io.ReadAll(file)
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}

		var cfg srpmprocpb.Cfg
		err = prototext.Unmarshal(directivesBytes, &cfg)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal directives")
		}

		errs := directives.Apply(&cfg, pd, md, wt, wt)
		// If there are errors, then we should return a reduced error
		if len(errs) > 0 {
			retErr := errors.New("failed to apply directives")
			for _, err := range errs {
				retErr = errors.Wrap(retErr, err.Error())
			}
			return retErr
		}
	}

	// Add sources to ignore to lookasideBlobs
	for _, source := range md.SourcesToIgnore {
		// Get the hash of the source
		hash, err := func() (string, error) {
			hash := sha256.New()
			file, err := wt.Filesystem.Open(source.Name)
			if err != nil {
				return "", errors.Wrap(err, "failed to open file")
			}
			defer file.Close()

			_, err = io.Copy(hash, file)
			if err != nil {
				return "", errors.Wrap(err, "failed to copy file")
			}

			return hex.EncodeToString(hash.Sum(nil)), nil
		}()
		if err != nil {
			return err
		}

		s.lookasideBlobs[source.Name] = hash
	}

	// Re-write the metadata file
	err = s.writeMetadataFile(wt.Filesystem)
	if err != nil {
		return errors.Wrap(err, "failed to write metadata file")
	}

	return nil
}

// Import imports the SRPM into the target repository.
func (s *State) Import(opts *git.CloneOptions, storer storage2.Storer, targetFS billy.Filesystem, lookaside storage.Storage, osRelease string) (*ImportOutput, error) {
	// Get the target repository.
	repo, branch, err := s.getRepo(opts, storer, targetFS, osRelease)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repo")
	}

	// Populate the target repository.
	err = s.populateTargetRepo(repo, targetFS, lookaside, branch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate target repo")
	}

	// Push the target repository.
	err = s.pushTargetRepo(repo, &git.PushOptions{
		Force: true,
		Auth:  opts.Auth,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%[1]s", branch)),
			config.RefSpec(fmt.Sprintf("refs/tags/imports/%s/*:refs/tags/imports/%[1]s/*", branch)),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to push target repo")
	}

	// Get latest commit
	head, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get HEAD")
	}

	// Get commit object
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get commit object")
	}

	return &ImportOutput{
		Commit: commit,
		Branch: branch,
		Tag:    s.tag,
	}, nil
}
