package helpers

import (
	"net/url"
	"strings"
)

// ParseArgsGet parses a map for GET requests and returns a query string
func ParseArgsGet(query map[string]string) string {
	qSlice := []string{}
	for k, v := range query {
		qSlice = append(qSlice, k+"="+v)
	}
	queryString := strings.Join(qSlice, "&")
	return queryString
}

// ParseArgsBody parses a map for POST/PUT/DELETE requests and returns a request body
func ParseArgsBody(query map[string]string) url.Values {
	reqBody := url.Values{}
	for k, v := range query {
		reqBody.Add(k, v)
	}
	return reqBody
}
