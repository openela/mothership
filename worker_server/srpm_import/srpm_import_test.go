// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package srpm_import

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	storage_memory "github.com/openela/mothership/base/storage/memory"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/openpgp"
)

func TestFromFile(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	require.Nil(t, s.Close())
}

func TestFromFile_SignatureOK(t *testing.T) {
	keyF, err := os.Open("testdata/RPM-GPG-KEY-Rocky-8")
	require.Nil(t, err)

	testKey, err := openpgp.ReadArmoredKeyRing(keyF)
	require.Nil(t, err)

	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false, testKey...)
	require.Nil(t, err)
	require.NotNil(t, s)
	require.Nil(t, s.Close())
}

func TestFromFile_SignatureFail(t *testing.T) {
	keyF, err := os.Open("testdata/RPM-GPG-KEY-Rocky-9")
	require.Nil(t, err)

	testKey, err := openpgp.ReadArmoredKeyRing(keyF)
	require.Nil(t, err)

	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false, testKey...)
	require.NotNil(t, err)
	require.Nil(t, s)
	require.Equal(t, "failed to verify RPM: keyid 15af5dac6d745a60 not found", err.Error())
}

func TestDetermineLookasideBlobs_Empty(t *testing.T) {
	s, err := FromFile("testdata/basesystem-11-5.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())
	require.Equal(t, 0, len(s.lookasideBlobs))
}

func TestDetermineLookasideBlobs_NotEmpty_Tarball(t *testing.T) {
	s, err := FromFile("testdata/bash-4.4.20-4.el8_6.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())
	require.Equal(t, 1, len(s.lookasideBlobs))
}

func TestUploadLookaside_Empty(t *testing.T) {
	s, err := FromFile("testdata/basesystem-11-5.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())

	// we can use memfs since we're not actually writing anything
	fs := memfs.New()
	lookaside := storage_memory.New(fs)
	require.Nil(t, s.uploadLookasideBlobs(lookaside))

	fi, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 0, len(fi))
}

func TestUploadLookaside_NotEmpty(t *testing.T) {
	s, err := FromFile("testdata/bash-4.4.20-4.el8_6.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())

	fs := osfs.New("/")
	lookaside := storage_memory.New(fs)
	require.Nil(t, s.uploadLookasideBlobs(lookaside))

	ok, err := lookaside.Exists("d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb")
	require.Nil(t, err)
	require.True(t, ok)
}

func TestUploadLookaside_NotEmpty_OnlyOnceForHash(t *testing.T) {
	s, err := FromFile("testdata/bash-4.4.20-4.el8_6.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())

	fs := osfs.New("/")
	lookaside := storage_memory.New(fs)
	require.Nil(t, s.uploadLookasideBlobs(lookaside))

	ok, err := lookaside.Exists("d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb")
	require.Nil(t, err)
	require.True(t, ok)

	_, err = lookaside.PutBytes("d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb", []byte("test"))
	require.Nil(t, err)

	require.Nil(t, s.uploadLookasideBlobs(lookaside))

	bts, err := lookaside.Get("d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb")
	require.Nil(t, err)
	// This verifies that the upload only checked for the hash once, and didn't
	// overwrite the existing blob.
	require.Equal(t, []byte("test"), bts)
}

func TestWriteMetadataFile(t *testing.T) {
	s, err := FromFile("testdata/bash-4.4.20-4.el8_6.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()
	require.Nil(t, s.determineLookasideBlobs())

	fs := memfs.New()
	require.Nil(t, s.writeMetadataFile(fs))

	fi, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 2, len(fi))
	require.Equal(t, ".bash.metadata", fi[0].Name())
	require.Equal(t, ".gitignore", fi[1].Name())

	f, err := fs.Open(".bash.metadata")
	require.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	require.Nil(t, err)

	require.Equal(t, "d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb SOURCES/bash-4.4.tar.gz\n", string(buf[:n]))

	f, err = fs.Open(".gitignore")
	require.Nil(t, err)

	buf = make([]byte, 1024)
	n, err = f.Read(buf)
	require.Nil(t, err)

	require.Equal(t, "SOURCES/bash-4.4.tar.gz\n", string(buf[:n]))
}

func TestExpandLayout(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	fs := memfs.New()
	require.Nil(t, s.ExpandLayout(fs))

	fi, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 2, len(fi))
	require.Equal(t, "SOURCES", fi[0].Name())
	require.Equal(t, "SPECS", fi[1].Name())

	fi, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)

	require.Equal(t, 2, len(fi))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", fi[0].Name())
	require.Equal(t, "efi-rpm-macros-3.tar.bz2", fi[1].Name())

	fi, err = fs.ReadDir("SPECS")
	require.Nil(t, err)

	require.Equal(t, 1, len(fi))
	require.Equal(t, "efi-rpm-macros.spec", fi[0].Name())
}

func TestExpandLayout_CopyMode(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	// Create a fake shell script
	shellPath := filepath.Join(s.tempDir, "test.sh")
	shellFile, err := os.OpenFile(shellPath, os.O_RDWR|os.O_CREATE, 0755)
	require.Nil(t, err)
	_, err = shellFile.WriteString("#!/bin/bash\necho 'test'\n")
	require.Nil(t, err)
	require.Nil(t, shellFile.Close())

	fs := memfs.New()
	require.Nil(t, s.ExpandLayout(fs))

	fi, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 2, len(fi))
	require.Equal(t, "SOURCES", fi[0].Name())
	require.Equal(t, "SPECS", fi[1].Name())

	fi, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 3, len(fi))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", fi[0].Name())
	require.Equal(t, "efi-rpm-macros-3.tar.bz2", fi[1].Name())
	require.Equal(t, "test.sh", fi[2].Name())
	// Verify mode
	info, err := fs.Stat(filepath.Join("SOURCES", "test.sh"))
	require.Nil(t, err)
	require.Equal(t, os.FileMode(0755), info.Mode())

	fi, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(fi))
	require.Equal(t, "efi-rpm-macros.spec", fi[0].Name())
}

