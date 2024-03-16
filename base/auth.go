package base

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getGithubReq(url string, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func executeGithubReq(url string, token string) (*http.Response, error) {
	client := &http.Client{}
	req, err := getGithubReq(url, token)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func getGithubUsername(token string) (string, error) {
	resp, err := executeGithubReq("https://api.github.com/user", token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	username, ok := data["login"].(string)
	if !ok {
		return "", errors.New("no username in response")
	}

	return username, nil
}

func getGithubMembership(token string, team string) (bool, error) {
	username, err := getGithubUsername(token)
	if err != nil {
		return false, err
	}

	resp, err := executeGithubReq("https://api.github.com/orgs/"+team+"/memberships/"+username, token)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, err
	}

	state, ok := data["state"].(string)
	if !ok {
		return false, errors.New("no state in response")
	}

	return state == "active", nil
}

// GithubGrpcInterceptor is a gRPC interceptor that checks for a valid GitHub token
// and that the user is in the given team
func GithubGrpcInterceptor(team string) (grpc.UnaryServerInterceptor, error) {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing auth token")
		}

		// verify the token
		isMember, err := getGithubMembership(token, team)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to verify membership")
		}
		if !isMember {
			return nil, status.Error(codes.PermissionDenied, "not in team")
		}

		return handler(ctx, req)
	}

	return interceptor, nil
}
