package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// sql.DB manages its own pool of connections so if we were to use the same DB across multiple go
// routines there would be no isolation (SQLite3) unless we set MaxOpenConns to 1.
func New(sqliteDsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", sqliteDsn)
	if err != nil {
		return nil, fmt.Errorf("could not open a connection to %q: %v", sqliteDsn, err)
	}

	db.SetMaxOpenConns(1)

	return db, nil
}

const tagsTable = `CREATE TABLE "tags" (
	id	INTEGER NOT NULL PRIMARY KEY,
	tag	TEXT UNIQUE
)`

const usersTable = `CREATE TABLE "users" (
	id	INTEGER NOT NULL,
	email	TEXT UNIQUE,
	PRIMARY KEY(id)
)`

const messagesTable = `CREATE TABLE "messages" (
	id	INTEGER NOT NULL,
	user_id	INTEGER NOT NULL,
	message	TEXT NOT NULL,
	created_at	INTEGER NOT NULL,
	PRIMARY KEY(id)
)`

const messagesIndexes = `CREATE INDEX user_id_idx ON messages (user_id)`

const messageTagTable = `CREATE TABLE message_tag (
	message_id	INTEGER NOT NULL,
	tag_id	INTEGER NOT NULL
)`

const messageTagIndex1 = `CREATE UNIQUE INDEX message_tag_unique ON message_tag (
	message_id,
	tag_id
)`

const messageTagIndex2 = `CREATE INDEX message_tag_message_id ON message_tag (
	message_id
)`

const messageTagIndex3 = `CREATE INDEX message_tag_tag_id ON message_tag (
	tag_id
)`

func LoadSchema(db *sql.DB) error {
	if _, err := db.Exec(tagsTable); err != nil {
		return fmt.Errorf("could not create tags table: %v", err)
	}
	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("could not create users table: %v", err)
	}
	if _, err := db.Exec(messagesTable); err != nil {
		return fmt.Errorf("could not create messages table: %v", err)
	}
	if _, err := db.Exec(messagesIndexes); err != nil {
		return fmt.Errorf("could not create messages indexes: %v", err)
	}
	if _, err := db.Exec(messageTagTable); err != nil {
		return fmt.Errorf("could not create message_tag table: %v", err)
	}
	if _, err := db.Exec(messageTagIndex1); err != nil {
		return fmt.Errorf("could not create message_tag index 1: %v", err)
	}
	if _, err := db.Exec(messageTagIndex2); err != nil {
		return fmt.Errorf("could not create message_tag index 2: %v", err)
	}
	if _, err := db.Exec(messageTagIndex3); err != nil {
		return fmt.Errorf("could not create message_tag index 3: %v", err)
	}

	return nil
}