func TestWriteMetadataExpandLayout(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	fs := memfs.New()
	require.Nil(t, s.determineLookasideBlobs())
	require.Nil(t, s.writeMetadataFile(fs))
	require.Nil(t, s.ExpandLayout(fs))

	fi, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(fi))
	require.Equal(t, ".efi-rpm-macros.metadata", fi[0].Name())
	require.Equal(t, ".gitignore", fi[1].Name())
	require.Equal(t, "SOURCES", fi[2].Name())
	require.Equal(t, "SPECS", fi[3].Name())

	fi, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)

	require.Equal(t, 1, len(fi))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", fi[0].Name())

	fi, err = fs.ReadDir("SPECS")
	require.Nil(t, err)

	require.Equal(t, 1, len(fi))
	require.Equal(t, "efi-rpm-macros.spec", fi[0].Name())

	f, err := fs.Open(".efi-rpm-macros.metadata")
	require.Nil(t, err)

	buf, err := io.ReadAll(f)
	require.Nil(t, err)

	require.Equal(t, "f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86 SOURCES/efi-rpm-macros-3.tar.bz2\n", string(buf))
}

func TestGetStreamSuffix(t *testing.T) {
	s, err := FromFile("testdata/nginx-1.14.1-9.module+el8.4.0+542+81547229.src.rpm", false)
	require.Nil(t, err)

	suffix, err := s.getStreamSuffix()
	require.Nil(t, err)
	require.Equal(t, "-stream-1.14", suffix)
}

func TestGetStreamSuffix_NotModuleComponent(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)

	suffix, err := s.getStreamSuffix()
	require.Nil(t, err)
	require.Equal(t, "", suffix)
}

func TestGetRepo_New(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)
}

func TestGetRepo_OSRelease(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "Rocky Linux release 8.8 (Green Obsidian)")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8.8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)
}

func TestGetRepo_OSRelease_Error(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "X invalid 1.1")
	require.NotNil(t, err)
	require.Nil(t, repo)
	require.Equal(t, "", branch)
	require.Equal(t, "invalid OS release X invalid 1.1", err.Error())
}

func TestGetRepo_New_Rolling(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", true)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)
}

func TestGetRepo_OSRelease_Rolling(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", true)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "Rocky Linux release 8.8 (Green Obsidian)")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)
}

