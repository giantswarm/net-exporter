# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.9.1] - 2020-08-19

### Added

- Added monitoring and common labels.

## [1.9.0] - 2020-06-29

## Added

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

[Unreleased]: https://github.com/giantswarm/net-exporter/compare/v1.9.1...HEAD
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
