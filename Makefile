tag=$(shell git describe --tags --abbrev=0)

.PHONY: deps build-linux docker cleanup

deps:
	go mod download

build-linux: deps
	mkdir bin/
	GOOS=linux GOARCH=amd64 go build -o bin/resque_exporter .

docker: build-linux 
	docker build -t gg1113/resque_exporter:$(tag) .
	docker push gg1113/resque_exporter:$(tag)

cleanup:
	rm -Rf bin/