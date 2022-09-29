FROM golang:1.17-alpine as builder
COPY ./ /app/
WORKDIR /app
RUN apk update && \
  apk add ffmpeg && \
  apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  apk del tzdata && \
  go mod download && \
  go build -ldflags "-s -w" -o kenkoukun
CMD ["/app/kenkoukun"]
