# net-exporter

Helm chart for net-exporter.

**Homepage:** <https://github.com/giantswarm/net-exporter>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| name | string | `"net-exporter"` |  |
| serviceType | string | `"managed"` |  |
| userID | int | `1000` |  |
| groupID | int | `1000` |  |
| port | int | `8000` |  |
| dns.port | int | `1053` |  |
| dns.label | string | `"coredns"` |  |
| dns.namespace | string | `"kube-system"` |  |
| dns.service | string | `"coredns"` |  |
| dnscache.port | int | `53` |  |
| cluster.kubernetes.DNS.IP | string | `"172.31.0.10"` |  |
| timeout | string | `"5s"` |  |
| image.registry | string | `"gsoci.azurecr.io"` |  |
| image.name | string | `"giantswarm/net-exporter"` |  |
| image.tag | string | `""` |  |
| kubectl.image.registry | string | `"gsoci.azurecr.io"` |  |
| kubectl.image.name | string | `"giantswarm/docker-kubectl"` |  |
| kubectl.image.tag | string | `"1.36.0"` |  |
| daemonset.priorityClassName | string | `"system-node-critical"` |  |
| resources.requests.cpu | string | `"30m"` |  |
| resources.requests.memory | string | `"60Mi"` |  |
| verticalPodAutoscaler.enabled | bool | `true` |  |
| verticalPodAutoscaler.minAllowed.cpu | string | `"30m"` |  |
| verticalPodAutoscaler.minAllowed.memory | string | `"60Mi"` |  |
| verticalPodAutoscaler.maxAllowed.memory | string | `"150Mi"` |  |
| controlPlaneSubnets | list | `[]` |  |
| NetExporter.Hosts | string | `""` |  |
| NetExporter.NTPServers | string | `""` |  |
| NetExporter.DNSCheck.TCP.Disabled | bool | `false` |  |
| ciliumNetworkPolicy.enabled | bool | `false` |  |
| serviceMonitor.enabled | bool | `true` |  |
| serviceMonitor.interval | duration | `"60s"` | Prometheus scrape interval. |
| serviceMonitor.scrapeTimeout | duration | `"45s"` | Prometheus scrape timeout. |
| podSecurityContext.runAsNonRoot | bool | `true` |  |
| podSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| securityContext.readOnlyRootFilesystem | bool | `true` |  |
| securityContext.allowPrivilegeEscalation | bool | `false` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
