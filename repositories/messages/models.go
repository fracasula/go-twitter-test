package messages

// MessageCreate is a model used when creating a new message (see POST /v1/messages)
type MessageCreate struct {
	ID      int64
	UserID  int64
	TagID   int64
	Message string
}

// MessageList is used when returning a list of messages (see GET /v1/messages)
type MessageList struct {
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	UserEmail string `json:"user_email"`
	Tag       string `json:"tag"`
}
