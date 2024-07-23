package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tade3910/blog_aggregator/internal/database"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}

func GetBody[T interface{}](r *http.Request, bodyStruct *T) error {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return fmt.Errorf("could not read body")
	}
	err = json.Unmarshal(body, bodyStruct)
	if err != nil {
		return fmt.Errorf("invalid body")
	}
	return nil
}

type tokenType string

const (
	ApiKey tokenType = "ApiKey"
)

func GetAuthToken(r *http.Request, tokenType tokenType) (string, error) {
	authorizationString := r.Header.Get("Authorization")
	splits := strings.Split(authorizationString, " ")
	if len(splits) != 2 || splits[0] != string(tokenType) {
		return "", fmt.Errorf("invalid authorization key")
	}
	return splits[1], nil

}

func GetUserByApiKey(apiKey string, db *database.Queries, r *http.Request) (database.User, error) {
	return db.GetUserByApiKey(r.Context(), apiKey)
}
