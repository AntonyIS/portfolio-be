build:
	go build -o bin/portifolio-be
servr: build
	./bin/portifolio-be
test:
	go test -v ./...