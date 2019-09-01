package routes

import (
	"bytes"
	"encoding/json"
	"go-twitter-test/container/mock"
	"go-twitter-test/repositories/messages"
	"go-twitter-test/repositories/messages/messagesfakes"
	"go-twitter-test/repositories/tags/tagsfakes"
	"go-twitter-test/repositories/users"
	"go-twitter-test/repositories/users/usersfakes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessagesRouter_CreateMessage(t *testing.T) {
	// Setting up container with mocked repositories
	c := mock.NewMockedContainer()
	tagsRepo := &tagsfakes.FakeRepository{}
	usersRepo := &usersfakes.FakeRepository{}
	messagesRepo := &messagesfakes.FakeRepository{}
	c.TagsRepositoryReturns(tagsRepo)
	c.UsersRepositoryReturns(usersRepo)
	c.MessagesRepositoryReturns(messagesRepo)

	// Setting up mocks
	const mockedTagID int64 = 123
	const mockedUserID int64 = 456
	const mockedMessageID int64 = 789
	usersRepo.GetReturns(&users.User{ID: mockedUserID, Email: "user@email.com"}, nil)
	tagsRepo.PutReturns(mockedTagID, nil)
	messagesRepo.CreateReturns(mockedMessageID, nil)

	// Setting up router and HTTP request
	router := NewRouter(c)
	request, err := http.NewRequest("POST", "/v1/messages", getRequestBody(t, message{
		Text: "A short message",
		Tag:  "my-test",
	}))
	require.Nil(t, err)

	request.Header.Set("X-User-ID", strconv.FormatInt(mockedUserID, 10))
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	// Assertions
	require.Equal(t, http.StatusCreated, responseRecorder.Code)
	require.Equal(t, mockedUserID, usersRepo.GetArgsForCall(0))
	require.Equal(t, "my-test", tagsRepo.PutArgsForCall(0))
	require.EqualValues(t, messages.MessageCreate{
		ID:      0,
		UserID:  mockedUserID,
		TagID:   mockedTagID,
		Message: "A short message",
	}, messagesRepo.CreateArgsForCall(0))
	require.Equal(
		t,
		"/v1/messages/"+strconv.FormatInt(mockedMessageID, 10),
		responseRecorder.Header().Get("Location"),
	)
}

func getRequestBody(t *testing.T, msg message) *bytes.Buffer {
	body, err := json.Marshal(msg)
	require.Nil(t, err)

	return bytes.NewBuffer(body)
}
