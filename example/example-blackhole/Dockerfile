FROM alpine:latest
RUN apk add --update ca-certificates

COPY ./example-blackhole /usr/bin/example-blackhole

ENTRYPOINT ["/usr/bin/example-blackhole"]