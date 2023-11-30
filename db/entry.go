package mothership_db

import (
	"database/sql"
	"github.com/openela/mothership/base"
	mothershippb "github.com/openela/mothership/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Entry struct {
	PikaTableName      string `pika:"entries"`
	PikaDefaultOrderBy string `pika:"-create_time"`

	Name           string                   `db:"name"`
	EntryID        string                   `db:"entry_id"`
	CreateTime     time.Time                `db:"create_time" pika:"omitempty"`
	OSRelease      string                   `db:"os_release"`
	Sha256Sum      string                   `db:"sha256_sum"`
	RepositoryName string                   `db:"repository_name"`
	WorkerID       sql.NullString           `db:"worker_id"`
	BatchName      sql.NullString           `db:"batch_name"`
	UserEmail      sql.NullString           `db:"user_email"`
	CommitURI      string                   `db:"commit_uri"`
	CommitHash     string                   `db:"commit_hash"`
	CommitBranch   string                   `db:"commit_branch"`
	CommitTag      string                   `db:"commit_tag"`
	State          mothershippb.Entry_State `db:"state"`
	PackageName    string                   `db:"package_name"`
}

func (e *Entry) GetID() string {
	return e.Name
}

func (e *Entry) ToPB() *mothershippb.Entry {
	return &mothershippb.Entry{
		Name:         e.Name,
		EntryId:      e.EntryID,
		CreateTime:   timestamppb.New(e.CreateTime),
		OsRelease:    e.OSRelease,
		Sha256Sum:    e.Sha256Sum,
		Repository:   e.RepositoryName,
		WorkerId:     base.SqlNullString(e.WorkerID),
		Batch:        base.SqlNullString(e.BatchName),
		UserEmail:    base.SqlNullString(e.UserEmail),
		CommitUri:    e.CommitURI,
		CommitHash:   e.CommitHash,
		CommitBranch: e.CommitBranch,
		CommitTag:    e.CommitTag,
		State:        e.State,
		Pkg:          e.PackageName,
	}
}
