package mothership_db

import (
	"database/sql"
	"github.com/openela/mothership/base"
	mothershippb "github.com/openela/mothership/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Batch struct {
	PikaTableName      string `pika:"batches"`
	PikaDefaultOrderBy string `pika:"-create_time"`

	Name          string         `db:"name"`
	BatchID       string         `db:"batch_id"`
	WorkerID      string         `db:"worker_id"`
	CreateTime    time.Time      `db:"create_time" pika:"omitempty"`
	UpdateTime    time.Time      `db:"create_time" pika:"omitempty"`
	SealTime      sql.NullTime   `db:"create_time"`
	BugtrackerURI sql.NullString `db:"bugtracker_uri"`
}

func (b *Batch) GetID() string {
	return b.Name
}

func (b *Batch) ToPB() *mothershippb.Batch {
	return &mothershippb.Batch{
		Name:          b.Name,
		BatchId:       b.BatchID,
		WorkerId:      b.WorkerID,
		CreateTime:    timestamppb.New(b.CreateTime),
		UpdateTime:    timestamppb.New(b.CreateTime),
		SealTime:      base.SqlNullTime(b.SealTime),
		BugtrackerUri: base.SqlNullString(b.BugtrackerURI),
	}
}
