// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	_ "embed"
	"net/http"
)

type FrontendInfo struct {
	// NoRun is a flag to disable running the frontend server
	NoRun bool

	// MuxHandler is the HTTP handler (can be nil)
	MuxHandler http.Handler

	// Title to add to the HTML page
	Title string

	// Port is the port to serve the frontend on
	Port int

	// Self is the URL to the frontend server
	Self string

	// NoAuth is a flag to disable authentication
	NoAuth bool

	// AllowUnauthenticated is a flag to allow unauthenticated users
	AllowUnauthenticated bool

	// OIDCIssuer is the issuer to use for authentication
	OIDCIssuer string

	// OIDCClientID is the client ID to use for authentication
	OIDCClientID string

	// OIDCClientSecret is the client secret to use for authentication
	OIDCClientSecret string

	// OIDCGroup is the group to check for authentication
	OIDCGroup string

	// OIDCUserInfoOverride is a flag to override the userinfo endpoint
	// todo(mustafa): we don't need to use it yet since RESF deploys cluster external Keycloak.
	OIDCUserInfoOverride string
}
