package mothership_worker_server

import (
	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/forge"
	"github.com/openela/mothership/base/storage"
	"golang.org/x/crypto/openpgp"
)

type Worker struct {
	db      *base.DB
	storage storage.Storage
	gpgKeys openpgp.EntityList
	forge   forge.Forge
	rolling bool
}

func New(db *base.DB, storage storage.Storage, gpgKeys openpgp.EntityList, forge forge.Forge, rolling bool) *Worker {
	return &Worker{
		db:      db,
		storage: storage,
		gpgKeys: gpgKeys,
		forge:   forge,
		rolling: rolling,
	}
}
