package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	"net/http"
)

func (s *server) renderWorkerList(w http.ResponseWriter, r *http.Request, alerts ...alert) {
	pt, token, err := getPt(r)
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Get workers from API
	workers, err := s.mshipAdmin.ListWorkers(requestToCtx(r), &mshipadminpb.ListWorkersRequest{
		PageSize:  int32(pt.PageSize),
		PageToken: token,
		Filter:    pt.Filter,
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Render template
	s.renderTemplate(w, "workers", newDataPt[*mshipadminpb.ListWorkersResponse](r, pt, workers, len(workers.Workers), alerts...))
}

func (s *server) workersListView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.renderWorkerList(w, r)
}

func (s *server) workerGetView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get worker from API
	worker, err := s.mshipAdmin.GetWorker(requestToCtx(r), &mshipadminpb.GetWorkerRequest{
		Name: "workers/" + ps.ByName("name"),
	})
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Render template
	s.renderTemplate(w, "worker", newData(r, worker))
}

func (s *server) workerCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		s.errorView(err.Error())(w, r, nil)
		return
	}

	// Check if worker_id is set
	if r.Form.Get("worker_id") == "" {
		s.renderWorkerList(w, r, alert{
			Variant: "danger",
			Title:   "You must specify a worker ID",
		})
		return
	}

	// Create worker
	worker, err := s.mshipAdmin.CreateWorker(requestToCtx(r), &mshipadminpb.CreateWorkerRequest{
		WorkerId: r.Form.Get("worker_id"),
	})
	if err != nil {
		s.renderWorkerList(w, r, alert{
			Variant:  "danger",
			Title:    "Failed to create worker",
			Subtitle: err.Error(),
		})
		return
	}

	// Render template
	s.renderWorkerList(w, r, alert{
		Variant:  "success",
		Title:    "Worker created",
		Subtitle: fmt.Sprintf("%s is the secret for this worker. This is the only time it will be shown, so make sure to save it somewhere safe.", worker.ApiSecret),
	})
}