func TestGetRepo_Existing(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	// Push a commit to the bare repo
	newTempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(newTempDir)

	osfs2 := osfs.New(newTempDir)
	dot2, err := osfs2.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp2 := filesystem.NewStorage(dot2, cache.NewObjectLRUDefault())

	repo, err := git.InitWithOptions(filesystemTemp2, osfs2, git.InitOptions{
		DefaultBranch: "refs/heads/el-8",
	})
	require.Nil(t, err)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{tempDir},
	})
	require.Nil(t, err)
	w, err := repo.Worktree()
	require.Nil(t, err)
	f, err := w.Filesystem.Create("testfile")
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)
	_, err = w.Add("testfile")
	require.Nil(t, err)
	_, err = w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@openela.org",
			When:  time.Now(),
		},
	})
	require.Nil(t, err)
	err = repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/el-8:refs/heads/el-8"},
	})
	require.Nil(t, err)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "test commit", obj.Message)
}

func TestCleanTargetRepo_Existing(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	// Push a commit to the bare repo
	newTempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(newTempDir)

	osfs2 := osfs.New(newTempDir)
	dot2, err := osfs2.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp2 := filesystem.NewStorage(dot2, cache.NewObjectLRUDefault())

	repo, err := git.InitWithOptions(filesystemTemp2, osfs2, git.InitOptions{
		DefaultBranch: "refs/heads/el-8",
	})
	require.Nil(t, err)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{tempDir},
	})
	require.Nil(t, err)
	w, err := repo.Worktree()
	require.Nil(t, err)
	f, err := w.Filesystem.Create("testfile")
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)
	_, err = w.Add("testfile")
	require.Nil(t, err)
	_, err = w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@openela.org",
			When:  time.Now(),
		},
	})
	require.Nil(t, err)
	err = repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/el-8:refs/heads/el-8"},
	})
	require.Nil(t, err)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	wt, err := repo.Worktree()
	require.Nil(t, err)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "test commit", obj.Message)

	// Clean repo
	require.Nil(t, s.cleanTargetRepo(wt, "."))

	// Verify empty
	ls, err := wt.Filesystem.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 0, len(ls))
}

func TestPopulateTargetRepo_New(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)

	// Populate repo
	inMemory := storage_memory.New(osfs.New("/"))
	require.Nil(t, s.populateTargetRepo(repo, fs, inMemory, branch))

	// Verify commit
	objIter, err = repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "import efi-rpm-macros-3-3.el8", obj.Message)

	// Verify tag
	tagIter, err := repo.Tags()
	require.Nil(t, err)
	tag, err := tagIter.Next()
	require.Nil(t, err)
	require.Equal(t, "imports/el-8/efi-rpm-macros-3-3.el8", tag.Name().Short())

	// Verify metadata
	f, err := fs.Open(".efi-rpm-macros.metadata")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86 SOURCES/efi-rpm-macros-3.tar.bz2\n", string(buf))

	// Verify layout
	ls, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(ls))
	require.Equal(t, ".efi-rpm-macros.metadata", ls[0].Name())
	require.Equal(t, ".gitignore", ls[1].Name())
	require.Equal(t, "SOURCES", ls[2].Name())
	require.Equal(t, "SPECS", ls[3].Name())

	ls, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", ls[0].Name())

	ls, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "efi-rpm-macros.spec", ls[0].Name())
}

