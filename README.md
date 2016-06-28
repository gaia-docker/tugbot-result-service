# tugbot-result-service

[![CircleCI](https://circleci.com/gh/gaia-docker/tugbot-result-service.svg?style=svg)](https://circleci.com/gh/gaia-docker/tugbot-result-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/gaia-docker/tugbot-result-service)](https://goreportcard.com/report/github.com/gaia-docker/tugbot-result-service)
[![codecov](https://codecov.io/gh/gaia-docker/tugbot-result-service/branch/master/graph/badge.svg)](https://codecov.io/gh/gaia-docker/tugbot-result-service)
[![Docker badge](https://img.shields.io/docker/pulls/gaiadocker/tugbot-result-service.svg)](https://hub.docker.com/r/gaiadocker/tugbot-result-service/)


Tugbot Result Service

Implements [Result Service API](https://github.com/gaia-docker/tugbot/blob/master/doc/proposal/Result%20Service%20API.md#api-design) 
and exposes websocket which present live stream of test results.

## Usage
`docker run gaiadocker/tugbot-result-service tugbot-result-service`

Open `http://result-service-host:8080` to view live stream.
