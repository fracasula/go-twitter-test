package container

import (
	"database/sql"
	"fmt"
	"go-twitter-test/repositories/messages"
	"go-twitter-test/repositories/tags"
	"go-twitter-test/repositories/users"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:generate counterfeiter . Container
type Container interface {
	MessagesRepository() messages.Repository
	UsersRepository() users.Repository
	TagsRepository() tags.Repository
	Logger() *log.Logger
}

type container struct {
	db                 *sql.DB
	logger             *log.Logger
	messagesRepository messages.Repository
	usersRepository    users.Repository
	tagsRepository     tags.Repository
}

func (c *container) MessagesRepository() messages.Repository {
	return c.messagesRepository
}

func (c *container) UsersRepository() users.Repository {
	return c.usersRepository
}

func (c *container) TagsRepository() tags.Repository {
	return c.tagsRepository
}

func (c *container) Logger() *log.Logger {
	return c.logger
}

func NewContainer(sqliteDsn string) (Container, error) {
	db, err := sql.Open("sqlite3", sqliteDsn)
	if err != nil {
		return nil, fmt.Errorf("could not open a connection to %q: %v", sqliteDsn, err)
	}

	db.SetMaxOpenConns(1)

	return &container{
		db:                 db,
		logger:             log.New(os.Stdout, "", log.LstdFlags),
		messagesRepository: messages.New(db),
		usersRepository:    users.New(db),
		tagsRepository:     tags.New(db),
	}, nil
}
