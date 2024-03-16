package mothership_db

import (
	"database/sql"
	"time"

	"github.com/openela/mothership/base"
	mothershippb "github.com/openela/mothership/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Batch struct {
	PikaTableName      string `pika:"batches"`
	PikaDefaultOrderBy string `pika:"-create_time"`

	Name          string         `db:"name"`
	BatchID       sql.NullString `db:"batch_id" pika:"omitempty"`
	WorkerID      string         `db:"worker_id"`
	CreateTime    time.Time      `db:"create_time" pika:"omitempty"`
	UpdateTime    time.Time      `db:"update_time" pika:"omitempty"`
	SealTime      sql.NullTime   `db:"seal_time"`
	BugtrackerURI sql.NullString `db:"bugtracker_uri"`
}

func (b *Batch) GetID() string {
	return b.Name
}

func (b *Batch) ToPB() *mothershippb.Batch {
	return &mothershippb.Batch{
		Name:          b.Name,
		BatchId:       b.BatchID.String,
		WorkerId:      b.WorkerID,
		CreateTime:    timestamppb.New(b.CreateTime),
		UpdateTime:    timestamppb.New(b.CreateTime),
		SealTime:      base.SqlNullTime(b.SealTime),
		BugtrackerUri: base.SqlNullString(b.BugtrackerURI),
	}
}

type BatchView struct {
	PikaTableName      string `pika:"batches_view"`
	PikaDefaultOrderBy string `pika:"-create_time"`

	Name          string         `db:"name"`
	BatchID       sql.NullString `db:"batch_id" pika:"omitempty"`
	WorkerID      string         `db:"worker_id"`
	CreateTime    time.Time      `db:"create_time" pika:"omitempty"`
	UpdateTime    time.Time      `db:"update_time" pika:"omitempty"`
	SealTime      sql.NullTime   `db:"seal_time"`
	BugtrackerURI sql.NullString `db:"bugtracker_uri"`
	EntryCount    int            `db:"entry_count" pika:"omitempty"`
}

func (b *BatchView) GetID() string {
	return b.Name
}

func (b *BatchView) ToPB() *mothershippb.Batch {
	return &mothershippb.Batch{
		Name:          b.Name,
		BatchId:       b.BatchID.String,
		WorkerId:      b.WorkerID,
		CreateTime:    timestamppb.New(b.CreateTime),
		UpdateTime:    timestamppb.New(b.CreateTime),
		SealTime:      base.SqlNullTime(b.SealTime),
		BugtrackerUri: base.SqlNullString(b.BugtrackerURI),
		EntryCount:    int32(b.EntryCount),
	}
}
