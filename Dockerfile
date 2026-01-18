# builder
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -o /bin/wormface-server ./cmd/wormface-server

# prod
FROM alpine:latest
EXPOSE 8080
ARG GIT_TAG
LABEL org.opencontainers.image.version=$GIT_TAG
LABEL org.opencontainers.image.source="https://github.com/OrkWard/wormface"
COPY --from=builder /bin/wormface-server /bin/

CMD ["/bin/womface-server"]
