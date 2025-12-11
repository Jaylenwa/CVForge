# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /app
RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o openresume ./main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata chromium nss freetype ttf-dejavu font-noto font-noto-cjk
WORKDIR /app
COPY --from=build /app/openresume /usr/local/bin/openresume
RUN mkdir -p /app/uploads
ENV PORT=8080
EXPOSE 8080
CMD ["openresume"]

