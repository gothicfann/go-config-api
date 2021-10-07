FROM golang:1.17-alpine AS builder
RUN apk add build-base
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app

FROM scratch
COPY --from=builder /app /app
ENTRYPOINT [ "/app" ]