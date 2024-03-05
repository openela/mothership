package mothership_worker_server

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/forge"
	storage_memory "github.com/openela/mothership/base/storage/memory"
	mothership_db "github.com/openela/mothership/db"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/worker_server/srpm_import"
	"github.com/stretchr/testify/require"
)

func TestGetRepo(t *testing.T) {
	s, err := srpm_import.FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
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

	// Check that the repo was created
	repo, err := getRepo("file://"+tempDir, nil)
	require.Nil(t, err)
	require.NotNil(t, repo)

	// Check that the repo was cloned
	commits, err := repo.CommitObjects()
	require.Nil(t, err)
	require.NotNil(t, commits)
	commit, err := commits.Next()
	require.Nil(t, err)
	require.NotNil(t, commit)
	require.Equal(t, "import efi-rpm-macros-3-3.el8", commit.Message)
}

func TestClonePathToFS(t *testing.T) {
	fromFS := memfs.New()
	toFS := memfs.New()

	f, err := fromFS.OpenFile("test", os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	err = clonePathToFS(fromFS, toFS, ".")
	require.Nil(t, err)

	f, err = toFS.OpenFile("test", os.O_RDONLY, 0644)
	require.Nil(t, err)
	b, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "test", string(b))
}

func TestClonePathToFS_RootPath(t *testing.T) {
	fromFS := memfs.New()
	toFS := memfs.New()

	f, err := fromFS.OpenFile("test", os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	f, err = fromFS.OpenFile("testdir/foo", os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	_, err = f.Write([]byte("bar"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	err = clonePathToFS(fromFS, toFS, "testdir")
	require.Nil(t, err)

	f, err = toFS.OpenFile("testdir/foo", os.O_RDONLY, 0644)
	require.Nil(t, err)
	b, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "bar", string(b))

	f, err = toFS.OpenFile("test", os.O_RDONLY, 0644)
	require.NotNil(t, err)
	require.Equal(t, "file does not exist", err.Error())
}

func TestPatchesToTemporaryFS(t *testing.T) {
	fromFS := memfs.New()

	f, err := fromFS.OpenFile("PATCHES/test", os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	_, err = f.Write([]byte("test"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	toFS, err := clonePatchesToTemporaryFS(fromFS)
	require.Nil(t, err)

	f, err = toFS.OpenFile("PATCHES/test", os.O_RDONLY, 0644)
	require.Nil(t, err)
	b, err := io.ReadAll(f)
	require.Nil(t, err)
	require.Equal(t, "test", string(b))
}

func TestResetRepoToPoint_OneCommit(t *testing.T) {
	s, err := srpm_import.FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
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
	firstImport, err := s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	repo, err := getRepo("file://"+tempDir, nil)
	require.Nil(t, err)

	// Get wt and checkout the correct branch
	wt, err := repo.Worktree()
	require.Nil(t, err)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(firstImport.Branch),
		Force:  true,
	})
	require.Nil(t, err)

	err = resetRepoToPoint(
		repo,
		&forge.Authenticator{AuthorName: "test", AuthorEmail: "test@rockylinux.org"},
		firstImport.Commit.Hash.String(),
		"el-8",
	)
	require.Nil(t, err)

	// Check that only rollback commit exists
	log, err := repo.Log(&git.LogOptions{})
	require.Nil(t, err)
	commit, err := log.Next()
	require.Nil(t, err)
	require.NotNil(t, commit)
	require.Equal(t, "Rollback to empty state", commit.Message)
	commit, err = log.Next()
	require.NotNil(t, err)
	require.Equal(t, "EOF", err.Error())
	require.Nil(t, commit)
}

func TestResetRepoToPoint_TwoCommits(t *testing.T) {
	s, err := srpm_import.FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
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
	firstImport, err := s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	storer2 := memory.NewStorage()
	fs2 := memfs.New()
	secondImport, err := s.Import(opts, storer2, fs2, lookaside, "")
	require.Nil(t, err)

	repo, err := getRepo("file://"+tempDir, nil)
	require.Nil(t, err)

	// Get wt and checkout the correct branch
	wt, err := repo.Worktree()
	require.Nil(t, err)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(secondImport.Branch),
		Force:  true,
	})
	require.Nil(t, err)

	err = resetRepoToPoint(
		repo,
		&forge.Authenticator{AuthorName: "test", AuthorEmail: "test@rockylinux.org"},
		secondImport.Commit.Hash.String(),
		"el-8",
	)
	require.Nil(t, err)

	// Check that only the first commit exists
	log, err := repo.Log(&git.LogOptions{})
	require.Nil(t, err)
	commit, err := log.Next()
	require.Nil(t, err)
	require.NotNil(t, commit)
	require.Equal(t, firstImport.Commit.Hash.String(), commit.Hash.String())
	commit, err = log.Next()
	require.NotNil(t, err)
	require.Equal(t, "EOF", err.Error())
	require.Nil(t, commit)
}

func TestResetRepoToPoint_TwoCommits_CommitAfterRetractPoint(t *testing.T) {
	s, err := srpm_import.FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
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
	firstImport, err := s.Import(opts, storer, fs, lookaside, "")
	require.Nil(t, err)

	storer2 := memory.NewStorage()
	fs2 := memfs.New()
	secondImport, err := s.Import(opts, storer2, fs2, lookaside, "")
	require.Nil(t, err)

	repo, err := getRepo("file://"+tempDir, nil)
	require.Nil(t, err)

	// Get wt and checkout the correct branch
	wt, err := repo.Worktree()
	require.Nil(t, err)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(secondImport.Branch),
		Force:  true,
	})
	require.Nil(t, err)

	// Create a commit after the commit we want to reset to
	f, err := wt.Filesystem.Create("PATCHES/test.cfg")
	require.Nil(t, err)
	_, err = f.Write([]byte("lookaside: { file: \"test.png\" }"))
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)

	_, err = wt.Add("PATCHES/test.cfg")
	require.Nil(t, err)

	stableTime := time.Now()
	_, err = wt.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Mustafa Gezen",
			Email: "mustafa@rockylinux.org",
			When:  stableTime,
		},
	})
	require.Nil(t, err)

	err = resetRepoToPoint(
		repo,
		&forge.Authenticator{AuthorName: "test", AuthorEmail: "test@rockylinux.org"},
		secondImport.Commit.Hash.String(),
		"el-8",
	)
	require.Nil(t, err)

	// Check that only the first commit exists
	log, err := repo.Log(&git.LogOptions{})
	require.Nil(t, err)
	commit, err := log.Next()
	require.Nil(t, err)
	require.NotNil(t, commit)
	msg := `Retract "import efi-rpm-macros-3-3.el8"

Fast-forwarded following commits:
test commit

Co-authored-by: Mustafa Gezen <mustafa@rockylinux.org>
`
	require.Equal(t, msg, commit.Message)
	commit, err = log.Next()
	require.Nil(t, err)
	require.Equal(t, firstImport.Commit.Hash.String(), commit.Hash.String())
	commit, err = log.Next()
	require.NotNil(t, err)
	require.Equal(t, "EOF", err.Error())
	require.Nil(t, commit)
}

