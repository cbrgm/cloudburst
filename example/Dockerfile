FROM alpine:latest
RUN apk add --update ca-certificates

COPY ./example /usr/bin/example

ENTRYPOINT ["/usr/bin/example"]