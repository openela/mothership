package github_bugtracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/openela/mothership/base/bugtracker"
	"github.com/openela/mothership/base/forge"
	github_forge "github.com/openela/mothership/base/forge/github"
)

type Bugtracker struct {
	repo          string
	appId         string
	appPrivateKey []byte
}

func New(repo string, appId string, appPrivateKey []byte) (*Bugtracker, error) {
	_, err := jwt.ParseRSAPrivateKeyFromPEM(appPrivateKey)
	if err != nil {
		return nil, err
	}

	return &Bugtracker{
		repo:          repo,
		appId:         appId,
		appPrivateKey: appPrivateKey,
	}, nil
}

// CreateTicket creates a ticket in the bug tracker.
// Returns the ticket ID or an error.
func (b *Bugtracker) CreateTicket(auth *forge.Authenticator, title string, body string, opts bugtracker.Options) (string, error) {
	// Cast AuthMethod to BasicAuth
	basicAuth := auth.AuthMethod.(*transport_http.BasicAuth)
	token := basicAuth.Password

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	mapBody := map[string]any{
		"title": title,
		"body":  body,
	}
	if opts.Labels != nil && len(opts.Labels) > 0 {
		mapBody["labels"] = opts.Labels
	}
	reqBody, err := json.Marshal(mapBody)
	if err != nil {
		return "", err
	}

	endpoint := "https://api.github.com/" + filepath.Join("repos", b.repo, "issues")
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to create ticket: %s", resp.Status)
	}

	respBody := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	idInt, ok := respBody["number"].(float64)
	if !ok {
		return "", fmt.Errorf("number not found in response")
	}
	id := strconv.FormatFloat(idInt, 'f', 0, 64)

	return id, nil
}

// EditTicket edits the ticket in the bug tracker.
// Returns an error if the ticket could not be edited.
func (b *Bugtracker) EditTicket(auth *forge.Authenticator, ticketID string, title string, body string, opts bugtracker.Options) error {
	// Cast AuthMethod to BasicAuth
	basicAuth := auth.AuthMethod.(*transport_http.BasicAuth)
	token := basicAuth.Password

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	mapBody := map[string]any{
		"title": title,
		"body":  body,
	}
	if opts.Labels != nil && len(opts.Labels) > 0 {
		mapBody["labels"] = opts.Labels
	}
	reqBody, err := json.Marshal(mapBody)
	if err != nil {
		return err
	}

	endpoint := "https://api.github.com/" + filepath.Join("repos", b.repo, "issues", ticketID)
	req, err := http.NewRequest("PATCH", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to close ticket: %s", resp.Status)
	}

	return nil
}

// CloseTicket closes the ticket in the bug tracker.
// Returns an error if the ticket could not be closed.
func (b *Bugtracker) CloseTicket(auth *forge.Authenticator, ticketID string) error {
	// Cast AuthMethod to BasicAuth
	basicAuth := auth.AuthMethod.(*transport_http.BasicAuth)
	token := basicAuth.Password

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	mapBody := map[string]any{
		"state":        "closed",
		"state_reason": "completed",
	}
	body, err := json.Marshal(mapBody)
	if err != nil {
		return err
	}

	endpoint := "https://api.github.com/" + filepath.Join("repos", b.repo, "issues", ticketID)
	req, err := http.NewRequest("PATCH", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to close ticket: %s", resp.Status)
	}

	return nil
}

// TicketURI returns the URI to the ticket in the bug tracker.
// Returns an error if the URI could not be generated.
func (b *Bugtracker) TicketURI(ticketID string) (string, error) {
	return fmt.Sprintf("https://github.com/%s/issues/%s", b.repo, ticketID), nil
}

// URIToTicket returns the ticket ID from the URI.
// Returns an error if the ticket ID could not be extracted.
func (b *Bugtracker) URIToTicket(uri string) (string, error) {
	uri = strings.TrimPrefix(uri, "https://")
	id := filepath.Base(uri)
	if id == "" {
		return "", fmt.Errorf("invalid URI: %s", uri)
	}

	return id, nil
}

// GetAuthenticator returns an authenticator for the bug tracker.
func (b *Bugtracker) GetAuthenticator() (*forge.Authenticator, error) {
	orgSplit := strings.Split(b.repo, "/")
	if len(orgSplit) != 2 {
		return nil, fmt.Errorf("invalid repo name: %s", b.repo)
	}
	org := orgSplit[0]
	ghForge, err := github_forge.New(org, b.appId, b.appPrivateKey, false)
	if err != nil {
		return nil, err
	}
	return ghForge.GetAuthenticator()
}
