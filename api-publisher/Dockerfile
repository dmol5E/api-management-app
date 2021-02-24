FROM golang:1.16.0-alpine as build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/api-publisher .

FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/out/api-publisher /app/api-publisher
RUN chown 10001:10001 /app/api-publisher
RUN chmod 755 /app/api-publisher

# This container exposes port 8080 to the outside world
EXPOSE 8080

USER 10001

# Run the binary program produced by `go install`
CMD ["/app/api-publisher"]