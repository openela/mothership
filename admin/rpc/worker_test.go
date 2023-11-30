package mothershipadmin_rpc

import (
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGetWorker_Empty(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	worker, err := s.GetWorker(testContext(), &mshipadminpb.GetWorkerRequest{})
	require.NotNil(t, err)
	require.Nil(t, worker)
	expectedErr := status.Error(codes.NotFound, "worker not found")
	require.Equal(t, expectedErr.Error(), err.Error())
}

func TestGetWorker_One(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Create(&mothership_db.Worker{
		Name:      "test",
		WorkerID:  "test-id",
		ApiSecret: "secret",
	}))
	defer func() {
		require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	}()

	worker, err := s.GetWorker(testContext(), &mshipadminpb.GetWorkerRequest{
		Name: "test",
	})
	require.Nil(t, err)
	require.Equal(t, "test", worker.Name)
	require.Equal(t, "test-id", worker.WorkerId)
	require.Empty(t, worker.ApiSecret)
}

func TestListWorkers_Empty(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	workers, err := s.ListWorkers(testContext(), &mshipadminpb.ListWorkersRequest{})
	require.Nil(t, err)
	require.Empty(t, workers.Workers)
}

func TestListWorkers_One(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Create(&mothership_db.Worker{
		Name:      "test",
		WorkerID:  "test-id",
		ApiSecret: "secret",
	}))
	defer func() {
		require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	}()

	workers, err := s.ListWorkers(testContext(), &mshipadminpb.ListWorkersRequest{})
	require.Nil(t, err)
	require.Len(t, workers.Workers, 1)
	require.Equal(t, "test", workers.Workers[0].Name)
	require.Equal(t, "test-id", workers.Workers[0].WorkerId)
	require.Empty(t, workers.Workers[0].ApiSecret)
}

func TestCreateWorker(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	worker, err := s.CreateWorker(testContext(), &mshipadminpb.CreateWorkerRequest{
		WorkerId: "test-id",
	})
	require.Nil(t, err)
	require.Equal(t, "test-id", worker.WorkerId)
	require.NotEmpty(t, worker.Name)
	require.NotEmpty(t, worker.ApiSecret)
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
}

func TestCreateWorker_Duplicate(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	_, err := s.CreateWorker(testContext(), &mshipadminpb.CreateWorkerRequest{
		WorkerId: "test-id",
	})
	require.Nil(t, err)
	_, err = s.CreateWorker(testContext(), &mshipadminpb.CreateWorkerRequest{
		WorkerId: "test-id",
	})
	require.NotNil(t, err)
	require.Equal(t, codes.AlreadyExists.String(), status.Code(err).String())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
}

func TestCreateWorker_ShortID(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	_, err := s.CreateWorker(testContext(), &mshipadminpb.CreateWorkerRequest{
		WorkerId: "id",
	})
	require.NotNil(t, err)
	require.Equal(t, codes.InvalidArgument.String(), status.Code(err).String())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
}

func TestDeleteWorker(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	worker, err := s.CreateWorker(testContext(), &mshipadminpb.CreateWorkerRequest{
		WorkerId: "test-id",
	})
	require.Nil(t, err)
	_, err = s.DeleteWorker(testContext(), &mshipadminpb.DeleteWorkerRequest{
		Name: worker.Name,
	})
	require.Nil(t, err)
	_, err = s.GetWorker(testContext(), &mshipadminpb.GetWorkerRequest{
		Name: worker.Name,
	})
	require.NotNil(t, err)
	require.Equal(t, codes.NotFound.String(), status.Code(err).String())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
}

func TestDeleteWorker_NotFound(t *testing.T) {
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
	_, err := s.DeleteWorker(testContext(), &mshipadminpb.DeleteWorkerRequest{
		Name: "test",
	})
	require.NotNil(t, err)
	require.Equal(t, codes.NotFound.String(), status.Code(err).String())
	require.Nil(t, base.Q[mothership_db.Worker](s.db).Delete())
}
