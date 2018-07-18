FROM scratch

ADD net-exporter /

ENTRYPOINT ["/net-exporter"]
