FROM alpine:3.3

WORKDIR /go/src/github.com/gaia-docker/tugbot-result-service

COPY .dist/tugbot-result-service /usr/bin/tugbot-result-service

RUN ls -l
RUN ls -l /user/bin

ENTRYPOINT ["/usr/bin/tugbot-result-service"]