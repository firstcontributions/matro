VERSION="1.0.0-alpha"

config:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

build:
	cd cmd/matro && go build -ldflags "-s -w \
	-X github.com/firstcontributions/matro/internal/commands.Version=${VERSION} \
	-X github.com/firstcontributions/matro/internal/commands.MinVersion=`git rev-parse HEAD` \
	-X github.com/firstcontributions/matro/internal/commands.BuildTime=`date +%FT%T%z` " \
	-o ${GOPATH}/bin/matro
	