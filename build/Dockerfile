# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=1

WORKDIR /build/
COPY ./ ./

RUN apk --no-cache add tzdata gcc musl-dev && \
    go mod download && \
    go install -ldflags='-s -w -extldflags "-static"' ./main.go

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM scratch

WORKDIR /kenkoukun/

COPY ./assets ./assets
COPY ./db ./db
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=builder /go/bin/main ./

CMD ["./main"]
