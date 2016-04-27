GOPATH=$(CURDIR)

export $GOPATH

all : build

install:
	go install seckilld

build:
	go build -o bin/seckilld seckilld

clean:
	@rm -f bin/seckilld
	@rm -rf log

