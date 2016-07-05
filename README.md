# tugbot-result-service

[![CircleCI](https://circleci.com/gh/gaia-docker/tugbot-result-service.svg?style=svg)](https://circleci.com/gh/gaia-docker/tugbot-result-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/gaia-docker/tugbot-result-service)](https://goreportcard.com/report/github.com/gaia-docker/tugbot-result-service)
[![codecov](https://codecov.io/gh/gaia-docker/tugbot-result-service/branch/master/graph/badge.svg)](https://codecov.io/gh/gaia-docker/tugbot-result-service)
[![Docker](https://img.shields.io/docker/pulls/gaiadocker/tugbot-result-service.svg)](https://hub.docker.com/r/gaiadocker/tugbot-result-service/)
[![Docker Image Layers](https://imagelayers.io/badge/gaiadocker/tugbot-result-service:latest.svg)](https://imagelayers.io/?images=gaiadocker/tugbot-result-service:latest 'Get your own badge on imagelayers.io')


Tugbot Result Service

Implements [Result Service API](https://github.com/gaia-docker/tugbot/blob/master/doc/proposal/Result%20Service%20API.md#api-design) 
and exposes websocket which present live stream of test results.

## Usage
```
$ tugbot-result-service help

NAME:
   tugbot-result-service - Implements Result Service API and exposes websocket which present live stream of test results.

USAGE:
   tugbot-result-service [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value   http service port (default: "8080")
   --loglevel value, -l value  log level (default: "debug")
   --help, -h                  show help
   --version, -v               print the version
```

## Run as docker container
`docker run -it --name result-service -p 8080:8080 gaiadocker/tugbot-result-service`

## View live stream
Open `http://result-service-host:8080`.