func TestWorker_RetractEntry(t *testing.T) {
	s, err := srpm_import.FromFile("testdata/efi-rpm-macros-3-3.el8.src.rpm", false)
	require.Nil(t, err)

	tempDir := filepath.Join(tempDirForge, "efi-rpm-macros")
	err = os.RemoveAll(tempDir)
	require.Nil(t, err)
	err = os.MkdirAll(tempDir, 0755)
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
	_, err = s.Import(opts, storer, fs, lookaside, "Rocky Linux release 8.8 (Green Obsidian)")
	require.Nil(t, err)

	storer2 := memory.NewStorage()
	fs2 := memfs.New()
	secondImport, err := s.Import(opts, storer2, fs2, lookaside, "Rocky Linux release 8.8 (Green Obsidian)")
	require.Nil(t, err)

	// Create entry
	entry := &mothership_db.Entry{
		Name:           base.NameGen("entries"),
		EntryID:        "efi-rpm-macros-3-3.el8.src",
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			Valid:  true,
			String: "test-worker",
		},
		CommitURI:    "file://" + tempDir,
		CommitHash:   secondImport.Commit.Hash.String(),
		CommitBranch: "el-8.8",
		CommitTag:    "imports/el-8.8/efi-rpm-macros-3-3.el8",
		State:        mothershippb.Entry_ARCHIVED,
		PackageName:  "efi-rpm-macros",
	}
	require.Nil(t, base.Q[mothership_db.Entry](testW.db).Create(entry))

	// Retract entry
	res, err := testW.RetractEntry(entry.Name)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, entry.Name, res.Name)
}
