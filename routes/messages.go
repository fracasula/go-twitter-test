package routes

import (
	"database/sql"
	"encoding/json"
	"go-twitter-test/repositories/messages"
	"go-twitter-test/repositories/tags"
	"go-twitter-test/repositories/users"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// NewMessagesRouter returns a router with the messages routes attached
func NewMessagesRouter(
	messagesRepository messages.Repository,
	usersRepository users.Repository,
	tagsRepository tags.Repository,
	logger *log.Logger,
) *chi.Mux {
	router := chi.NewRouter()
	msgs := &messagesRouter{
		messagesRepository: messagesRepository,
		usersRepository:    usersRepository,
		tagsRepository:     tagsRepository,
		logger:             logger,
	}

	router.Get("/", msgs.GetMessages)
	router.Post("/", msgs.CreateMessage)

	return router
}

type messagesRouter struct {
	messagesRepository messages.Repository
	usersRepository    users.Repository
	tagsRepository     tags.Repository
	logger             *log.Logger
}

func (mr *messagesRouter) GetMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var err error
	var tagID int64
	if tag := query.Get("tag"); tag != "" {
		if tagID, err = mr.tagsRepository.GetID(tag); err != nil {
			if err == sql.ErrNoRows {
				RenderError(w, r, "Tag not found", http.StatusNotFound)
			} else {
				RenderError(w, r, "Could not get tag ID", http.StatusInternalServerError)
				mr.logger.Printf("Could not get tag ID: %v", err)
			}

			return
		}
	}

	dateStart := query.Get("dateStart")
	dateEnd := query.Get("dateEnd")
	if dateStart != "" || dateEnd != "" {
		if dateStart == "" || dateEnd == "" {
			RenderError(w, r, "dateStart and dateEnd must be used together or not at all", http.StatusBadRequest)
			return
		}
	}

	var unixStart, unixEnd int64
	if dateStart != "" && dateEnd != "" {
		re := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
		if !re.MatchString(dateStart) {
			RenderError(w, r, "Invalid date start YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		if !re.MatchString(dateEnd) {
			RenderError(w, r, "Invalid date end YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		layout := "2006-01-02T15:04:05.000Z"
		timeStart, err := time.Parse(layout, dateStart+"T00:00:00.000Z")
		if err != nil {
			RenderError(w, r, "Cannot parse date start YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		timeEnd, err := time.Parse(layout, dateEnd+"T23:59:59.000Z")
		if err != nil {
			RenderError(w, r, "Cannot parse date end YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		unixStart = timeStart.Unix()
		unixEnd = timeEnd.Unix()
	}

	var responseBody interface{}
	if query.Get("count") == "1" {
		count, err := mr.messagesRepository.CountMessages(tagID, unixStart, unixEnd)
		if err != nil {
			RenderError(w, r, "Could not count messages", http.StatusInternalServerError)
			mr.logger.Printf("Could not count messages: %v", err)
			return
		}

		responseBody = count
	} else {
		list, err := mr.messagesRepository.GetMessages(tagID, unixStart, unixEnd)
		if err != nil {
			RenderError(w, r, "Could not get messages", http.StatusInternalServerError)
			mr.logger.Printf("Could not get messages: %v", err)
			return
		}

		responseBody = list
	}

	render.JSON(w, r, responseBody)
}

func (mr *messagesRouter) CreateMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int64)
	if userID == 0 {
		// if this happens there's most likely an issue with the API Gateway, we should log and
		// return a 500, since there's no API Gateway at the moment I'll just return a 401 for now
		RenderError(w, r, "No user ID provided", http.StatusUnauthorized)
		return
	}

	user, err := mr.usersRepository.Get(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// if this happens there's most likely an issue with the API Gateway, we should log and
			// return a 500 like above
			RenderError(w, r, "User not found", http.StatusForbidden)
		} else {
			RenderError(w, r, "Repository error", http.StatusInternalServerError)
			mr.logger.Printf("Could not get user %q: %v", userID, err)
		}

		return
	}

	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderError(w, r, "Invalid request body", http.StatusBadRequest)
		return
	}

	var body message
	if err := json.Unmarshal(jsonData, &body); err != nil {
		RenderError(w, r, "Request body is not a valid message", http.StatusBadRequest)
		return
	}

	tagID, err := mr.tagsRepository.Put(body.Tag)
	if err != nil {
		RenderError(w, r, "Could not create tag", http.StatusInternalServerError)
		mr.logger.Printf("Could not put tag %q: %v", body.Tag, err)
		return
	}

	msg := messages.MessageCreate{
		TagID:   tagID,
		UserID:  user.ID,
		Message: body.Text,
	}

	msgID, err := mr.messagesRepository.Create(msg)
	if err != nil {
		RenderError(w, r, "Could not create message", http.StatusInternalServerError)
		mr.logger.Printf("Could not create message %+v: %v", msg, err)
		return
	}

	w.Header().Set("Location", "/v1/messages/"+strconv.FormatInt(msgID, 10))
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

type message struct {
	Text string `json:"text"`
	Tag  string `json:"tag"`
}
