default: run


build:
	go build -o shorturl ./src


run: build
	./shorturl


clean:
	rm -f shorturl

