.PHONY: build clean

build:
	go build -o build/pallete.exe

clean:
	rm build/*