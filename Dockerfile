FROM alpine:3.3

COPY .dist/tugbot-result-service /usr/bin/tugbot-result-service

ENTRYPOINT ["/usr/bin/tugbot-result-service"]