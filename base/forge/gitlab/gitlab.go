// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	transport_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/openela/mothership/base/forge"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Forge struct {
	host                 string
	group                string
	username             string
	password             string
	authorName           string
	authorEmail          string
	shouldMakeRepoPublic bool
}

func New(host string, group string, username string, password string, authorName string, authorEmail string, shouldMakeRepoPublic bool) *Forge {
	return &Forge{
		host:                 host,
		group:                group,
		username:             username,
		password:             password,
		authorName:           authorName,
		authorEmail:          authorEmail,
		shouldMakeRepoPublic: shouldMakeRepoPublic,
	}
}

func (f *Forge) GetAuthenticator() (*forge.Authenticator, error) {
	transporter := &transport_http.BasicAuth{
		Username: f.username,
		Password: f.password,
	}

	// We're assuming never expiring tokens for now
	// Set it to 100 years from now
	expires := time.Now().AddDate(100, 0, 0)

	return &forge.Authenticator{
		AuthMethod:  transporter,
		AuthorName:  f.authorName,
		AuthorEmail: f.authorEmail,
		Expires:     expires,
	}, nil
}

func (f *Forge) GetRemote(repo string) string {
	return fmt.Sprintf("https://%s/%s/%s", f.host, f.group, repo)
}

func (f *Forge) GetCommitViewerURL(repo string, commit string) string {
	return fmt.Sprintf(
		"https://%s/%s/%s/-/commit/%s",
		f.host,
		f.group,
		repo,
		commit,
	)
}

func (f *Forge) EnsureRepositoryExists(auth *forge.Authenticator, repo string) error {
	// Cast AuthMethod to BasicAuth
	basicAuth := auth.AuthMethod.(*transport_http.BasicAuth)
	token := basicAuth.Password

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Check if the repo exists
	urlEncodedPath := url.PathEscape(fmt.Sprintf("%s/%s", f.group, repo))
	endpoint := fmt.Sprintf("https://%s/api/v4/projects/%s", f.host, urlEncodedPath)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		// Repo exists, we're done
		return nil
	}

	// Repo doesn't exist, create it
	// First get namespace id
	endpoint = fmt.Sprintf("https://%s/api/v4/namespaces/%s", f.host, url.PathEscape(f.group))
	req, err = http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("namespace %s does not exist", f.group)
	}

	mapBody := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&mapBody)
	if err != nil {
		return err
	}

	namespaceId := mapBody["id"].(float64)

	mapBody = map[string]any{
		"name":         repo,
		"namespace_id": namespaceId,
	}
	if f.shouldMakeRepoPublic {
		mapBody["visibility"] = "public"
	} else {
		mapBody["visibility"] = "private"
	}

	endpoint = fmt.Sprintf("https://%s/api/v4/projects", f.host)
	body, err := json.Marshal(mapBody)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create repo %s: %s", repo, string(body))
	}

	return nil
}

func (f *Forge) WithNamespace(namespace string) forge.Forge {
	newF := *f
	newF.group = namespace
	return &newF
}
