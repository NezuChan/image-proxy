FROM golang:1.22-alpine as build-stage

WORKDIR /tmp/build

COPY . .

# Install needed deps
RUN apk add --no-cache libc-dev vips-dev gcc g++ make

# Build the project
RUN go build cmd/server/main.go

FROM alpine:3

LABEL name "NezukoChan Image Proxy"
LABEL maintainer "KagChi"

WORKDIR /app

# Install needed deps
RUN apk add --no-cache vips tini

COPY --from=build-stage /tmp/build/main main

ENTRYPOINT ["tini", "--"]
CMD ["/app/main"]
