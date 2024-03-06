package mothership_worker_server

import (
	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/bugtracker"
	"github.com/openela/mothership/base/forge"
	"github.com/openela/mothership/base/storage"
	"golang.org/x/crypto/openpgp"
)

type Worker struct {
	db         *base.DB
	storage    storage.Storage
	gpgKeys    openpgp.EntityList
	forge      forge.Forge
	bugtracker bugtracker.Bugtracker
	rolling    bool
	publicURI  string
}

// New creates a new Worker
// todo(mustafa): This is really ugly, we should probably just use the struct above directly
func New(db *base.DB, storage storage.Storage, gpgKeys openpgp.EntityList, forge forge.Forge, bugtracker bugtracker.Bugtracker, rolling bool, publicURI string) *Worker {
	return &Worker{
		db:         db,
		storage:    storage,
		gpgKeys:    gpgKeys,
		forge:      forge,
		bugtracker: bugtracker,
		rolling:    rolling,
		publicURI:  publicURI,
	}
}
