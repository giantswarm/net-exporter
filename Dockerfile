FROM alpine:latest

ADD net-exporter /

ENTRYPOINT ["/net-exporter"]
