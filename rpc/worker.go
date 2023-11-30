package mothership_rpc

import (
	"context"
	"fmt"
	"github.com/openela/mothership/base"
	mothership_db "github.com/openela/mothership/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// getWorkerIdentity returns the identity of the worker that the request is
// coming from. Returns an error if the worker is not found or unauthenticated.
func (s *Server) getWorkerIdentity(ctx context.Context) (*mothership_db.Worker, error) {
	// Get x-mship-worker-secret
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	fmt.Println(md)

	secrets := md["x-mship-worker-secret"]
	if len(secrets) != 1 {
		return nil, status.Error(codes.Unauthenticated, "missing worker secret")
	}

	secret := secrets[0]
	worker, err := base.Q[mothership_db.Worker](s.db).F("api_secret", secret).GetOrNil()
	if err != nil {
		base.LogErrorf("failed to get worker: %v", err)
		return nil, status.Error(codes.Internal, "failed to get worker")
	}

	if worker == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid worker secret")
	}

	return worker, nil
}
