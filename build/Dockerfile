# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=1
WORKDIR /build/

RUN apk --no-cache add tzdata gcc musl-dev

COPY ./ ./

RUN \
    go mod tidy && \
    go build -ldflags='-s -w -extldflags "-static"' -o main

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM scratch

WORKDIR /kenkoukun/

COPY ./assets ./assets
COPY ./db ./db
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=builder /build/main ./

CMD ["./main"]
