package users

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tade3910/blog_aggregator/internal/database"
	"github.com/tade3910/blog_aggregator/internal/util"
	"github.com/tade3910/blog_aggregator/middleware"
)

type users struct {
	db *database.Queries
}

func NewUsersHandler(db *database.Queries) *users {
	return &users{
		db: db,
	}
}

func (handler *users) getUser(w http.ResponseWriter, r *http.Request) {
	key, exist := r.Context().Value(middleware.ApiKey).(string)
	if !exist {
		util.RespondWithError(w, 500, "error passing apiKey")
	}
	user, err := util.GetUserByApiKey(key, handler.db, r)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, user)

}

func (handler *users) postUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name string
	}
	body := reqBody{}
	err := util.GetBody(r, &body)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	user, err := handler.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
	})
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, user)
}

func (handler *users) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.postUser(w, r)
	case http.MethodGet:
		handler.getUser(w, r)
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
