FROM golang:1.17-alpine as builder
COPY ./ /app/
WORKDIR /app
RUN apk update \
  && apk add ffmpeg \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o kenkoukun main.go
CMD ["/app/kenkoukun"]
