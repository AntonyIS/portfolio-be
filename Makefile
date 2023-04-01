build:
	go build -o bin/portifolio-be
serve: build
	./bin/portifolio-be
test:
	go test -v ./...