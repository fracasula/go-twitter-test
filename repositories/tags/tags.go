package tags

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	Put(tag string) (int64, error)
	GetID(tag string) (int64, error)
}

type tagsRepository struct {
	db *sql.DB
}

func (r *tagsRepository) GetID(tag string) (int64, error) {
	tag = r.tokenize(tag)

	var tagID int64
	err := r.db.QueryRow("SELECT id FROM tags WHERE tag = ?", tag).Scan(&tagID)
	if err != nil {
		return 0, err
	}

	return tagID, nil
}

func (r *tagsRepository) Put(tag string) (int64, error) {
	tag = r.tokenize(tag)

	res, err := r.db.Exec("INSERT OR IGNORE INTO tags (tag) VALUES (?)", tag)
	if err != nil {
		return 0, fmt.Errorf("could not insert tag %q: %v", tag, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last inserted id when creating tag %q: %v", tag, err)
	}

	if id > 0 {
		// the tag didn't exist, returning the ID
		return id, nil
	}

	// the tag already existed, let's select its ID
	if err = r.db.QueryRow("SELECT id FROM tags WHERE tag = ?", tag).Scan(&id); err != nil {
		return 0, fmt.Errorf("could not get id for tag %q: %v", tag, err)
	}

	return id, nil
}

func (r *tagsRepository) tokenize(tag string) string {
	// let's tokenize the tag to avoid duplicates as much as possible
	// this could potentially be more complicated but for now I'll
	// just make it lowercase and I'll replace spaces with dashes
	tag = strings.ToLower(tag)
	tag = strings.Trim(tag, " ")
	tag = strings.ReplaceAll(tag, " ", "-")
	return tag
}

func New(db *sql.DB) Repository {
	return &tagsRepository{
		db: db,
	}
}
