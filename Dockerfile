FROM golang:1.20-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . ./

RUN  go build  -ldflags="-w -s" -o /build/EthTweet .

FROM alpine:3

WORKDIR /
RUN apk update --no-cache && apk upgrade && apk add --no-cache ca-certificates

COPY Bootstrap.txt ./Bootstrap.txt
COPY --from=builder /build/EthTweet /EthTweet

EXPOSE 4001
EXPOSE 8080
ENTRYPOINT   ["/EthTweet"]
