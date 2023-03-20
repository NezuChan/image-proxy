FROM golang:1.20-alpine as build-stage

LABEL name "NezukoChan Image Proxy (Docker Build)"
LABEL maintainer "KagChi"

WORKDIR /tmp/build

COPY . .

# Install needed deps
RUN apk add libc-dev vips-dev gcc g++ make

# Build the project
RUN go build cmd/server/main.go

FROM golang:1.20-alpine

LABEL name "NezukoChan Image Proxy"
LABEL maintainer "KagChi"

WORKDIR /app

# Install needed deps
RUN apk add vips

COPY --from=build-stage /tmp/build/main main

CMD ["./main"]
