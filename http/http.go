package http

import (
	"net/http"
)

const (
	GET  = "GET"
	POST = "POST"
	HEAD = "HEAD"
)

type (
	Request        = http.Request
	Response       = http.Response
	ResponseWriter = http.ResponseWriter
)
