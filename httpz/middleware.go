package httpz

import (
	"context"
	"net/http"
)

type DoFunc func(ctx context.Context, req *http.Request) (*http.Response, error)

type Middleware func(DoFunc) DoFunc
