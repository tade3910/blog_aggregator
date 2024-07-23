package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tade3910/blog_aggregator/internal/util"
)

type V1 int

func splitURL(url string, baseUrl string) (string, error) {
	splits := strings.Split(url, baseUrl)
	if len(splits) != 2 {
		return "", fmt.Errorf("invalid path")
	}
	return splits[1], nil
}

func (v1 V1) handleHealth(w http.ResponseWriter) {
	util.RespondWithJSON(w, 200, map[string]string{"status": "ok"})
}

func (v1 V1) handleErr(w http.ResponseWriter) {
	util.RespondWithError(w, 500, "Internal Server Error")
}

func (v1 V1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := splitURL(r.URL.Path, "/v1/")
	if err != nil {
		util.RespondWithError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}
	if r.Method != http.MethodGet {
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}
	switch path {
	case "healthz":
		v1.handleHealth(w)
	case "err":
		v1.handleErr(w)
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
