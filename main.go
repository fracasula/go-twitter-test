package main

import (
	"go-twitter-test/container"
	"go-twitter-test/routes"
	"log"
	"net/http"
	"os"
	"regexp"
)

const sqliteDsn = "db.sqlite" // @TODO make this one come from an environment variable

func main() {
	// We could use a library to unmarshal env vars (e.g. Netflix/go-env) but since I'm doing custom validation on these
	// and I only have two variables I decided to do it manually.
	// If you end up having more than two environment variables may be worth using a library.
	httpPortRe := "[0-9]{2,5}"
	httpPort := os.Getenv("HTTP_PORT")
	if !regexp.MustCompile(httpPortRe).MatchString(httpPort) {
		log.Fatalf("Invalid HTTP port supplied %q (%v)", httpPort, httpPortRe)
	}

	c, err := container.NewContainer(sqliteDsn)
	if err != nil {
		log.Fatalf("Could not initialize container: %v", err)
	}

	router := routes.NewRouter(c)

	log.Fatal(http.ListenAndServe(":"+httpPort, router))
}
