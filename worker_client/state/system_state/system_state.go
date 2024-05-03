// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package system_state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/openela/mothership/base/storage"
	"github.com/openela/mothership/worker_client/state"
	"github.com/pkg/errors"
)

// State implements the System State mode, where the worker
// state is stored in a JSON file on disk.
type State struct {
	mutex sync.Mutex

	state   *state.PackageState
	storage storage.Storage

	majorVersion      int
	pathToSrcs        string
	reposToSync       []string
	dirtyPackageState map[string]string
	filePath          string
}

type Args struct {
	FilePath     string          `yaml:"file_path"`
	PathToSrcs   string          `yaml:"path_to_srcs"`
	ReposToSync  []string        `yaml:"repos_to_sync"`
	WorkerSecret string          `yaml:"worker_secret"`
	Storage      storage.Storage `yaml:"-"`
}

var ignoreList = []string{
	"redhat-logos",
	"redhat-release",
	"kernel-rt",
	"rhncfg",
	"rhn-custom-info",
	"rhnpush",
	"shim",
}

func New(args *Args) (*State, error) {
	s := &State{
		storage:           args.Storage,
		pathToSrcs:        args.PathToSrcs,
		reposToSync:       args.ReposToSync,
		dirtyPackageState: map[string]string{},
		filePath:          args.FilePath,
	}

	_, err := os.Stat(args.FilePath)
	if os.IsNotExist(err) {
		// File does not exist, create it
		f, err := os.OpenFile(args.FilePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		// Create the state
		s.state = &state.PackageState{
			Packages: make(map[string]string),
		}

		// Write the state to disk
		err = json.NewEncoder(f).Encode(s.state)
		if err != nil {
			return nil, err
		}

		err = f.Close()
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	// Check if path to srcs exists
	_, err = os.Stat(args.PathToSrcs)
	if os.IsNotExist(err) {
		err = os.MkdirAll(args.PathToSrcs, 0755)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.OpenFile(args.FilePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Unmarshal the state
	var st state.PackageState
	err = json.NewDecoder(f).Decode(&st)
	if err != nil {
		return nil, err
	}

	if st.Packages == nil {
		st.Packages = make(map[string]string)
	}

	s.state = &st

	return s, nil
}

func (s *State) writeToDisk() error {
	f, err := os.OpenFile(s.filePath, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(s.state)
}

func (s *State) modifyPackages(merge map[string]string) error {
	for k, v := range merge {
		s.state.Packages[k] = v
	}

	err := s.writeToDisk()
	if err != nil {
		return err
	}

	return nil
}

func (s *State) GetDirtyObjects() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var dirtyObjects []string
	for _, v := range s.dirtyPackageState {
		dirtyObjects = append(dirtyObjects, "/"+v)
	}

	sort.Strings(dirtyObjects)

	return dirtyObjects
}

// FetchNewPackageState does a reposync to update the state of the packages.
// The pathToSrcs is used as the base directory for the reposync.
// The reposToSync are the repositories to sync.
// Then all paths are walked and the SHA256 hash of each file is calculated.
// Any changed files are updated in the state and uploaded to storage with path `/<majorVersion>/<relativePath>`.
func (s *State) FetchNewPackageState() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	args := []string{"-n", "--refresh", "--source"}
	for _, repo := range s.reposToSync {
		args = append(args, "--repo", repo)
	}

	slog.Info("Running reposync", "args", args)

	cmd := exec.Command("reposync", args...)
	cmd.Dir = s.pathToSrcs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	// After running reposync, we can map the full SHA256 hash of each file
	// to the relative path of sync directory.
	// We can then compare this to the state and upload any changed files, update
	// the state and upload the files to storage.
	allPaths := map[string]string{}
	err = filepath.WalkDir(s.pathToSrcs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// If not rpm, skip
		if filepath.Ext(path) != ".rpm" {
			return nil
		}

		// Get the base path
		basePath := filepath.Base(path)

		// Skip if in ignore list
		for _, ignore := range ignoreList {
			if strings.HasPrefix(basePath, ignore) {
				return nil
			}
		}

		// If an entry for an NVR exists, then we can skip
		if _, ok := s.state.Packages[basePath]; ok {
			return nil
		}

		// Get the SHA256 hash of the file
		hash, err := sha256OfFile(path)
		if err != nil {
			return err
		}

		allPaths[path] = hash

		return nil
	})
	if err != nil {
		return err
	}

	// Upload the changed files to storage
	var wg sync.WaitGroup
	var lock sync.Mutex
	var syncErr error
	for path, hash := range allPaths {
		wg.Add(1)

		go func(path string, hash string) {
			defer wg.Done()

			// Get the base path
			basePath := filepath.Base(path)

			objectPath := "/" + hash

			// File has changed, upload it
			slog.Info("uploading file", "path", path, "objectPath", objectPath)

			// Check if it exists, if so skip
			exists, err := s.storage.Exists(objectPath)
			if err != nil {
				if syncErr == nil {
					syncErr = err
				} else {
					syncErr = errors.Wrap(syncErr, "failed to check if object exists: "+objectPath)
				}
			}
			if exists {
				slog.Info("object already exists", "objectPath", objectPath)
				lock.Lock()
				s.dirtyPackageState[basePath] = hash
				lock.Unlock()
				return
			}

			_, err = s.storage.Put(objectPath, path)
			if err != nil {
				if syncErr == nil {
					syncErr = err
				} else {
					syncErr = errors.Wrap(syncErr, "failed to upload object: "+objectPath)
				}
				return
			}

			// Update the state
			lock.Lock()
			s.dirtyPackageState[basePath] = hash
			lock.Unlock()
		}(path, hash)
	}

	wg.Wait()
	if syncErr != nil {
		return syncErr
	}

	return nil
}

func (s *State) WritePackageState() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.dirtyPackageState) == 0 {
		return nil
	}

	err := s.modifyPackages(s.dirtyPackageState)
	if err != nil {
		return err
	}

	s.dirtyPackageState = map[string]string{}

	return nil
}

func (s *State) GetState() *state.PackageState {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.state
}

func sha256OfFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
