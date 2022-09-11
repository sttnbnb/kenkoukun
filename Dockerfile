FROM golang:1.17-alpine as builder
ENV APPDIR /go/app
COPY go.mod go.sum $APPDIR/
WORKDIR $APPDIR
RUN go mod download
WORKDIR /
COPY ./ $APPDIR/
WORKDIR $APPDIR
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o kenkoukun main.go

FROM gcr.io/distroless/static
COPY --from=builder /go/app/kenkoukun ./app/
CMD ["/app/kenkoukun"]
