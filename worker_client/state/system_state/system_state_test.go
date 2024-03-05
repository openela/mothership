package system_state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	storage_memory "github.com/openela/mothership/base/storage/memory"
	"github.com/openela/mothership/worker_client/state"
	"github.com/stretchr/testify/require"
)

var initPathVar string

func resetEnvNew(newBinDir string) {
	if initPathVar == "" {
		initPathVar = os.Getenv("PATH")
	}
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", newBinDir, initPathVar)); err != nil {
		panic(err)
	}
}

func newReposync(script string) string {
	tempDir, err := os.MkdirTemp("", "reposync")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s/reposync", tempDir), []byte(script), 0755)
	if err != nil {
		panic(err)
	}

	resetEnvNew(tempDir)

	return tempDir
}

func newFiles(in map[string]string) (string, map[string]string) {
	tempDir, err := os.MkdirTemp("", "files")
	if err != nil {
		panic(err)
	}

	out := make(map[string]string)
	for k, v := range in {
		outPath := fmt.Sprintf("%s/%s", tempDir, k)
		dir := filepath.Dir(outPath)
		if dir != tempDir {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		}

		err = os.WriteFile(outPath, []byte(v), 0644)
		if err != nil {
			panic(err)
		}

		hash, err := sha256OfFile(outPath)
		if err != nil {
			panic(err)
		}
		out[k] = "/" + hash
	}

	return tempDir, out
}

func sha256Hash(in string) string {
	h := sha256.New()
	_, err := h.Write([]byte(in))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeState(file string, s *state.PackageState) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(s); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestState_FetchNewPackageState_AllNew(t *testing.T) {
	dir, err := os.MkdirTemp("", "state_1")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	filesDir, out := newFiles(map[string]string{
		"test1.rpm": "test1",
		"test2.rpm": "test2",
	})
	defer os.RemoveAll(filesDir)

	reposyncDir := newReposync(`#!/bin/sh
exit 0`)
	defer os.RemoveAll(reposyncDir)

	storage := storage_memory.New(osfs.New("/"), "")
	systemState, err := New(&Args{
		FilePath:    filepath.Join(dir, "state.json"),
		PathToSrcs:  filesDir,
		ReposToSync: []string{"test"},
		Storage:     storage,
	})

	require.Nil(t, systemState.FetchNewPackageState())
	dirtyObjects := systemState.GetDirtyObjects()
	require.Equal(t, []string{"/" + sha256Hash("test1"), "/" + sha256Hash("test2")}, dirtyObjects)
	err = systemState.WritePackageState()
	require.Nil(t, err)

	test1, err := storage.Get(out["test1.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test1", string(test1))

	test2, err := storage.Get(out["test2.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test2", string(test2))

	packageState := systemState.GetState()
	require.Equal(t, out["test1.rpm"][1:], packageState.Packages["test1.rpm"])
	require.Equal(t, out["test2.rpm"][1:], packageState.Packages["test2.rpm"])
}

func TestState_FetchNewPackageState_ExistingModifiedNew(t *testing.T) {
	dir, err := os.MkdirTemp("", "state_2")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	filesDir, out := newFiles(map[string]string{
		"test1.rpm": "test1",
		"test2.rpm": "test2-changed",
		"test3.rpm": "test3-new",
	})
	defer os.RemoveAll(filesDir)

	reposyncDir := newReposync(`#!/bin/sh
exit 0`)
	defer os.RemoveAll(reposyncDir)

	storage := storage_memory.New(osfs.New("/"), "")

	_, err = storage.PutBytes(out["test1.rpm"], []byte("test1"))
	require.Nil(t, err)

	_, err = storage.PutBytes("/"+sha256Hash("test2"), []byte("test2"))
	require.Nil(t, err)

	filePath := filepath.Join(dir, "state.json")
	writeState(filePath, &state.PackageState{
		Packages: map[string]string{
			"test1.rpm": sha256Hash("test1"),
			"test2.rpm": sha256Hash("test2"),
		},
	})

	systemState, err := New(&Args{
		FilePath:    filePath,
		PathToSrcs:  filesDir,
		ReposToSync: []string{"test"},
		Storage:     storage,
	})

	require.Nil(t, systemState.FetchNewPackageState())
	dirtyObjects := systemState.GetDirtyObjects()
	require.Equal(t, []string{"/" + sha256Hash("test2-changed"), "/" + sha256Hash("test3-new")}, dirtyObjects)
	err = systemState.WritePackageState()
	require.Nil(t, err)

	test1, err := storage.Get(out["test1.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test1", string(test1))

	test2, err := storage.Get(out["test2.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test2-changed", string(test2))

	test3, err := storage.Get(out["test3.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test3-new", string(test3))

	packageState := systemState.GetState()
	require.Equal(t, out["test1.rpm"][1:], packageState.Packages["test1.rpm"])
	require.Equal(t, out["test2.rpm"][1:], packageState.Packages["test2.rpm"])
	require.Equal(t, out["test3.rpm"][1:], packageState.Packages["test3.rpm"])
}

func TestState_FetchNewPackageState_AllExisting(t *testing.T) {
	dir, err := os.MkdirTemp("", "state_3")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	filesDir, out := newFiles(map[string]string{
		"test1.rpm": "test1",
		"test2.rpm": "test2",
	})
	defer os.RemoveAll(filesDir)

	reposyncDir := newReposync(`#!/bin/sh
exit 0`)
	defer os.RemoveAll(reposyncDir)

	storage := storage_memory.New(osfs.New("/"), "")

	_, err = storage.PutBytes(out["test1.rpm"], []byte("test1"))
	require.Nil(t, err)

	_, err = storage.PutBytes(out["test2.rpm"], []byte("test2"))
	require.Nil(t, err)

	filePath := filepath.Join(dir, "state.json")
	writeState(filePath, &state.PackageState{
		Packages: map[string]string{
			"test1.rpm": sha256Hash("test1"),
			"test2.rpm": sha256Hash("test2"),
		},
	})

	systemState, err := New(&Args{
		FilePath:    filePath,
		PathToSrcs:  filesDir,
		ReposToSync: []string{"test"},
		Storage:     storage,
	})

	require.Nil(t, systemState.FetchNewPackageState())
	dirtyObjects := systemState.GetDirtyObjects()
	require.Len(t, dirtyObjects, 0)
	err = systemState.WritePackageState()
	require.Nil(t, err)

	test1, err := storage.Get(out["test1.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test1", string(test1))

	test2, err := storage.Get(out["test2.rpm"])
	require.Nil(t, err)
	require.Equal(t, "test2", string(test2))

	packageState := systemState.GetState()
	require.Equal(t, out["test1.rpm"][1:], packageState.Packages["test1.rpm"])
	require.Equal(t, out["test2.rpm"][1:], packageState.Packages["test2.rpm"])
}
