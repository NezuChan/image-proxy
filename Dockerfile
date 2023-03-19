FROM golang:1.20-alpine as build-stage

LABEL name "NezukoChan Image Proxy (Docker Build)"
LABEL maintainer "KagChi"

WORKDIR /tmp/build

RUN apk add --no-cache build-base git python3

COPY package.json .

RUN npm install

COPY . .

RUN npm run build

RUN npm prune --production

FROM golang:1.20-alpine

LABEL name "NezukoChan Image Proxy"
LABEL maintainer "KagChi"

WORKDIR /app

RUN apk add --no-cache tzdata git

COPY --from=build-stage /tmp/build/image-proxy image-proxy

CMD ./image-proxy