FROM golang:1.24 AS build
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV http_proxy=$HTTP_PROXY
ENV https_proxy=$HTTPS_PROXY
ENV no_proxy=$NO_PROXY
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends git build-essential ca-certificates && rm -rf /var/lib/apt/lists/*
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/openresume ./main.go
ENV http_proxy= https_proxy= no_proxy=

FROM debian:bookworm-slim
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV http_proxy=$HTTP_PROXY
ENV https_proxy=$HTTPS_PROXY
ENV no_proxy=$NO_PROXY
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata chromium \
    fonts-dejavu fonts-noto fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*
ENV http_proxy= https_proxy= no_proxy=
WORKDIR /app
COPY --from=build /app/openresume /usr/local/bin/openresume
RUN mkdir -p /app/uploads
ENV PORT=8080
EXPOSE 8080
CMD ["openresume"]
