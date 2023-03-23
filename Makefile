build:
	go build -o bin/portifolio-be
run: build
	./bin/portifolio-be
test:
	go test -v ./...