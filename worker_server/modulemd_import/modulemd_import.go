package modulemd_import

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage"
	"github.com/pkg/errors"
	"github.com/rocky-linux/srpmproc/modulemd"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	elDistRegex  = regexp.MustCompile(`el\d+`)
	releaseRegex = regexp.MustCompile(`.*release (\d+\.\d+).*`)
	gitlabURI    = "https://gitlab.com/redhat/centos-stream"
)

type State struct {
	// moduleName is the name of the module
	moduleName string

	// authorName is the name of the author of the commit
	authorName string

	// authorEmail is the email of the author of the commit
	authorEmail string

	// rolling determines how the branch is named.
	// if true, the branch is named "elX" where X is the major release
	// if false, the branch is named "el-X.Y" where X.Y is the full release
	rolling bool
}

type ImportOutput struct {
	// Commit is the commit object
	Commit *object.Commit

	// Branch is the branch name
	Branch string

	// Tag is the tag name
	Tag string
}

// getURI returns the URI for the module source
func (s *State) getURI(pkg string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(gitlabURI, "/"), pkg)
}

// getUpstreamRepo returns the source repo from CentOS Stream gitlab
// It finds the branch of format "stream-MODULE_NAME-STREAM-rhel-X.Y.0" where X.Y is <= osRelease X.Y
// Let's find ALL branches for MODULE_NAME and X.Y, the STREAM can be anything (as in multiple streams for one module)
func (s *State) getUpstreamRepo(opts *git.CloneOptions, storer storage.Storer, targetFS billy.Filesystem, pkg string, osRelease string) (*git.Repository, []string, error) {
	// Clone the repository, or fail if it doesn't exist
	repo, err := git.Clone(storer, targetFS, opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to clone repo")
	}

	// First let's determine the X.Y
	matches := releaseRegex.FindStringSubmatch(osRelease)
	if len(matches) != 2 {
		return nil, nil, errors.New("failed to determine release")
	}
	partXY := matches[1]

	// Now let's find the branches
	branches, err := repo.Branches()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get branches")
	}

	var branchNames []string
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		// Let's see if the branch matches our pattern
		streamRegex, err := regexp.Compile(fmt.Sprintf(`stream-%s-.*-rhel-%s\.0`, pkg, partXY))
		if err != nil {
			return errors.Wrap(err, "failed to compile regex")
		}

		matches := streamRegex.FindStringSubmatch(ref.Name().String())
		if len(matches) != 3 {
			return nil
		}
		branchNames = append(branchNames, ref.Name().String())
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to iterate branches")
	}

	return repo, branchNames, nil
}

// getRepo returns the target repository for the module.
// This is where the payload is uploaded to.
func (s *State) getRepo(opts *git.CloneOptions, storer storage.Storer, targetFS billy.Filesystem, branches []string) (*git.Repository, error) {
	// Clone the repository, to the target filesystem.
	// We do an init, then a fetch, then a checkout
	// If the repo doesn't exist, then we init only
	repo, err := git.Init(storer, targetFS)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init repo")
	}

	var refSpecs []config.RefSpec
	for _, branch := range branches {
		refSpecs = append(refSpecs, config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%[1]s", branch)))
	}

	// Create a new remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{opts.URL},
		Fetch: refSpecs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote")
	}

	return repo, nil
}

// copyFromUpstream copies the contents of the upstream repo to the target repo
func (s *State) copyFromUpstream(opts *git.CloneOptions, repo *git.Repository, upstreamRepo *git.Repository, branches []string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}
	sourceWt, err := upstreamRepo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get upstream worktree")
	}

	var refSpecs []config.RefSpec
	for _, branch := range branches {
		refSpecs = append(refSpecs, config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%[1]s", branch)))
	}

	// Fetch the remote
	err = repo.Fetch(&git.FetchOptions{
		Auth:       opts.Auth,
		RemoteName: "origin",
		RefSpecs:   refSpecs,
	})

	// Fetch upstream
	err = upstreamRepo.Fetch(&git.FetchOptions{
		Auth:       opts.Auth,
		RemoteName: "origin",
		RefSpecs:   refSpecs,
	})
	if err != nil {
		return errors.Wrap(err, "failed to fetch upstream")
	}

	// Checkout the branch
	for _, branch := range branches {
		refName := plumbing.NewBranchReferenceName(branch)

		if err != nil {
			h := plumbing.NewSymbolicReference(plumbing.HEAD, refName)
			if err := repo.Storer.CheckAndSetReference(h, nil); err != nil {
				return errors.Wrap(err, "failed to checkout branch")
			}
		} else {
			err = wt.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(branch),
				Force:  true,
			})
			if err != nil {
				return errors.Wrap(err, "failed to checkout branch")
			}
		}

		// Copy the contents of the upstream repo to the target repo
		err = sourceWt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		})

		// Search for a {NAME}.yaml file, otherwise get the first YAML file
		yamlFileName := fmt.Sprintf("%s.yaml", s.moduleName)
		yamlFile, err := sourceWt.Filesystem.Open(yamlFileName)
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.Wrap(err, "failed to open yaml file")
			}

			// Get the first YAML file
			dir, err := sourceWt.Filesystem.ReadDir(".")
			if err != nil {
				return errors.Wrap(err, "failed to read directory")
			}

			for _, file := range dir {
				if strings.HasSuffix(file.Name(), ".yaml") {
					yamlFileName = file.Name()
					break
				}
			}

			// If we didn't find a YAML file, then we can't continue
			if yamlFileName == "" {
				return errors.New("failed to find yaml file")
			}

			yamlFile, err = sourceWt.Filesystem.Open(yamlFileName)
			if err != nil {
				return errors.Wrap(err, "failed to open yaml file")
			}
		}

		// Fully read the file
		yamlBytes, err := io.ReadAll(yamlFile)
		if err != nil {
			return errors.Wrap(err, "failed to read yaml file")
		}

		// Remove XMD metadata
		md, err := modulemd.Parse(yamlBytes)
		if err != nil {
			return errors.Wrap(err, "failed to parse yaml file")
		}

		if md.V3 != nil {
			md.V3.Data.Xmd = nil
		} else if md.V2 != nil {
			md.V2.Data.Xmd = nil
		}

		// Write the file to the target repo
		err = md.Marshal(wt.Filesystem, yamlFileName)
		if err != nil {
			return errors.Wrap(err, "failed to marshal yaml file")
		}
	}

	return nil
}
