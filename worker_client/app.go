package worker_client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/worker_client/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getRedHatRelease() (string, error) {
	err := exec.Command("dnf", "install", "--refresh", "-y", "redhat-release").Run()
	if err != nil {
		return "", err
	}

	bts, err := os.ReadFile("/etc/redhat-release")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bts)), nil
}

func Run(ctx context.Context, rootURI string, forceRelease string, s state.State, srpmArchiver mothershippb.SrpmArchiverClient) error {
	_, err := srpmArchiver.WorkerPing(ctx, &emptypb.Empty{})
	if err != nil {
		slog.Error("failed to ping mothership", "error", err)
	}

	redHatRelease := forceRelease
	if redHatRelease == "" {
		redHatRelease, err = getRedHatRelease()
		if err != nil {
			return err
		}
	}
	slog.Info("got red hat release", "release", redHatRelease)

	err = s.FetchNewPackageState()
	if err != nil {
		return err
	}

	dirtyObjects := s.GetDirtyObjects()
	if len(dirtyObjects) == 0 {
		return nil
	}

	// Create batch
	batch, err := srpmArchiver.CreateBatch(
		ctx,
		&mothershippb.CreateBatchRequest{
			Batch: &mothershippb.Batch{},
		},
	)
	if err != nil {
		return err
	}
	slog.Info("created batch", "batch", batch.Name)

	var allOperationNames []string
	defer func() {
		slog.Info("sealing batch", "batch", batch.Name)
		_, _ = srpmArchiver.SealBatch(
			ctx,
			&mothershippb.SealBatchRequest{
				Name:           batch.Name,
				OperationNames: allOperationNames,
			},
		)
	}()

	for _, obj := range dirtyObjects {
		entry, err := srpmArchiver.SubmitEntry(
			ctx,
			&mothershippb.SubmitEntryRequest{
				ProcessRpmRequest: &mothershippb.ProcessRPMRequest{
					RpmUri:     fmt.Sprintf("%s/%s", rootURI, strings.TrimPrefix(obj, "/")),
					OsRelease:  redHatRelease,
					Checksum:   strings.TrimPrefix(obj, "/"),
					Repository: "",
					Batch:      batch.Name,
				},
			},
		)
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Code() == codes.AlreadyExists {
					slog.Info("entry already exists", "obj", obj)
					continue
				}
			}
			slog.Error("failed to submit entry", "error", err)
			return err
		}

		allOperationNames = append(allOperationNames, entry.Name)

		slog.Info("submitted entry", "entry", entry.Name)
	}

	err = s.WritePackageState()
	if err != nil {
		return err
	}

	slog.Info("wrote package state")

	return nil
}
