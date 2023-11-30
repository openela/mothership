package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type pageToken struct {
	Offset     int    `json:"offset"`
	Filter     string `json:"filter"`
	OrderBy    string `json:"order_by"`
	PageSize   int    `json:"page_size"`
	PrevOffset int    `json:"prev_offset"`
	NextOffset int    `json:"next_offset"`
	NextQuery  string `json:"next_query"`
	PrevQuery  string `json:"prev_query"`
	Range      string `json:"range"`

	GetSizeQuery func(size int) string `json:"-"`
}

func (p *pageToken) ToB64() (string, error) {
	x, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	// Convert to base64
	y := base64.StdEncoding.EncodeToString(x)

	return y, nil
}

func getPt(r *http.Request) (*pageToken, string, error) {
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, "", err
	}

	q := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")

	pageSize := r.URL.Query().Get("size")
	if pageSize == "" {
		pageSize = "100"
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, "", err
	}

	pt := &pageToken{
		Offset:     offsetInt,
		Filter:     q,
		OrderBy:    order,
		PageSize:   pageSizeInt,
		PrevOffset: offsetInt - pageSizeInt,
		NextOffset: offsetInt + pageSizeInt,
	}
	if pt.PrevOffset < 0 {
		pt.PrevOffset = 0
	}
	if pt.Offset == pt.PrevOffset {
		pt.PrevOffset = -1
	}

	var token string
	if offsetInt > 0 {
		token, err = pt.ToB64()
		if err != nil {
			return nil, "", err
		}
	}

	// Set query strings
	nextQuery := url.Values{}
	prevQuery := url.Values{}
	if q != "" {
		nextQuery.Set("q", q)
		prevQuery.Set("q", q)
	}
	if order != "" {
		nextQuery.Set("order", order)
		prevQuery.Set("order", order)
	}
	nextQuery.Set("offset", strconv.Itoa(pt.NextOffset))
	prevQuery.Set("offset", strconv.Itoa(pt.PrevOffset))
	nextQuery.Set("size", strconv.Itoa(pt.PageSize))
	prevQuery.Set("size", strconv.Itoa(pt.PageSize))

	pt.NextQuery = "?" + nextQuery.Encode()
	pt.PrevQuery = "?" + prevQuery.Encode()

	pt.GetSizeQuery = func(size int) string {
		currentQuery := r.URL.Query()
		currentQuery.Set("size", strconv.Itoa(size))
		return "?" + currentQuery.Encode()
	}

	return pt, token, nil
}
