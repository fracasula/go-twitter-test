###############
# FIRST STAGE #
###############
FROM golang:1.12-alpine as builder

# Installing dependencies
RUN apk add git gcc g++ libc-dev musl-dev sqlite --update
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN go get -u github.com/maxbrunsfeld/counterfeiter

# Bootstrapping modules dependencies
RUN mkdir -p /src/go-twitter-test
WORKDIR /src/go-twitter-test
COPY go.mod go.mod
COPY go.sum go.sum
RUN go get -d

# Copying source files after `go get` to retain modules cache as often as possible
COPY . /src/go-twitter-test

# Running tests
RUN go generate ./...
RUN go test ./...

# Compiling binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o api


################
# SECOND STAGE #
################
FROM scratch

ARG HTTP_PORT
ENV HTTP_PORT ${HTTP_PORT}
EXPOSE ${HTTP_PORT}

# Copy api binary from first step
COPY --from=builder /src/go-twitter-test/api api

CMD ["./api"]
