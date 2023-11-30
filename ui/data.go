package main

import (
	"github.com/gorilla/csrf"
	"github.com/openela/mothership/base"
	"html/template"
	"net/http"
	"strconv"
)

type getPageToken interface {
	GetNextPageToken() string
}

type data[T any] struct {
	Consistent *consistentContext
	Custom     T
}

type alert struct {
	Variant  string
	Title    string
	Subtitle string
	Icon     string
	Closable bool
}

type consistentContext struct {
	Request    *http.Request
	User       base.UserInfo
	Pagination *pageToken
	Alerts     []alert
	Csrf       template.HTML
}

func newData[T any](r *http.Request, custom T, alerts ...alert) *data[T] {
	c := &consistentContext{
		Request: r,
		Alerts:  alerts,
		Csrf:    csrf.TemplateField(r),
	}

	user, ok := r.Context().Value(base.UserContextKey).(base.UserInfo)
	if ok {
		c.User = user
	}

	return &data[T]{
		Consistent: c,
		Custom:     custom,
	}
}

func newDataPt[T getPageToken](r *http.Request, pt *pageToken, custom T, lenEntries int, alerts ...alert) *data[T] {
	ret := newData[T](r, custom, alerts...)

	existingPt := custom.GetNextPageToken()
	if existingPt == "" {
		pt.NextQuery = ""
	}
	if pt.PrevOffset == -1 {
		pt.PrevQuery = ""
	}

	// Set range, start from offset, then add page size (or just amount of entries if less than page size)
	pt.Range = strconv.Itoa(pt.Offset+1) + "-" + strconv.Itoa(pt.Offset+lenEntries)

	ret.Consistent.Pagination = pt

	return ret
}
