package middleware

import (
	"context"
	"net/http"

	"github.com/tade3910/blog_aggregator/internal/util"
)

type middleWare struct{}

func NewMiddleWare() *middleWare {
	return &middleWare{}
}

type contextKey string

const (
	ApiKey contextKey = "apiKey"
)

func (middleware *middleWare) EnsureAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ignoreRoute(r) {
			next.ServeHTTP(w, r)
			return
		}
		apiKey, err := util.GetAuthToken(r, util.ApiKey)
		if err != nil {
			util.RespondWithError(w, 401, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), ApiKey, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ignoreRoute(r *http.Request) bool {
	url := r.URL.Path
	switch r.Method {
	case http.MethodPost:
		if url == "/users" {
			return true
		}
	case http.MethodGet:
		if url == "/feeds" {
			return true
		}
	}
	return false
}