func TestPopulateTargetRepo_Existing(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	// Push a commit to the bare repo
	newTempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(newTempDir)

	osfs2 := osfs.New(newTempDir)
	dot2, err := osfs2.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp2 := filesystem.NewStorage(dot2, cache.NewObjectLRUDefault())

	repo, err := git.InitWithOptions(filesystemTemp2, osfs2, git.InitOptions{
		DefaultBranch: "refs/heads/el-8",
	})
	require.Nil(t, err)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{tempDir},
	})
	require.Nil(t, err)
	w, err := repo.Worktree()
	require.Nil(t, err)
	f, err := w.Filesystem.Create("testfile")
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)
	_, err = w.Add("testfile")
	require.Nil(t, err)
	_, err = w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@openela.org",
			// We're subtracting an hour here because the commit time is
			// truncated to the nearest second, and we want to make sure
			// that the commit time is different from the one we're going
			// to create below.
			When: time.Now().Add(-1 * time.Hour),
		},
	})
	require.Nil(t, err)
	err = repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/el-8:refs/heads/el-8"},
	})
	require.Nil(t, err)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "test commit", obj.Message)

	// Populate repo
	inMemory := storage_memory.New(osfs.New("/"))
	require.Nil(t, s.populateTargetRepo(repo, fs, inMemory, branch))

	// Verify commit (second one)
	var sortedCommits []*object.Commit
	objIter, err = repo.CommitObjects()
	require.Nil(t, err)
	obj, err = objIter.Next()
	require.Nil(t, err)
	sortedCommits = append(sortedCommits, obj)
	obj, err = objIter.Next()
	require.Nil(t, err)
	sortedCommits = append(sortedCommits, obj)

	sort.Slice(sortedCommits, func(i, j int) bool {
		return sortedCommits[i].Author.When.After(sortedCommits[j].Author.When)
	})
	require.Equal(t, "import efi-rpm-macros-3-3.el8", sortedCommits[0].Message)

	// Verify tag
	tagIter, err := repo.Tags()
	require.Nil(t, err)
	tag, err := tagIter.Next()
	require.Nil(t, err)
	require.Equal(t, "imports/el-8/efi-rpm-macros-3-3.el8", tag.Name().Short())

	// Verify metadata
	f, err = fs.Open(".efi-rpm-macros.metadata")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86 SOURCES/efi-rpm-macros-3.tar.bz2\n", string(buf))

	// Verify layout
	ls, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(ls))
	require.Equal(t, ".efi-rpm-macros.metadata", ls[0].Name())
	require.Equal(t, ".gitignore", ls[1].Name())
	require.Equal(t, "SOURCES", ls[2].Name())
	require.Equal(t, "SPECS", ls[3].Name())

	ls, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", ls[0].Name())

	ls, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "efi-rpm-macros.spec", ls[0].Name())
}

func TestPatchTargetRepo(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)

	wt, err := repo.Worktree()
	require.Nil(t, err)

	f, err := wt.Filesystem.Create("PATCHES/test.cfg")
	require.Nil(t, err)
	_, err = f.Write([]byte(`
spec_change {
  changelog {
    author_name: "Mustafa Gezen"
    author_email: "mustafa@rockylinux.org"
    message: "test patch1"
  }
}
`))
	require.Nil(t, err)

	// Populate repo, this will also run patchTargetRepo
	inMemory := storage_memory.New(osfs.New("/"))
	require.Nil(t, s.populateTargetRepo(repo, fs, inMemory, branch))

	// Verify spec file contains "test patch1"
	f, err = fs.Open("SPECS/efi-rpm-macros.spec")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Contains(t, string(buf), "test patch1")
}

func TestPatchTargetRepo_CarryOver(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "Rocky Linux release 8.8 (Green Obsidian)")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8.8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)

	wt, err := repo.Worktree()
	require.Nil(t, err)

	f, err := wt.Filesystem.Create("PATCHES/test.cfg")
	require.Nil(t, err)
	_, err = f.Write([]byte(`
spec_change {
  changelog {
    author_name: "Mustafa Gezen"
    author_email: "mustafa@rockylinux.org"
    message: "test patch1"
  }
}
`))
	require.Nil(t, err)

	// Populate repo, this will also run patchTargetRepo
	inMemory := storage_memory.New(osfs.New("/"))
	require.Nil(t, s.populateTargetRepo(repo, fs, inMemory, branch))

	// Verify spec file contains "test patch1"
	f, err = fs.Open("SPECS/efi-rpm-macros.spec")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Contains(t, string(buf), "test patch1")

	fs = memfs.New()
	repo, branch, err = s.getRepo(opts, storer, fs, "Rocky Linux release 8.9 (Green Obsidian)")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8.9", branch)

	// Populate repo, this will also run patchTargetRepo
	inMemory = storage_memory.New(osfs.New("/"))
	require.Nil(t, s.populateTargetRepo(repo, fs, inMemory, branch))

	// Verify spec file contains "test patch1"
	f, err = fs.Open("SPECS/efi-rpm-macros.spec")
	require.Nil(t, err)
	buf, err = io.ReadAll(f)
	require.Nil(t, err)
	require.Contains(t, string(buf), "test patch1")
}

