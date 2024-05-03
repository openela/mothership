// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_worker_server

import (
	"database/sql"
	"errors"
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/temporal"
	"time"
)

func (s *UnitTestSuite) TestProcessRPMWorkflow_FullSuccess1() {
	s.env.OnActivity(testW.VerifyResourceExists, "memory://efi-rpm-macros-3-3.el8.src.rpm").Return(nil)
	s.env.OnActivity(testW.SetWorkerLastCheckinTime, mock.Anything).Return(nil)

	entry := (&mothership_db.Entry{
		Name:           base.NameGen("entries"),
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			String: "test-worker",
			Valid:  true,
		},
		State: mothershippb.Entry_ARCHIVING,
	}).ToPB()
	s.env.OnActivity(testW.CreateEntry, mock.Anything).Return(entry, nil)

	entry.EntryId = "efi-rpm-macros-3-3.el8.src"
	entry.Sha256Sum = "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28"
	s.env.OnActivity(testW.SetEntryIDFromRPM, entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum).Return(entry, nil)

	importRpmRes := &mothershippb.ImportRPMResponse{
		CommitHash:   "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		CommitUri:    testW.forge.GetCommitViewerURL("efi-rpm-macros", "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83"),
		CommitBranch: "el-8.8",
		CommitTag:    "imports/el-8.8/efi-rpm-macros-3-3.el8",
		Nevra:        "efi-rpm-macros-0:3-3.el8.aarch64",
		Pkg:          "efi-rpm-macros",
	}
	s.env.OnActivity(testW.ImportRPM, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum, entry.OsRelease).Return(importRpmRes, nil)

	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ARCHIVED, importRpmRes).Return(entry, nil)

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   entry.Sha256Sum,
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	s.env.ExecuteWorkflow(ProcessRPMWorkflow, args)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var res mothershippb.ProcessRPMResponse
	s.NoError(s.env.GetWorkflowResult(&res))
	s.Equal(entry.Name, res.Entry.Name)
	s.Equal(entry.EntryId, res.Entry.EntryId)
}

func (s *UnitTestSuite) TestProcessRPMWorkflow_OnHold_Cancel() {
	s.env.OnActivity(testW.VerifyResourceExists, "memory://efi-rpm-macros-3-3.el8.src.rpm").Return(nil)
	s.env.OnActivity(testW.SetWorkerLastCheckinTime, mock.Anything).Return(nil)

	entry := (&mothership_db.Entry{
		Name:           base.NameGen("entries"),
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			String: "test-worker",
			Valid:  true,
		},
		State: mothershippb.Entry_ARCHIVING,
	}).ToPB()
	s.env.OnActivity(testW.CreateEntry, mock.Anything).Return(entry, nil)

	entry.EntryId = "efi-rpm-macros-3-3.el8.src"
	entry.Sha256Sum = "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28"
	s.env.OnActivity(testW.SetEntryIDFromRPM, entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum).Return(entry, nil)

	importErr := errors.New("import error")
	s.env.OnActivity(testW.ImportRPM, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum, entry.OsRelease).Return(nil, importErr)

	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ON_HOLD, mock.Anything).Return(entry, nil)
	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_CANCELLED, mock.Anything).Return(entry, nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, 500*time.Millisecond)

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   entry.Sha256Sum,
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	s.env.ExecuteWorkflow(ProcessRPMWorkflow, args)

	s.True(s.env.IsWorkflowCompleted())
	s.ErrorContains(s.env.GetWorkflowError(), "canceled")
}

func (s *UnitTestSuite) TestProcessRPMWorkflow_OnHold_Success() {
	s.env.OnActivity(testW.VerifyResourceExists, "memory://efi-rpm-macros-3-3.el8.src.rpm").Return(nil)
	s.env.OnActivity(testW.SetWorkerLastCheckinTime, mock.Anything).Return(nil)

	entry := (&mothership_db.Entry{
		Name:           base.NameGen("entries"),
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			String: "test-worker",
			Valid:  true,
		},
		State: mothershippb.Entry_ARCHIVING,
	}).ToPB()
	s.env.OnActivity(testW.CreateEntry, mock.Anything).Return(entry, nil)

	entry.EntryId = "efi-rpm-macros-3-3.el8.src"
	entry.Sha256Sum = "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28"
	s.env.OnActivity(testW.SetEntryIDFromRPM, entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum).Return(&*entry, nil)

	importErr := errors.New("import error")
	importRpmRes := &mothershippb.ImportRPMResponse{
		CommitHash:   "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		CommitUri:    testW.forge.GetCommitViewerURL("efi-rpm-macros", "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83"),
		CommitBranch: "el-8.8",
		CommitTag:    "imports/el-8.8/efi-rpm-macros-3-3.el8",
		Nevra:        "efi-rpm-macros-0:3-3.el8.aarch64",
		Pkg:          "efi-rpm-macros",
	}
	shouldErrImport := true
	s.env.OnActivity(testW.ImportRPM, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum, entry.OsRelease).
		Return(func(uri string, checksum string, osRelease string) (*mothershippb.ImportRPMResponse, error) {
			if shouldErrImport {
				return nil, importErr
			}
			return importRpmRes, nil
		})

	entry.State = mothershippb.Entry_ON_HOLD
	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ON_HOLD, mock.Anything).Return(&*entry, nil)

	entry.State = mothershippb.Entry_ARCHIVED
	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ARCHIVING, mock.Anything).Return(&*entry, nil)

	entry.State = mothershippb.Entry_ARCHIVED
	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ARCHIVED, importRpmRes).Return(&*entry, nil)

	s.env.RegisterDelayedCallback(func() {
		shouldErrImport = false
		s.env.SignalWorkflow("rescue", true)
	}, 500*time.Millisecond)

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   entry.Sha256Sum,
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	s.env.ExecuteWorkflow(ProcessRPMWorkflow, args)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var res mothershippb.ProcessRPMResponse
	s.NoError(s.env.GetWorkflowResult(&res))
	s.Equal(entry.Name, res.Entry.Name)
	s.Equal(entry.EntryId, res.Entry.EntryId)
}

