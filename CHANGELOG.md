# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.7.0] 2020-03-20

### Changed

- Dial error if the Pod doesn't exist anymore will be ignored.

## [v1.6.0] 2020-01-29

### Changed

- Allow to disable DNS TCP check.
- Allow custom internal domain configuration for dns collector.

## [v1.5.1] 2020-01-08

### Changed

- Changed Priority Class to `system-node-critical` for net-exporter deployed in TC.

## [v1.5.0] 2020-01-07

### Changed

- Changed Priority Class to `system-node-critical`.

## [v1.4.3] 2019-12-27

### Changed

- Fixed invalid image reference.

## [v1.4.2] 2019-12-24

### Changed

- Restore CPU requests.

## [v1.4.1] 2019-11-29

### Changed

- Make image registry configurable for namespace labeler init container.

## [v1.4.0] 2019-11-21

### Changed

- Change ping behavior to use a ring instead of a mesh
- Reduce number of latency buckets from 15 to 5
- Fix DNS 5 dots issue for test installations
- Update README to align with other apps

## [v1.3.0] 2019-10-24

### Changed

- Add net-exporter to default app catalog

## [v1.2.0] 2019-07-17

### Changed

- Tolerations changed to tolerate all taints.
- Change prioty class to `giantswarm-critical`.

[v1.7.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.7.0
[v1.6.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.6.0
[v1.5.1]: https://github.com/giantswarm/net-exporter/releases/tag/v1.5.1
[v1.5.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.5.0
[v1.4.3]: https://github.com/giantswarm/net-exporter/releases/tag/v1.4.3
[v1.4.2]: https://github.com/giantswarm/net-exporter/releases/tag/v1.4.2
[v1.4.1]: https://github.com/giantswarm/net-exporter/releases/tag/v1.4.1
[v1.4.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.4.0
[v1.3.0]: https://github.com/giantswarm/net-exporter/releases/tag/v1.3.0
