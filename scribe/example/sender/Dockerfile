# Multi-stage build: compile a static binary, ship it on a minimal runtime.
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 go build -o /sender .

FROM alpine:3.20
COPY --from=build /sender /sender
ENTRYPOINT ["/sender"]
