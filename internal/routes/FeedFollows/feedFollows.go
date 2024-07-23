package feedfollows

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tade3910/blog_aggregator/internal/database"
	"github.com/tade3910/blog_aggregator/internal/util"
	"github.com/tade3910/blog_aggregator/middleware"
)

type feedfollows struct {
	db *database.Queries
}

func NewFeedsHandler(db *database.Queries) *feedfollows {
	return &feedfollows{
		db: db,
	}
}

func (handler *feedfollows) handleGet(w http.ResponseWriter, r *http.Request) {
	key, exist := r.Context().Value(middleware.ApiKey).(string)
	if !exist {
		util.RespondWithError(w, 500, "error passing apiKey")
	}
	user, err := util.GetUserByApiKey(key, handler.db, r)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	feedFollows, err := handler.db.GetAllUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, feedFollows)
}

func (handler *feedfollows) handlePost(w http.ResponseWriter, r *http.Request) {
	key, exist := r.Context().Value(middleware.ApiKey).(string)
	if !exist {
		util.RespondWithError(w, 500, "error passing apiKey")
	}
	type reqBody struct {
		Feed_id string
	}
	body := reqBody{}
	err := util.GetBody(r, &body)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	feedId, err := uuid.Parse(body.Feed_id)
	if err != nil {
		util.RespondWithError(w, 500, "invalid feed id provided")
	}
	user, err := util.GetUserByApiKey(key, handler.db, r)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	post, err := handler.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedId,
		UserID:    user.ID,
	})
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 200, post)
}

func (handler *feedfollows) handleDelete(w http.ResponseWriter, r *http.Request, feedFollowId string) {
	uuid, err := uuid.Parse(feedFollowId)
	if err != nil {
		util.RespondWithError(w, 500, "invalid feed id provided")
	}
	err = handler.db.DeleteFeedFollow(r.Context(), uuid)
	if err != nil {
		util.RespondWithError(w, 500, err.Error())
		return
	}
	util.RespondWithJSON(w, 201, nil)
}

func (handler *feedfollows) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url_split := strings.Split(r.URL.Path, "/")
	switch len(url_split) {
	case 2:
		switch r.Method {
		case http.MethodPost:
			handler.handlePost(w, r)
		case http.MethodGet:
			handler.handleGet(w, r)
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	case 3:
		switch r.Method {
		case http.MethodDelete:
			handler.handleDelete(w, r, url_split[2])
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
