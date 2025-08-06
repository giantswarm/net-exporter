# Cilium uses bpftool to find supported ebpf maps features. Our bpftool image is statically linked.
ARG BPFTOOL_IMAGE=ghcr.io/castai/egressd/bpftool@sha256:d2cf7a30c598e1b39c8b04660d6f1f9ab0925af2951c09216d87eb0d3de0f27b
FROM ${BPFTOOL_IMAGE} as bpftool-dist

FROM quay.io/giantswarm/alpine:3.20.3-giantswarm
# FROM scratch
#
# COPY --from=0 /etc/passwd /etc/passwd
# COPY --from=0 /etc/group /etc/group
COPY --from=bpftool-dist /bin/bpftool /bin/bpftool

ADD net-exporter /
USER giantswarm

ENTRYPOINT ["/net-exporter"]
