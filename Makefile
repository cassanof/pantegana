.PHONY: all test clean

pantegana-dir := $(shell pwd)

all: build-client build-server

build-client:
	mkdir -p out; \
		cd ./main/client; \
		go build -o client.bin; \
		mv client.bin $(pantegana-dir)/out/client.bin;

build-server:
	mkdir -p out; \
		cd ./main/server; \
		go build -o server.bin; \
		mv server.bin $(pantegana-dir)/out/server.bin;

clean:
	rm -fr ./out
