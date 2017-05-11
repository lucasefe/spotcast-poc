all: bin/client bin/server

bin/client: bin
	go build -o $@ src/client/*.go

bin/server: bin
	go build -o $@ src/server/*.go

bin pkg:
	mkdir -p $@

cross-compile: bin/client bin/server
	script/cross-compile client
	script/cross-compile server

clean:
	rm -rf pkg bin

.PHONY: cross-compile clean
