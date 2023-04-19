build:
	go build -o bin/portifolio-be

serve-pro: build
	./bin/portifolio-be -env=pro

serve-dev: build
	./bin/portifolio-be -env=dev

test:
	go test -v ./...