func TestPatchTargetRepo_Failed(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	storer := memory.NewStorage()
	fs := memfs.New()
	opts := &git.CloneOptions{
		URL: "file://" + tempDir,
	}
	repo, branch, err := s.getRepo(opts, storer, fs, "")
	require.Nil(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "el-8", branch)

	// Verify empty
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	_, err = objIter.Next()
	require.Equal(t, io.EOF, err)

	wt, err := repo.Worktree()
	require.Nil(t, err)

	f, err := wt.Filesystem.Create("PATCHES/test.patch")
	require.Nil(t, err)
	_, err = f.Write([]byte(`
diff --git a/SPECS/efi-rpm-macros.spec b/SPECS/efi-rpm-macros.spec
index 0f2b2f6..f002f60 100644
--- a/SPECS/efi-rpm-macros.spec
+++ b/SPECS/efi-rpm-macros.spec
@@ -1,6 +1,6 @@
asd
`))
	require.Nil(t, err)

	f, err = wt.Filesystem.Create("PATCHES/test.cfg")
	require.Nil(t, err)
	_, err = f.Write([]byte(`
patch {
  file: "PATCHES/test.patch"
}
`))
	require.Nil(t, err)

	// Populate repo, this will also run patchTargetRepo
	inMemory := storage_memory.New(osfs.New("/"))
	err = s.populateTargetRepo(repo, fs, inMemory, branch)
	require.NotNil(t, err)
	require.Equal(t, "failed to patch target repo: COULD_NOT_PARSE_PATCH_FILE:PATCHES/test.patch: failed to apply directives", err.Error())
}

func TestPushTargetRepo(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)
	require.NotNil(t, s)
	defer func() {
		require.Nil(t, s.Close())
	}()

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	// Push a commit to the bare repo
	newTempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(newTempDir)

	osfs2 := osfs.New(newTempDir)
	dot2, err := osfs2.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp2 := filesystem.NewStorage(dot2, cache.NewObjectLRUDefault())

	repo, err := git.InitWithOptions(filesystemTemp2, osfs2, git.InitOptions{
		DefaultBranch: "refs/heads/el-8",
	})
	require.Nil(t, err)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{tempDir},
	})
	require.Nil(t, err)
	w, err := repo.Worktree()
	require.Nil(t, err)
	f, err := w.Filesystem.Create("testfile")
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)
	_, err = w.Add("testfile")
	require.Nil(t, err)
	_, err = w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@openela.org",
			When:  time.Now(),
		},
	})
	require.Nil(t, err)
	require.Nil(t, s.pushTargetRepo(repo, &git.PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/el-8:refs/heads/el-8"},
	}))

	// Verify testfile is still there
	f, err = osfs2.Open("testfile")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "test", string(buf))
}

func TestImport1_New(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	opts := &git.CloneOptions{
		URL: tempDir,
	}
	storer := memory.NewStorage()
	fs := memfs.New()
	lookaside := storage_memory.New(osfs.New("/"))
	_, err = s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	// Open repo
	repo, err := git.PlainOpen(tempDir)
	require.Nil(t, err)
	// Switch to el8 branch
	w, err := repo.Worktree()
	require.Nil(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/el-8",
	})
	require.Nil(t, err)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "import efi-rpm-macros-3-3.el8", obj.Message)

	// Verify tag
	tagIter, err := repo.Tags()
	require.Nil(t, err)
	tag, err := tagIter.Next()
	require.Nil(t, err)
	require.Equal(t, "imports/el-8/efi-rpm-macros-3-3.el8", tag.Name().Short())

	// Verify metadata
	f, err := fs.Open(".efi-rpm-macros.metadata")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86 SOURCES/efi-rpm-macros-3.tar.bz2\n", string(buf))

	// Verify layout
	ls, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(ls))
	require.Equal(t, ".efi-rpm-macros.metadata", ls[0].Name())
	require.Equal(t, ".gitignore", ls[1].Name())
	require.Equal(t, "SOURCES", ls[2].Name())
	require.Equal(t, "SPECS", ls[3].Name())

	ls, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", ls[0].Name())

	ls, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "efi-rpm-macros.spec", ls[0].Name())

	// Verify lookaside
	ok, err := lookaside.Exists("f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86")
	require.Nil(t, err)
	require.True(t, ok)
}

