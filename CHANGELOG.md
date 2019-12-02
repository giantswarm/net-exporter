# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

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

[Unreleased]: https://github.com/giantswarm/net-exporter/compare/v1.2.0...HEAD
