# builder
FROM golang:1.25.4-alpine AS builder
ARG GIT_TAG
LABEL org.opencontainers.image.version=$GIT_TAG
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-w -s" -o /bin/wormface-server ./cmd/wormface-server

# prod
FROM alpine:latest
COPY --from=builder /bin/wormface-server /bin/
EXPOSE 8080
CMD ["/bin/womface-server"]
