package system_state

import (
	"encoding/json"
	"github.com/openela/mothership/base/storage"
	"github.com/openela/mothership/worker_client/state"
	"log/slog"
	"os"
	"os/exec"
	"sync"
)

// State implements the System State mode, where the worker
// state is stored in a JSON file on disk.
type State struct {
	mutex sync.Mutex
	f     *os.File

	state   *state.PackageState
	storage storage.Storage

	majorVersion int
	pathToSrcs   string
	reposToSync  []string
}

type Args struct {
	FilePath    string          `yaml:"file_path"`
	PathToSrcs  string          `yaml:"path_to_srcs"`
	ReposToSync []string        `yaml:"repos_to_sync"`
	Storage     storage.Storage `yaml:"-"`
}

func New(args *Args) (*State, error) {
	s := &State{
		storage:     args.Storage,
		pathToSrcs:  args.PathToSrcs,
		reposToSync: args.ReposToSync,
	}

	_, err := os.Stat(args.FilePath)
	if os.IsNotExist(err) {
		// File does not exist, create it
		f, err := os.OpenFile(args.FilePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		s.f = f
		s.state = &state.PackageState{
			Packages: make(map[string]string),
		}

		return s, nil
	}

	f, err := os.OpenFile(args.FilePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// Unmarshal the state
	var st state.PackageState
	err = json.NewDecoder(f).Decode(&st)
	if err != nil {
		return nil, err
	}

	if st.Packages == nil {
		st.Packages = make(map[string]string)
	}

	s.f = f
	s.state = &st

	return s, nil
}

func (s *State) writeToDisk() error {
	_, err := s.f.Seek(0, 0)
	if err != nil {
		return err
	}

	return json.NewEncoder(s.f).Encode(s.state)
}

func (s *State) modifyPackages(merge map[string]string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k, v := range merge {
		s.state.Packages[k] = v
	}

	err := s.writeToDisk()
	if err != nil {
		return err
	}

	return nil
}

func (s *State) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.f.Close()
}

// UpdatePackageState does a reposync to update the state of the packages.
// The pathToSrcs is used as the base directory for the reposync.
// The reposToSync are the repositories to sync.
// Then all paths are walked and the SHA256 hash of each file is calculated.
// Any changed files are updated in the state and uploaded to storage with path `/<majorVersion>/<relativePath>`.
func (s *State) UpdatePackageState() error {
	args := []string{"-n", "--refresh"}
	for _, repo := range s.reposToSync {
		args = append(args, "--repo", repo)
	}

	slog.Info("Running reposync", "args", args)

	cmd := exec.Command("reposync", args...)
	cmd.Dir = s.pathToSrcs
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
