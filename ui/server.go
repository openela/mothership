package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
	"github.com/openela/mothership/base"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"golang.org/x/oauth2"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

//go:embed node_modules/@shoelace-style/shoelace/dist/assets
var shoelace embed.FS

//go:embed static/*
var static embed.FS

type server struct {
	srpmArchiver mothershippb.SrpmArchiverClient
	mshipAdmin   mshipadminpb.MshipAdminClient

	frontendInfo   *base.FrontendInfo
	oauth2Config   *oauth2.Config
	oidcProvider   base.OidcProvider
	templateBundle map[string]*template.Template
}

func (s *server) routes() http.Handler {
	router := httprouter.New()
	router.GET("/_/healthz", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	router.GET("/vendor/shoelace/assets/*filepath", serveDir(&shoelace, "node_modules/@shoelace-style/shoelace/dist/assets"))
	if version != "DEV" {
		router.GET("/static/*filepath", serveDir(&static, "static"))
	} else {
		router.GET("/static/*filepath", serveDir(nil, "static"))
	}

	// Auth
	router.GET("/auth/login", s.authLoginHandler)
	router.GET("/auth/callback", s.authCallbackHandler)
	router.GET("/auth/logout", s.authLogoutHandler)

	// Entries
	router.GET("/", redirectTo("/entries"))
	router.GET("/entries", s.entriesListView)
	router.GET("/entries/:name", s.entryGetView)
	router.POST("/entries/:name/rescue", s.ensureLoggedIn(s.entryRescue))
	router.POST("/entries/:name/retract", s.ensureLoggedIn(s.entryRetract))

	// Workers
	router.GET("/workers", s.ensureLoggedIn(s.workersListView))
	router.POST("/workers", s.ensureLoggedIn(s.workerCreate))
	router.GET("/workers/:name", s.ensureLoggedIn(s.workerGetView))

	// Not found
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.renderView("404")(w, r, nil)
	})

	return csrf.Protect([]byte(os.Getenv("CSRF_SECRET")))(s.authMiddleware(router))
}

func (s *server) run() error {
	if err := s.loadTemplates(); err != nil {
		return err
	}

	ctx := context.TODO()
	provider2, err := oidc.NewProvider(ctx, s.frontendInfo.OIDCIssuer)
	if err != nil {
		return fmt.Errorf("failed to create oidc provider: %w", err)
	}

	provider := &base.OidcProviderImpl{Provider: provider2}

	redirectURL := s.frontendInfo.Self + "/auth/callback"

	oauth2Config := oauth2.Config{
		ClientID:     s.frontendInfo.OIDCClientID,
		ClientSecret: s.frontendInfo.OIDCClientSecret,
		Endpoint:     provider2.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}
	s.oauth2Config = &oauth2Config
	s.oidcProvider = provider

	router := s.routes()

	slog.Info("Starting server", "port", s.frontendInfo.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.frontendInfo.Port), router)
}
