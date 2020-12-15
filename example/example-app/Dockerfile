FROM alpine:latest
RUN apk add --update ca-certificates

COPY ./example-app /usr/bin/example-app

ENTRYPOINT ["/usr/bin/example-app"]