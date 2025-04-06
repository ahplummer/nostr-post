.PHONY: default
default: displayhelp ;

displayhelp:
	@echo Use "clean, showcoverage, tests, build, docker or run" with make, por favor.

showcoverage: tests
	@echo Running Coverage output
	go tool cover -html=coverage.out

tests: clean
	@echo Running Tests
	go test --coverprofile=coverage.out ./...

run: build
	@echo Running program
	LOG_LEVEL=DEBUG ./bin/nostr-post

build: clean
	@echo Running build command
	go build -o bin/nostr-post src/main.go

clean:
	@echo Removing binary TODO
	rm -rf ./bin ./vendor Gopkg.lock

docker:
	podman build -t nostr-post:latest . -f Dockerfile
	podman run -e CMC_API=$$CMC_API -e NSEC=$$NSEC -it nostr-post:latest

