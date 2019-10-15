all: bin/gograz-meetup

bin:
	mkdir -p bin

bin/gograz-meetup: $(shell find . -name '*.go') bin
	go build -o bin/gograz-meetup 

clean:
	rm -rf bin

.PHONY: clean all