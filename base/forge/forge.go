package forge

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"time"
)

type Authenticator struct {
	transport.AuthMethod

	AuthorName  string
	AuthorEmail string
	// Expires is the time when the token expires.
	// So it can be used to cache the token.
	Expires time.Time
}

type Forge interface {
	GetAuthenticator() (*Authenticator, error)
	GetRemote(repo string) string
	GetCommitViewerURL(repo string, commit string) string
	EnsureRepositoryExists(auth *Authenticator, repo string) error
	WithNamespace(namespace string) Forge
}
