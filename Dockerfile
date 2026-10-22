FROM --platform=$BUILDPLATFORM gsoci.azurecr.io/giantswarm/alpine:3.24.2 AS certs
RUN apk add --no-cache ca-certificates \
 && addgroup -g 1000 giantswarm \
 && adduser -D -H -u 1000 -G giantswarm -s /sbin/nologin giantswarm

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=certs /etc/passwd /etc/passwd
COPY --from=certs /etc/group /etc/group

ARG TARGETARCH
ADD net-exporter-linux-${TARGETARCH} /net-exporter
USER giantswarm

ENTRYPOINT ["/net-exporter"]
