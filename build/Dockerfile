FROM golang:1.19.5-alpine AS builder
WORKDIR /build/
COPY ./ ./
RUN apk --no-cache add tzdata && \
  go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -o kenkoukun

FROM scratch
WORKDIR /root/
COPY ./assets ./assets
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=builder /build/kenkoukun ./
CMD ["./kenkoukun"]
