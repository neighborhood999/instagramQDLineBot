GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: deps test

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) github.com/line/line-bot-sdk-go/linebot
	$(GOGET) github.com/stretchr/testify/assert
