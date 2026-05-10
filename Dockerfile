FROM --platform=$BUILDPLATFORM gsoci.azurecr.io/giantswarm/alpine:3.20.3-giantswarm
FROM scratch

COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/group /etc/group

ARG TARGETARCH
ADD net-exporter-linux-${TARGETARCH} /net-exporter
USER giantswarm

ENTRYPOINT ["/net-exporter"]
