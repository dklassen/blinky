clean:
	if [ -f "./blinky" ]; then rm ./blinky; fi

build: clean
	go build

run:
	./blinky