func (s *UnitTestSuite) TestProcessRPMWorkflow_OnHold_Error() {
	s.env.OnActivity(testW.VerifyResourceExists, "memory://efi-rpm-macros-3-3.el8.src.rpm").Return(nil)
	s.env.OnActivity(testW.SetWorkerLastCheckinTime, mock.Anything).Return(nil)

	entry := (&mothership_db.Entry{
		Name:           base.NameGen("entries"),
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			String: "test-worker",
			Valid:  true,
		},
		State: mothershippb.Entry_ARCHIVING,
	}).ToPB()
	s.env.OnActivity(testW.CreateEntry, mock.Anything).Return(entry, nil)

	entry.EntryId = "efi-rpm-macros-3-3.el8.src"
	entry.Sha256Sum = "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28"
	s.env.OnActivity(testW.SetEntryIDFromRPM, entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum).Return(entry, nil)

	importErr := errors.New("import error")
	s.env.OnActivity(testW.ImportRPM, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum, entry.OsRelease).Return(nil, importErr)

	s.env.OnActivity(testW.SetEntryState, entry.Name, mothershippb.Entry_ON_HOLD, mock.Anything).Return(entry, nil)

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   entry.Sha256Sum,
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	s.env.ExecuteWorkflow(ProcessRPMWorkflow, args)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) TestProcessRPMWorkflow_Error_DeleteEntry() {
	s.env.OnActivity(testW.VerifyResourceExists, "memory://efi-rpm-macros-3-3.el8.src.rpm").Return(nil)
	s.env.OnActivity(testW.SetWorkerLastCheckinTime, mock.Anything).Return(nil)

	entry := (&mothership_db.Entry{
		Name:           base.NameGen("entries"),
		CreateTime:     time.Now(),
		OSRelease:      "Rocky Linux release 8.8 (Green Obsidian)",
		Sha256Sum:      "518a9418fec1deaeb4c636615d8d81fb60146883c431ea15ab1127893d075d28",
		RepositoryName: "BaseOS",
		WorkerID: sql.NullString{
			String: "test-worker",
			Valid:  true,
		},
		State: mothershippb.Entry_ARCHIVING,
	}).ToPB()
	s.env.OnActivity(testW.CreateEntry, mock.Anything).Return(entry, nil)

	checksumErr := temporal.NewNonRetryableApplicationError(
		"checksum does not match",
		"checksumDoesNotMatch",
		errors.New("client submitted a checksum that does not match the resource"),
	)
	s.env.OnActivity(testW.SetEntryIDFromRPM, entry.Name, "memory://efi-rpm-macros-3-3.el8.src.rpm", entry.Sha256Sum).Return(nil, checksumErr)

	s.env.OnActivity(testW.DeleteEntry, entry.Name).Return(nil)

	args := &mothershippb.ProcessRPMArgs{
		Request: &mothershippb.ProcessRPMRequest{
			RpmUri:     "memory://efi-rpm-macros-3-3.el8.src.rpm",
			OsRelease:  "Rocky Linux release 8.8 (Green Obsidian)",
			Checksum:   entry.Sha256Sum,
			Repository: "BaseOS",
		},
		InternalRequest: &mothershippb.ProcessRPMInternalRequest{
			WorkerId: "test-worker",
		},
	}
	s.env.ExecuteWorkflow(ProcessRPMWorkflow, args)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) TestRetractEntryWorkflow_Success() {
	entry := base.NameGen("entries")
	s.env.OnActivity(testW.SetEntryState, entry, mothershippb.Entry_RETRACTING, mock.Anything).Return(nil, nil)

	res := &mshipadminpb.RetractEntryResponse{
		Name: entry,
	}
	s.env.OnActivity(testW.RetractEntry, entry).Return(res, nil)

	s.env.OnActivity(testW.SetEntryState, entry, mothershippb.Entry_RETRACTED, mock.Anything).Return(nil, nil)

	s.env.ExecuteWorkflow(RetractEntryWorkflow, entry)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UnitTestSuite) TestRetractEntryWorkflow_Failed_RevertToArchived() {
	entry := base.NameGen("entries")
	s.env.OnActivity(testW.SetEntryState, entry, mothershippb.Entry_RETRACTING, mock.Anything).Return(nil, nil)

	anyErr := errors.New("any error")
	s.env.OnActivity(testW.RetractEntry, entry).Return(nil, anyErr)

	s.env.OnActivity(testW.SetEntryState, entry, mothershippb.Entry_ARCHIVED, mock.Anything).Return(nil, nil)

	s.env.ExecuteWorkflow(RetractEntryWorkflow, entry)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}
