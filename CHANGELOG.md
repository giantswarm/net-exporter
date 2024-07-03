# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Update Go to 1.22.

## [1.20.0] - 2024-06-17

### Changed

- Add `node` and `app` labels in ServiceMonitor.


## [1.19.0] - 2024-04-23

### Added

- Add `/blackbox` endpoint.

## [1.18.2] - 2023-12-13

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.

## [1.18.1] - 2023-12-01

### Changed

- Update quay.io/giantswarm/alpine Docker tag to v3.18.5 ([#319](https://github.com/giantswarm/net-exporter/pull/319))

## [1.18.0] - 2023-09-28

### Changed

- Enable PSP resource deployment based on global value.

## [1.17.1] - 2023-09-20

### Fixed

- Fix kubernetes network policy to allow scraping on cilium-less clusters.

## [1.17.0] - 2023-06-27

### Changed

- Add security context values to make chart comply to PodSecurityStandard restricted profile.

## [1.16.2] - 2023-06-13

### Changed

- Reduce CPU and Mem requests.

## [1.16.1] - 2023-06-02

### Added

- Add service monitor to be scraped by Prometheus Agent.

## [1.16.0] - 2023-05-31

### Added

- Add values for the daemonset resources ([#280](https://github.com/giantswarm/net-exporter/pull/280)).

## [1.15.0] - 2023-05-04

### Changed

- Allow requests from the api-server.
- Disable PSPs for k8s 1.25 and newer.

## [1.14.1] - 2023-04-24

### Changed

- Update icon.

## [1.14.0] - 2023-04-04

### Added

- Add `Cilium Network Policy` to net-exporter.

### Changed

- Don't push net-exporter to capa-app-collection because it's already a default app.
- Don't push net-exporter to cloud-director-app-collection and vsphere-app-collection because it's already in default app bundles.
- Adjust VPA resources.

## [1.13.0] - 2022-12-19

### Added

- Add helm chart values schema

### Changed

- Update to Go 1.18
- Update github.com/giantswarm/k8sclient to v7.0.1
- Update github.com/giantswarm/micrologger to v1.0.0
- Update github.com/miekg/dns to v1.1.50
- Update k8s.io deps to v0.26.0
- Update docker-kubectl to 1.25.4

## [1.12.0] - 2022-03-16

### Changed

- Use parameter for CoreDNS namespace (defaulted to kube-system)

## [1.11.0] - 2022-03-07

### Added

- Add networkpolicy to allow egress towards `k8s-dns-node-cache-app` endpoints.

## [1.10.3] - 2021-08-12

### Changed

- Prepare helm values to configuration management.
- Update architect-orb to v4.0.0.

## [1.10.2] - 2021-05-20

### Changed

- Allow to customize dns service.
- Only check pod existence on dial errors. Check pod deletion directly by IP instead of listing pods and searching.

## [1.10.1] - 2021-04-29

## [1.10.0] - 2021-04-20

### Changed

- Add label selector for pods to help lower memory usage

## [1.9.3] - 2021-03-26

### Changed

- Set docker.io as the default registry
- Update kubectl image to v1.18.8.

## [1.9.2] - 2020-08-21

### Changed

- Updated backward incompatible Kubernetes dependencies to v1.18.5.

### Fixed

- Fixed indentation problem with the daemonset template.

## [1.9.1] - 2020-08-19

### Added

- Added monitoring and common labels.

## [1.9.0] - 2020-06-29

### Added

- Add `ntp` collector.

## [1.8.1] - 2020-06-17

### Changed

- Added 100.64.0.0/10 to the allowed egress subnets in NetworkPolicy.

## [1.8.0]

### Changed

- Deploy as a unique app in app collection.

## [1.7.1] 2020-04-01

### Changed

- Change daemonset to use release revision not time for Helm 3 support.
- Only set hosts arg if a value is present.
- Remove label from role ref in cluster role binding.

## [1.7.0] 2020-03-20

### Changed

- Dial error if the Pod doesn't exist anymore will be ignored.

## [1.6.0] 2020-01-29

### Changed

- Allow to disable DNS TCP check.
- Allow custom internal domain configuration for dns collector.

## [1.5.1] 2020-01-08

### Changed

- Changed Priority Class to `system-node-critical` for net-exporter deployed in TC.

## [1.5.0] 2020-01-07

### Changed

- Changed Priority Class to `system-node-critical`.

## [1.4.3] 2019-12-27

### Changed

- Fixed invalid image reference.

## [1.4.2] 2019-12-24

### Changed

- Restore CPU requests.

## [1.4.1] 2019-11-29

### Changed

- Make image registry configurable for namespace labeler init container.

## [1.4.0] 2019-11-21

### Changed

- Change ping behavior to use a ring instead of a mesh
- Reduce number of latency buckets from 15 to 5
- Fix DNS 5 dots issue for test installations
- Update README to align with other apps

## [1.3.0] 2019-10-24

### Changed

- Add net-exporter to default app catalog

## [1.2.0] 2019-07-17

### Changed

- Tolerations changed to tolerate all taints.
- Change prioty class to `giantswarm-critical`.

[Unreleased]: https://github.com/giantswarm/net-exporter/compare/v1.20.0...HEAD
[1.20.0]: https://github.com/giantswarm/net-exporter/compare/v1.19.0...v1.20.0
[1.19.0]: https://github.com/giantswarm/net-exporter/compare/v1.18.2...v1.19.0
[1.18.2]: https://github.com/giantswarm/net-exporter/compare/v1.18.1...v1.18.2
[1.18.1]: https://github.com/giantswarm/net-exporter/compare/v1.18.0...v1.18.1
[1.18.0]: https://github.com/giantswarm/net-exporter/compare/v1.17.1...v1.18.0
[1.17.1]: https://github.com/giantswarm/net-exporter/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/giantswarm/net-exporter/compare/v1.16.2...v1.17.0
[1.16.2]: https://github.com/giantswarm/net-exporter/compare/v1.16.1...v1.16.2
[1.16.1]: https://github.com/giantswarm/net-exporter/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/giantswarm/net-exporter/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/giantswarm/net-exporter/compare/v1.14.1...v1.15.0
[1.14.1]: https://github.com/giantswarm/net-exporter/compare/v1.14.0...v1.14.1
[1.14.0]: https://github.com/giantswarm/net-exporter/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/giantswarm/net-exporter/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/giantswarm/net-exporter/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/giantswarm/net-exporter/compare/v1.10.3...v1.11.0
[1.10.3]: https://github.com/giantswarm/net-exporter/compare/v1.10.2...v1.10.3
[1.10.2]: https://github.com/giantswarm/net-exporter/compare/v1.10.1...v1.10.2
[1.10.1]: https://github.com/giantswarm/net-exporter/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/giantswarm/net-exporter/compare/v1.9.3...v1.10.0
[1.9.3]: https://github.com/giantswarm/net-exporter/compare/v1.9.2...v1.9.3
[1.9.2]: https://github.com/giantswarm/net-exporter/compare/v1.9.1...v1.9.2
[1.9.1]: https://github.com/giantswarm/net-exporter/compare/v1.9.0...v1.9.1
[1.9.0]: https://github.com/giantswarm/net-exporter/compare/v1.8.1...v1.9.0
[1.8.1]: https://github.com/giantswarm/net-exporter/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/giantswarm/net-exporter/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/giantswarm/net-exporter/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/giantswarm/net-exporter/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/giantswarm/net-exporter/compare/v1.5.1...v1.6.0
[1.5.1]: https://github.com/giantswarm/net-exporter/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/giantswarm/net-exporter/compare/v1.4.3...v1.5.0
[1.4.3]: https://github.com/giantswarm/net-exporter/compare/v1.4.2...v1.4.3
[1.4.2]: https://github.com/giantswarm/net-exporter/compare/v1.4.1...v1.4.2
[1.4.1]: https://github.com/giantswarm/net-exporter/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/giantswarm/net-exporter/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.3.0
