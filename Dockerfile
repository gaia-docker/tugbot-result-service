FROM alpine:3.3

ENV RESULT_SERVICE_DIR /go/src/github.com/gaia-docker/tugbot-result-service

WORKDIR $RESULT_SERVICE_DIR

COPY .dist/tugbot-result-service /usr/bin/tugbot-result-service

CMD ["/usr/bin/tugbot-result-service"]