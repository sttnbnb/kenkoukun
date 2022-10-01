include .env

DOCKER_FILE=./build/Dockerfile
ENV_FILE=.env
APP_NAME=kenkoukun

docker/run: docker/rm docker/build
	docker run -d --name $(APP_NAME) --env-file=$(ENV_FILE) kenkoukun:latest

docker/start:
	docker start $(APP_NAME)

docker/stop:
	docker stop $(APP_NAME)

docker/rm:
	docker rm -f $(APP_NAME)

docker/build:
	go mod vendor
	docker build -t $(APP_NAME) -f $(DOCKER_FILE) .

docker/logs:
	docker logs -f $(APP_NAME)
