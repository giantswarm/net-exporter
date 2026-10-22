FROM --platform=$BUILDPLATFORM gsoci.azurecr.io/giantswarm/alpine:3.24.3 AS base
RUN addgroup -g 1000 giantswarm \
 && adduser -D -H -u 1000 -G giantswarm -s /sbin/nologin giantswarm

FROM scratch

COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

ARG TARGETARCH
ADD net-exporter-linux-${TARGETARCH} /net-exporter
USER giantswarm

ENTRYPOINT ["/net-exporter"]
