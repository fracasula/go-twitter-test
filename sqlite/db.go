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

func LoadSchema(db *sql.DB) error {
	_, err := db.Exec(tagsTable)
	if err != nil {
		return fmt.Errorf("could not create tags table: %v", err)
	}

	return nil
}
