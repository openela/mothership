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

func Run(ctx context.Context, rootURI string, s state.State, srpmArchiver mothershippb.SrpmArchiverClient) error {
	redHatRelease, err := getRedHatRelease()
	if err != nil {
		return err
	}

	err = s.FetchNewPackageState()
	if err != nil {
		return err
	}

	dirtyObjects := s.GetDirtyObjects()
	if len(dirtyObjects) == 0 {
		return nil
	}

	for _, obj := range dirtyObjects {
		entry, err := srpmArchiver.SubmitEntry(
			ctx,
			&mothershippb.SubmitEntryRequest{
				ProcessRpmRequest: &mothershippb.ProcessRPMRequest{
					RpmUri:     fmt.Sprintf("%s/%s", rootURI, strings.TrimPrefix(obj, "/")),
					OsRelease:  redHatRelease,
					Checksum:   strings.TrimPrefix(obj, "/"),
					Repository: "",
					Batch:      "",
				},
			},
		)
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Code() == codes.AlreadyExists {
					slog.Info("entry already exists", "entry", entry.Name)
					continue
				}
			}
			slog.Error("failed to submit entry", "error", err)
			return err
		}

		slog.Info("submitted entry", "entry", entry.Name)
	}

	err = s.WritePackageState()
	if err != nil {
		return err
	}

	return nil
}