func TestImport1_New_Rolling(t *testing.T) {
	s, err := FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", true)
	require.Nil(t, err)

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	opts := &git.CloneOptions{
		URL: tempDir,
	}
	storer := memory.NewStorage()
	fs := memfs.New()
	lookaside := storage_memory.New(osfs.New("/"))
	_, err = s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	// Open repo
	repo, err := git.PlainOpen(tempDir)
	require.Nil(t, err)
	// Switch to el8 branch
	w, err := repo.Worktree()
	require.Nil(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/el8",
	})
	require.Nil(t, err)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "import efi-rpm-macros-3-3.el8", obj.Message)

	// Verify tag
	tagIter, err := repo.Tags()
	require.Nil(t, err)
	tag, err := tagIter.Next()
	require.Nil(t, err)
	require.Equal(t, "imports/el8/efi-rpm-macros-3-3.el8", tag.Name().Short())

	// Verify metadata
	f, err := fs.Open(".efi-rpm-macros.metadata")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86 SOURCES/efi-rpm-macros-3.tar.bz2\n", string(buf))

	// Verify layout
	ls, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(ls))
	require.Equal(t, ".efi-rpm-macros.metadata", ls[0].Name())
	require.Equal(t, ".gitignore", ls[1].Name())
	require.Equal(t, "SOURCES", ls[2].Name())
	require.Equal(t, "SPECS", ls[3].Name())

	ls, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "0001-macros.efi-srpm-make-all-of-our-macros-always-expand.patch", ls[0].Name())

	ls, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "efi-rpm-macros.spec", ls[0].Name())

	// Verify lookaside
	ok, err := lookaside.Exists("f002f60baed7a47ca3e98b8dd7ece2f7352dac9ffab7ae3557eb56b481ce2f86")
	require.Nil(t, err)
	require.True(t, ok)
}

func TestImport2_New(t *testing.T) {
	s, err := FromFile("testdata/bash-4.4.20-4.el8_6.src.rpm", false)
	require.Nil(t, err)

	tempDir, err := os.MkdirTemp("", "peridot-srpm-import-test-*")
	require.Nil(t, err)

	// Create a bare repo in tempDir
	osfsTemp := osfs.New(tempDir)
	dot, err := osfsTemp.Chroot(".git")
	require.Nil(t, err)
	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	require.Nil(t, filesystemTemp.Init())
	_, err = git.Init(filesystemTemp, nil)
	require.Nil(t, err)

	opts := &git.CloneOptions{
		URL: tempDir,
	}
	storer := memory.NewStorage()
	fs := memfs.New()
	lookaside := storage_memory.New(osfs.New("/"))
	_, err = s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	// Open repo
	repo, err := git.PlainOpen(tempDir)
	require.Nil(t, err)
	// Switch to el8 branch
	w, err := repo.Worktree()
	require.Nil(t, err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/el-8",
	})
	require.Nil(t, err)

	// Verify commit
	objIter, err := repo.CommitObjects()
	require.Nil(t, err)
	obj, err := objIter.Next()
	require.Nil(t, err)
	require.Equal(t, "import bash-4.4.20-4.el8_6", obj.Message)

	// Verify tag
	tagIter, err := repo.Tags()
	require.Nil(t, err)
	tag, err := tagIter.Next()
	require.Nil(t, err)
	require.Equal(t, "imports/el-8/bash-4.4.20-4.el8_6", tag.Name().Short())

	// Verify metadata
	f, err := fs.Open(".bash.metadata")
	require.Nil(t, err)
	buf, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb SOURCES/bash-4.4.tar.gz\n", string(buf))

	// Verify layout
	ls, err := fs.ReadDir(".")
	require.Nil(t, err)
	require.Equal(t, 4, len(ls))
	require.Equal(t, ".bash.metadata", ls[0].Name())
	require.Equal(t, ".gitignore", ls[1].Name())
	require.Equal(t, "SOURCES", ls[2].Name())
	require.Equal(t, "SPECS", ls[3].Name())

	ls, err = fs.ReadDir("SOURCES")
	require.Nil(t, err)
	require.Equal(t, 61, len(ls))

	ls, err = fs.ReadDir("SPECS")
	require.Nil(t, err)
	require.Equal(t, 1, len(ls))
	require.Equal(t, "bash.spec", ls[0].Name())

	// Verify lookaside
	ok, err := lookaside.Exists("d86b3392c1202e8ff5a423b302e6284db7f8f435ea9f39b5b1b20fd3ac36dfcb")
	require.Nil(t, err)
	require.True(t, ok)
}
