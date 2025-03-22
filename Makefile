all: build run clean

build:
	go build -o main cmd/main.go

run:
	./main

clean:
	rm main