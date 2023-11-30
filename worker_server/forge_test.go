package mothership_worker_server

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/openela/mothership/base/forge"
	"path/filepath"
	"time"
)

type inMemoryForge struct {
	localTempDir        string
	repos               map[string]bool
	remoteBaseURL       string
	invalidUsernamePass bool
	noAuthMethod        bool
	namespace           string
}

func (f *inMemoryForge) GetAuthenticator() (*forge.Authenticator, error) {
	ret := &forge.Authenticator{
		AuthMethod: &transport_http.BasicAuth{
			Username: "user",
			Password: "pass",
		},
		AuthorName:  "Test User",
		AuthorEmail: "test@openela.org",
		Expires:     time.Now().Add(time.Hour),
	}

	if f.noAuthMethod {
		ret.AuthMethod = nil
	} else if f.invalidUsernamePass {
		ret.AuthMethod = &transport_http.BasicAuth{
			Username: "invalid",
			Password: "invalid",
		}
	}

	return ret, nil
}

func (f *inMemoryForge) GetRemote(repo string) string {
	return fmt.Sprintf("file://%s/%s%s", f.localTempDir, f.namespace, repo)
}

func (f *inMemoryForge) GetCommitViewerURL(repo string, commit string) string {
	return f.remoteBaseURL + "/" + f.namespace + repo + "/commit/" + commit
}

func (f *inMemoryForge) EnsureRepositoryExists(auth *forge.Authenticator, repo string) error {
	// Try casting auth.AuthMethod to *transport_http.BasicAuth
	// If it fails, return an error
	authx, ok := auth.AuthMethod.(*transport_http.BasicAuth)
	if !ok {
		return errors.New("auth failed")
	}
	if authx.Username != "user" || authx.Password != "pass" {
		return errors.New("username or password incorrect")
	}

	if f.repos[repo] {
		return nil
	}

	osfsTemp := osfs.New(filepath.Join(f.localTempDir, repo))
	dot, err := osfsTemp.Chroot(".git")
	if err != nil {
		return err
	}

	filesystemTemp := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	err = filesystemTemp.Init()
	if err != nil {
		return err
	}

	_, err = git.Init(filesystemTemp, nil)
	if err != nil {
		return err
	}

	f.repos[repo] = true
	return nil
}

func (f *inMemoryForge) WithNamespace(ns string) forge.Forge {
	x := *f
	x.namespace = ns
	if x.namespace[len(x.namespace)-1] != '/' {
		x.namespace += "/"
	}
	return &x
}
