FROM golang:1.17 as build_service
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /crypto-watcher
COPY .env /crypto-watcher/cmd/crypto-watcher

WORKDIR /crypto-watcher/cmd/crypto-watcher
RUN go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:3.15
RUN addgroup -g 1000 -S crypto && \
    adduser -u 1000 -h /app -G crypto -S crypto

COPY --from=build_service /crypto-watcher/cmd/crypto-watcher /app
COPY --from=build_service /crypto-watcher/cmd/crypto-watcher/.env /app/
WORKDIR /app
USER crypto
CMD ["./crypto-watcher"]
