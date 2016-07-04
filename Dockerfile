FROM alpine:3.3

ENV RESULT_SERVICE_DIR /go/src/github.com/gaia-docker/tugbot-result-service

WORKDIR $RESULT_SERVICE_DIR

ADD views $RESULT_SERVICE_DIR/views

ADD .dist/tugbot-result-service $RESULT_SERVICE_DIR/tugbot-result-service

CMD ["./tugbot-result-service"]