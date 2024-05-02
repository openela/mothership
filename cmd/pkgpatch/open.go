package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/openela/mothership/worker_server/srpm_import"
	"github.com/urfave/cli/v2"
)

const baseSrpmURL = "https://ax8edlmsvvfp.compat.objectstorage.us-phoenix-1.oraclecloud.com/mship-srpm1"

type entryChecksum struct {
	Sha256Sum string `json:"sha256Sum"`
}

func downloadResource(url, dest string) error {
	slog.Info("Downloading resource", "url", url, "dest", dest)
	// Fetch the resource from the URL
	// Write the response body to a file in the dest directory
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func open(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return cli.Exit("usage: pkgpatch open ENTRY_ID", 1)
	}

	entryID := ctx.Args().First()

	// Fetch from https://imports.openela.org/api/v1/entries/ENTRY_ID
	// Parse the response body as JSON
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	slog.Info("Fetching entry", "entryID", entryID)

	req, err := http.NewRequest(http.MethodGet, "https://imports.openela.org/api/v1/entries/"+entryID, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var checksum entryChecksum
	if err := json.NewDecoder(res.Body).Decode(&checksum); err != nil {
		return err
	}

	slog.Info("Entry fetched", "entryID", entryID, "checksum", checksum.Sha256Sum)

	// Create a temporary directory and download the resource
	tmpDir, err := os.MkdirTemp("", "mothership-worker-server-import-rpm-*")
	if err != nil {
		return err
	}

	dest := tmpDir + "/resource.rpm"
	downloadURL := baseSrpmURL + "/" + checksum.Sha256Sum
	if err := downloadResource(downloadURL, dest); err != nil {
		return err
	}

	// Do an SRPM expand
	slog.Info("Expanding SRPM", "dest", tmpDir)
	state, err := srpm_import.FromFile(dest, true)
	if err != nil {
		return err
	}

	// Create expand directory
	expandDir := tmpDir + "/expand"
	if err := os.Mkdir(expandDir, 0755); err != nil {
		return err
	}

	expandFS := osfs.New(expandDir)
	err = state.ExpandLayout(expandFS)
	if err != nil {
		return err
	}

	// Expand the largest tarball, that will be a separate .git repo
	// but first we need to find the largest tarball
	var largestTarball string
	var largestSize int64
	ls, err := expandFS.ReadDir("SOURCES")
	if err != nil {
		return err
	}
	for _, file := range ls {
		if strings.Contains(file.Name(), ".tar") && file.Size() > largestSize {
			largestSize = file.Size()
			largestTarball = file.Name()
		}
	}

	if largestSize > 0 {
		slog.Info("Tarball found", "tarball", largestTarball, "size", largestSize)

		// Create a base tarball directory
		tarballDir := expandDir + "/tarball"
		if err := os.Mkdir(tarballDir, 0755); err != nil {
			return err
		}

		slog.Info("Expanding tarball and creating a git repository (this might take a while)")

		// Expand the largest tarball
		cmd := exec.Command("tar", "-x", "--strip-components=1", "-C", tarballDir, "-f", filepath.Join(expandDir, "SOURCES", largestTarball))
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		// Delete the tarball
		err = expandFS.Remove("SOURCES/" + largestTarball)
		if err != nil {
			return err
		}

		// Init a new git repository in that directory
		err = expandFS.MkdirAll("tarball/.git", 0755)
		if err != nil {
			return err
		}

		dot, err := expandFS.Chroot("tarball/.git")
		if err != nil {
			return err
		}
		storer := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
		tarballFS, err := expandFS.Chroot("tarball")
		if err != nil {
			return err
		}
		repo, err := git.Init(storer, tarballFS)
		if err != nil {
			return err
		}

		// Add all files to the git repository
		w, err := repo.Worktree()
		if err != nil {
			return err
		}

		_, err = w.Add(".")
		if err != nil {
			return err
		}

		_, err = w.Commit("Initial commit", &git.CommitOptions{})
		if err != nil {
			return err
		}
	}

	// Init a new git repository in that directory
	err = expandFS.MkdirAll(".git", 0755)
	if err != nil {
		return err
	}
	dot, err := expandFS.Chroot(".git")
	if err != nil {
		return err
	}
	storer := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	repo, err := git.Init(storer, expandFS)
	if err != nil {
		return err
	}

	// Create a gitignore file for the tarball directory
	gitignore := []byte("tarball/\n")
	f, err := expandFS.OpenFile(".gitignore", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(gitignore); err != nil {
		return err
	}

	// Create a .pkgpatchroot file
	f, err = expandFS.OpenFile(".pkgpatchroot", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Add all files to the git repository
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit("Initial commit", &git.CommitOptions{})
	if err != nil {
		return err
	}

	expandDirBold := color.New(color.Bold).Sprint(expandDir)
	generateBold := color.New(color.Bold).Sprint("pkgpatch generate MESSAGE")
	color.Green(`
Success!
A new workspace has been created for entry %s

Open %s in your favorite editor.
When you're done run %s to generate the patches.

NOTICE: Currently, the "generate" command is not implemented.
You need to generate the patches manually in both repositories.
If you make a change in "expand", then commit the changes and run "git format-patch HEAD^1"
same applies to the "tarball" repository.
`, entryID, expandDirBold, generateBold)

	return nil
}
