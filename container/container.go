package container

import (
	"database/sql"
	"fmt"
	"go-twitter-test/repositories/messages"
	"go-twitter-test/repositories/tags"
	"go-twitter-test/repositories/users"
	"go-twitter-test/sqlite"
	"log"
	"os"
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
	db, err := sqlite.New(sqliteDsn)
	if err != nil {
		return nil, fmt.Errorf("container could not initialize db: %v", err)
	}

	return &container{
		db:                 db,
		logger:             log.New(os.Stdout, "", log.LstdFlags),
		messagesRepository: messages.New(db),
		usersRepository:    users.New(db),
		tagsRepository:     tags.New(db),
	}, nil
}
