package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/urfave/cli/v2"
)

func openRepo(rootFS billy.Filesystem, path string) (*git.Repository, error) {
	dot, err := rootFS.Chroot(path + "/.git")
	if err != nil {
		return nil, err
	}
	storer := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	repo, err := git.Open(storer, rootFS)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func getRootPatch(rootFS billy.Filesystem) (string, error) {
	// Open the git repository
	repo, err := openRepo(rootFS, ".")
	if err != nil {
		return "", err
	}

	// Check if there are any changes
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	status, err := wt.Status()
	if err != nil {
		return "", err
	}

	return status.String(), nil
}

func getTarballPatch(rootFS billy.Filesystem) (string, error) {
	// Open the git repository
	repo, err := openRepo(rootFS, "tarball")
	if err != nil {
		return "", err
	}

	// Check if there are any changes
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	status, err := wt.Status()
	if err != nil {
		return "", err
	}

	return status.String(), nil
}

func generate(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.Exit("usage: pkgpatch generate MESSAGE", 1)
	}
	msg := ctx.Args().First()
	fmt.Println(msg)

	// Verify that this is ran from a .pkgpatchroot directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(".pkgpatchroot"); os.IsNotExist(err) {
		return fmt.Errorf("not a .pkgpatchroot directory")
	}

	// Open the root filesystem
	rootFS := osfs.New(currentDir)

	// Get the root patch
	rootPatch, err := getRootPatch(rootFS)
	if err != nil {
		return err
	}
	fmt.Println(rootPatch)

	// Get the tarball patch
	if _, err := rootFS.Stat("tarball"); err == nil {
		tarballPatch, err := getTarballPatch(rootFS)
		if err != nil {
			return err
		}
		fmt.Println(tarballPatch)
	}

	return nil
}
