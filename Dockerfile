FROM golang:1.18 as builder
WORKDIR /matro/
COPY . /matro/
ARG VERSION=v1.0.1-alpha
ENV VERSION=${VERSION}

WORKDIR /matro/cmd/matro
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
        -X github.com/firstcontributions/matro/internal/commands.Version=${VERSION} \
        -X github.com/firstcontributions/matro/internal/commands.MinVersion=`git rev-parse HEAD` \
        -X github.com/firstcontributions/matro/internal/commands.BuildTime=`date +%FT%T%z` " \
        -o ./build/matro -mod vendor  ./cmd/matro

FROM alpine:3.16 as deploy
COPY --from=builder /matro/build/* ./ 


