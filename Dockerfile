FROM golang:1.20-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

WORKDIR /build

COPY . ./
RUN go mod download
COPY *.go ./

RUN  go build  -ldflags="-w -s" -o /build/EthTweet .

FROM alpine:3.17

WORKDIR /
RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY tweet.yaml ./tweet.yaml
COPY --from=builder /build/EthTweet /EthTweet

EXPOSE 4001
EXPOSE 8080
ENTRYPOINT   ["/EthTweet"]
