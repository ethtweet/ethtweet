FROM golang:1.21-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . ./

RUN  go build  -ldflags="-w -s" -o /build/EthTweet .

FROM alpine:3

ARG VERSION
ENV VERSION=$VERSION

WORKDIR /app
RUN apk update --no-cache && apk upgrade && apk add --no-cache ca-certificates

COPY templates.zip /app/
COPY Bootstrap.txt /app/Bootstrap.txt
COPY --from=builder /build/EthTweet /app/EthTweet

EXPOSE 4001/tcp
EXPOSE 4001/udp
EXPOSE 8080
ENTRYPOINT   ["/app/EthTweet"]
