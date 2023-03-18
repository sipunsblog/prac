build:
	go build -o bin/main main.go

run:
	go run main.go

hello:
	echo "hello world"


all: hello run