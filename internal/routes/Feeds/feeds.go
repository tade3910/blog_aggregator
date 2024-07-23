package feeds

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tade3910/blog_aggregator/internal/database"
	"github.com/tade3910/blog_aggregator/internal/util"
	"github.com/tade3910/blog_aggregator/middleware"
)

type feeds struct {
	db *database.Queries
}

func NewFeedsHandler(db *database.Queries) *feeds {
	return &feeds{
		db: db,
	}
}

func (handler *feeds) handleGet(w http.ResponseWriter, r *http.Request) {
	posts, err := handler.db.GetAllFeeds(r.Context())
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, posts)
}

func (handler *feeds) handlePost(w http.ResponseWriter, r *http.Request) {
	key, exist := r.Context().Value(middleware.ApiKey).(string)
	if !exist {
		util.RespondWithError(w, 500, "error passing apiKey")
	}
	type reqBody struct {
		Name string
		Url  string
	}
	body := reqBody{}
	err := util.GetBody(r, &body)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	user, err := util.GetUserByApiKey(key, handler.db, r)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	post, err := handler.db.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      body.Name,
		Url:       body.Url,
		UserID:    user.ID,
	})
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, post)
}

func (handler *feeds) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.handlePost(w, r)
	case http.MethodGet:
		handler.handleGet(w, r)
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
