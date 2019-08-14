FROM quay.io/giantswarm/alpine:3.9-giantswarm

USER root

ADD nic-exporter /
RUN apk add iproute2 && rm -rf /var/cache/apk/*

USER giantswarm

ENTRYPOINT ["/net-exporter"]
