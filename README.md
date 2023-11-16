# kenkoukun

## Specification

- Go 1.20

## Setup

> **Note**  
> BOT requires `SERVER MEMBERS INTENT` and `MESSAGE CONTENT INTENT`.

```bash
$ cp .env.sample .env

# create docker container (build included)
$ make docker/run
```

## Make commands

```bash
# start container
$ make docker/start

# stop container
$ make docker/stop

# remove container
$ make docker/rm

# remove container and volume
$ make docker/reset

# view logs
$ make docker/logs
```
