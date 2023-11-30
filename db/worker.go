package mothership_db

import (
	"database/sql"
	"github.com/openela/mothership/base"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Worker struct {
	PikaTableName      string `pika:"workers"`
	PikaDefaultOrderBy string `pika:"-create_time"`

	Name            string       `db:"name"`
	CreateTime      time.Time    `db:"create_time" pika:"omitempty"`
	WorkerID        string       `db:"worker_id"`
	LastCheckinTime sql.NullTime `db:"last_checkin_time"`
	ApiSecret       string       `db:"api_secret"`
}

func (w *Worker) GetID() string {
	return w.Name
}

func (w *Worker) ToPB() *mshipadminpb.Worker {
	return &mshipadminpb.Worker{
		Name:            w.Name,
		WorkerId:        w.WorkerID,
		CreateTime:      timestamppb.New(w.CreateTime),
		LastCheckinTime: base.SqlNullTime(w.LastCheckinTime),
	}
}
