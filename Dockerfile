ARG TARGETPLATFORM
ARG TARGETARCH

FROM golang:1.24 AS build
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    git build-essential ca-certificates \
 && rm -rf /var/lib/apt/lists/*

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=$TARGETARCH \
    go build -o /app/cvforge ./main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/cvforge /usr/local/bin/cvforge
RUN mkdir -p /app/uploads

ENV PORT=8080
EXPOSE 8080
CMD ["cvforge"]
