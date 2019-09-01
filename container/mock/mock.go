package mock

import (
	"go-twitter-test/container/containerfakes"
	"io/ioutil"
	"log"
)

func NewMockedContainer() *containerfakes.FakeContainer {
	c := &containerfakes.FakeContainer{}
	c.LoggerReturns(nullLogger())
	return c
}

func nullLogger() *log.Logger {
	return log.New(ioutil.Discard, "", log.LstdFlags)
}
