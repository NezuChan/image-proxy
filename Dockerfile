FROM golang:1.20 as build-stage

LABEL name "NezukoChan Image Proxy (Docker Build)"
LABEL maintainer "KagChi"

WORKDIR /tmp/build

COPY . .

RUN go build

FROM golang:1.20

LABEL name "NezukoChan Image Proxy"
LABEL maintainer "KagChi"

WORKDIR /app

COPY --from=build-stage /tmp/build/image-proxy image-proxy

CMD ./image-proxy
