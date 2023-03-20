FROM golang:1.20-alpine as build-stage

LABEL name "NezukoChan Image Proxy (Docker Build)"
LABEL maintainer "KagChi"

WORKDIR /tmp/build

COPY . .

# Install deps for vips
RUN apk add vips gcc

# Build the project
RUN go build cmd/server/main.go

FROM golang:1.20-alpine

LABEL name "NezukoChan Image Proxy"
LABEL maintainer "KagChi"

WORKDIR /app

COPY --from=build-stage /tmp/build/main main

CMD ["./main"]