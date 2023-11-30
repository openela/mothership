package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"net/http"
)

func (s *server) entriesListView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pt, token, err := getPt(r)
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Get entries from API
	entries, err := s.srpmArchiver.ListEntries(r.Context(), &mothershippb.ListEntriesRequest{
		PageSize:  int32(pt.PageSize),
		PageToken: token,
		Filter:    pt.Filter,
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Render template
	s.renderTemplate(w, "entries", newDataPt[*mothershippb.ListEntriesResponse](r, pt, entries, len(entries.Entries)))
}

func (s *server) entryGetView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get entry from API
	entry, err := s.srpmArchiver.GetEntry(r.Context(), &mothershippb.GetEntryRequest{
		Name: "entries/" + ps.ByName("name"),
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Render template
	s.renderTemplate(w, "entry", newData(r, entry))
}

func (s *server) entryRescue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Rescue entry
	_, err := s.mshipAdmin.RescueEntryImport(requestToCtx(r), &mshipadminpb.RescueEntryImportRequest{
		Name: "entries/" + ps.ByName("name"),
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Redirect to entry
	http.Redirect(w, r, fmt.Sprintf("/entries/%s", ps.ByName("name")), http.StatusFound)
}

func (s *server) entryRetract(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, err := s.mshipAdmin.RetractEntry(requestToCtx(r), &mshipadminpb.RetractEntryRequest{
		Name: "entries/" + ps.ByName("name"),
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Redirect to entry
	http.Redirect(w, r, fmt.Sprintf("/entries/%s", ps.ByName("name")), http.StatusFound)
}
