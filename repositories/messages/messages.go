package messages

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository represents a contract for querying messages from an arbitrary data source
type Repository interface {
	Create(msg *MessageCreate) error
	GetMessages(tagID, dateStart, dateEnd int64) ([]MessageList, error)
}

type messagesRepository struct {
	db *sql.DB
}

func (r *messagesRepository) Create(msg *MessageCreate) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction for creating a new message: %v", err)
	}

	res, err := tx.Exec(
		"INSERT INTO messages (user_id, message, created_at) VALUES (?, ?, ?)",
		msg.UserID, msg.Message, time.Now().Unix(), // https://www.sqlite.org/datatype3.html#datetime
	)
	if err != nil {
		return fmt.Errorf("could not create message with user ID %d and message %q: %v", msg.UserID, msg.Message, err)
	}

	msgID, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("could not get last inserted message ID: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("could not rollback transaction (after %v): %v", err, rollbackErr)
		}

		return err
	}

	_, err = tx.Exec("INSERT INTO message_tag (message_id, tag_id) VALUES (?, ?)", msgID, msg.TagID)
	if err != nil {
		err = fmt.Errorf("could not link tag to message: %v", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = fmt.Errorf("could not rollback transaction (after %v): %v", err, rollbackErr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction while creating new message: %v", err)
	}

	msg.ID = msgID

	return nil
}

func (r *messagesRepository) GetMessages(tagID, dateStart, dateEnd int64) ([]MessageList, error) {
	// this query without pagination could be dangerous, avoid in production
	var args []interface{}
	query := `SELECT m.id, m.message, m.created_at, u.email, t.tag
		FROM messages AS m
		INNER JOIN users AS u ON m.user_id = u.id
		INNER JOIN message_tag AS mt ON m.id = mt.message_id
		INNER JOIN tags AS t ON mt.tag_id = t.id
		WHERE 1 = 1`
	if tagID != 0 {
		query += " AND t.id = ?"
		args = append(args, tagID)
	}
	if dateStart != 0 && dateEnd != 0 {
		query += " AND m.created_at BETWEEN ? AND ?"
		args = append(args, dateStart, dateEnd)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not get messages: %v", err)
	}

	var list []MessageList
	for rows.Next() {
		msg := MessageList{}
		var createdAt int64
		if err := rows.Scan(&msg.ID, &msg.Message, &createdAt, &msg.UserEmail, &msg.Tag); err != nil {
			return nil, fmt.Errorf("could not scan message row: %v", err)
		}

		msg.CreatedAt = time.Unix(createdAt, 0).Format("2006-01-02T15:04:05")

		list = append(list, msg)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("could not close rows: %v", err)
	}

	return list, nil
}

func New(db *sql.DB) Repository {
	return &messagesRepository{
		db: db,
	}
}