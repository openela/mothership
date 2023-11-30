package mothership_worker_server

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/forge"
	mothership_db "github.com/openela/mothership/db"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"strings"
	"time"
)

// getRepo gets a git repository from a remote
// It clones into an in-memory filesystem
func getRepo(remote string, auth transport.AuthMethod) (*git.Repository, error) {
	// Just use in memory storage for all repos
	storer := memory.NewStorage()
	fs := memfs.New()
	repo, err := git.Init(storer, fs)
	if err != nil {
		return nil, err
	}

	// Add a new remote
	refspec := config.RefSpec("refs/*:refs/*")
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{remote},
		Fetch: []config.RefSpec{refspec},
	})
	if err != nil {
		return nil, err
	}

	// Fetch all the refs from the remote
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		RefSpecs:   []config.RefSpec{refspec},
		Tags:       git.AllTags,
		Auth:       auth,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// clonePathToFS clones a path from one filesystem to another
func clonePathToFS(fromFS billy.Filesystem, toFS billy.Filesystem, rootPath string) error {
	// check if root directory exists
	_, err := fromFS.Stat(rootPath)
	if err != nil {
		// we don't care if the directory doesn't exist
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	// read the root directory
	rootDir, err := fromFS.ReadDir(rootPath)
	if err != nil {
		return err
	}

	// iterate over the files
	for _, file := range rootDir {
		// get the file path
		filePath := rootPath + "/" + file.Name()

		// check if the file is a directory
		if file.IsDir() {
			// create the directory in the toFS
			err = toFS.MkdirAll(filePath, 0755)
			if err != nil {
				return err
			}

			// recursively call this function
			err = clonePathToFS(fromFS, toFS, filePath)
			if err != nil {
				return err
			}
		} else {
			// open the file
			f, err := fromFS.OpenFile(filePath, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			// create the file in the toFS
			toFile, err := toFS.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			defer toFile.Close()

			// copy the file contents
			_, err = io.Copy(toFile, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// clonePatchesToTemporaryFS clones the PATCHES directory to a temporary filesystem
// PATCHES directory is the only directory that should survive a retraction
func clonePatchesToTemporaryFS(currentFS billy.Filesystem) (billy.Filesystem, error) {
	// create a new in-memory filesystem
	fs := memfs.New()

	// clone the current filesystem to the new filesystem
	err := clonePathToFS(currentFS, fs, "PATCHES")
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func resetRepoToPoint(repo *git.Repository, authenticator *forge.Authenticator, commit string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	// Let's find out the commit before the commit we want to revert
	log, err := repo.Log(&git.LogOptions{
		From:  plumbing.NewHash(commit),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get log")
	}

	// log.Next() x2 should be the commit we want to revert
	targetCommit, err := log.Next()
	if err != nil {
		return errors.Wrap(err, "failed to get next commit x1")
	}
	resetToCommit, err := log.Next()
	if err != nil {
		return errors.Wrap(err, "failed to get next commit x2")
	}

	// Also get all commits that touches the PATCHES directory
	// until the commit we want to revert
	firstLog, err := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		PathFilter: func(s string) bool {
			// Only include PATCHES
			if strings.HasPrefix(s, "PATCHES") {
				return true
			}

			return false
		},
		// Limit to until the commit we want to revert
		Since: &resetToCommit.Author.When,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get log")
	}

	// Get all authors of the commits, since we're going to copy PATCHES
	// back into the repo
	var commits []*object.Commit
	for {
		c, err := firstLog.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to get next commit")
		}
		// If the commit was created before the resetToCommit, then we don't want to include it
		if c.Author.When.Before(targetCommit.Author.When) {
			// Breaking because the commits are sorted by date, so if we encounter a commit that was created before the resetToCommit,
			// then all the following commits will be created before the resetToCommit
			break
		}
		commits = append(commits, c)
	}

	// Copy PATCHES into a temporary filesystem
	patchesFS, err := clonePatchesToTemporaryFS(wt.Filesystem)
	if err != nil {
		return errors.Wrap(err, "failed to clone PATCHES")
	}

	// reset the repo
	err = wt.Reset(&git.ResetOptions{
		Commit: resetToCommit.Hash,
		Mode:   git.HardReset,
	})
	if err != nil {
		return errors.Wrap(err, "failed to reset repo")
	}

	// Copy PATCHES back into the repo
	err = clonePathToFS(patchesFS, wt.Filesystem, "PATCHES")
	if err != nil {
		return errors.Wrap(err, "failed to copy PATCHES")
	}

	// If there are diffs, then create a commit consisting of the joined messages and authors
	// of the commits that touches the PATCHES directory
	if len(commits) > 0 {
		// Add the files
		_, err = wt.Add(".")
		if err != nil {
			return errors.Wrap(err, "failed to add PATCHES")
		}

		// Get the commit message
		commitMsg := "Retract \"" + targetCommit.Message + "\"\n\nFast-forwarded following commits:\n"
		for _, c := range commits {
			commitMsg += c.Message + "\n"
		}
		commitMsg += "\n"

		// Add the authors as "Co-authored-by"
		authors := make(map[string]bool)
		for _, c := range commits {
			authors[c.Author.Name+" <"+c.Author.Email+">"] = true
		}
		for author := range authors {
			commitMsg += "Co-authored-by: " + author + "\n"
		}

		// Create the commit
		_, err = wt.Commit(commitMsg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  authenticator.AuthorName,
				Email: authenticator.AuthorEmail,
				When:  time.Now(),
			},
		})
		if err != nil {
			return errors.Wrap(err, "failed to commit")
		}
	}

	return err
}

func (w *Worker) RetractEntry(name string) (*mshipadminpb.RetractEntryResponse, error) {
	entry, err := base.Q[mothership_db.Entry](w.db).F("name", name).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get entry: %v", err)
		return nil, status.Error(codes.Internal, "failed to get entry")
	}

	if entry == nil {
		return nil, temporal.NewNonRetryableApplicationError(
			"entry not found",
			"entryNotFound",
			nil,
		)
	}

	// Get the repo
	remote := w.forge.GetRemote(entry.PackageName)
	auth, err := w.forge.GetAuthenticator()
	if err != nil {
		base.LogErrorf("failed to get forge authenticator: %v", err)
		return nil, status.Error(codes.Internal, "failed to get forge authenticator")
	}
	repo, err := getRepo(remote, auth.AuthMethod)
	if err != nil {
		base.LogErrorf("failed to get repo: %v", err)
		return nil, status.Error(codes.Internal, "failed to get repo")
	}

	// Checkout the entry branch
	wt, err := repo.Worktree()
	if err != nil {
		base.LogErrorf("failed to get worktree: %v", err)
		return nil, status.Error(codes.Internal, "failed to get worktree")
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(entry.CommitBranch),
		Force:  true,
	})
	if err != nil {
		base.LogErrorf("failed to checkout branch: %v", err)
		return nil, status.Error(codes.Internal, "failed to checkout branch")
	}

	// Reset the repo to the commit before the commit we want to revert
	err = resetRepoToPoint(repo, auth, entry.CommitHash)
	if err != nil {
		base.LogErrorf("failed to reset repo: %v", err)
		return nil, status.Error(codes.Internal, "failed to reset repo")
	}

	// Push the changes
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Force:      true,
		Auth:       auth.AuthMethod,
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/" + entry.CommitBranch + ":refs/heads/" + entry.CommitBranch),
		},
	})
	if err != nil {
		base.LogErrorf("failed to push changes: %v", err)
		return nil, status.Error(codes.Internal, "failed to push changes")
	}

	return &mshipadminpb.RetractEntryResponse{
		Name: entry.Name,
	}, nil
}
