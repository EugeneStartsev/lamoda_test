all: build up

build:
	go build -o bin/ ./ && \
	cd .. && \
	docker build -t alpine3.16_worker-test -f Dockerfile .

up:
	cd .. && \
	docker-compose -f docker-compose.yml up
