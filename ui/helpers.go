package main

import (
	"context"
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/openela/mothership/base"
	"google.golang.org/grpc/metadata"
	"net/http"
	"path/filepath"
)

func requestToCtx(r *http.Request) context.Context {
	ctx := metadata.NewOutgoingContext(r.Context(), metadata.Pairs())

	// Check if token is in context, if so create proper outgoing context
	if token, ok := r.Context().Value(base.TokenContextKey).(string); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}

	return ctx
}

func redirectTo(path string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, path, http.StatusFound)
	}
}

func serveDir(fs *embed.FS, addPrefix string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filePath := ps.ByName("filepath")
		// Rewrite path to match filePath
		r.URL.Path = filepath.Join(addPrefix, filePath)

		if fs == nil {
			http.ServeFile(w, r, r.URL.Path)
		} else {
			http.FileServer(http.FS(fs)).ServeHTTP(w, r)
		}
	}
}

func (s *server) renderView(templateName string, alerts ...alert) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.renderTemplate(w, templateName, newData[any](r, nil, alerts...))
	}
}

func (s *server) errorView(err string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.renderTemplate(w, "error", newData[string](r, err))
	}
}

func (s *server) ensureLoggedIn(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		_, ok := r.Context().Value(base.TokenContextKey).(string)
		if !ok {
			s.renderTemplate(w, "error", newData[string](r, "Please sign in to access this page."))
			return
		}

		next(w, r, ps)
	}
}
