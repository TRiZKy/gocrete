# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Gocrete
- `gocrete init` command for project initialization
- `gocrete add` command for adding modules
- Base project template with Chi/Gin/Fiber support
- PostgreSQL module with pgx and repository pattern
- MongoDB module with official driver
- OpenAPI module with gen and manual modes
- Docker module with multi-stage Dockerfile
- Goose migrations support
- Structured logging with slog
- Health and readiness endpoints
- Request ID middleware
- Panic recovery middleware
- Environment-based configuration
- Comprehensive documentation

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [1.0.0] - 2024-01-XX

### Added
- Initial public release
- Core generator engine
- Module system architecture
- Embedded templates
- CLI with Cobra
- Comprehensive test suite
- Examples and documentation

---

## Release Process

1. Update version in code
2. Update CHANGELOG.md
3. Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. Create GitHub release with notes from changelog

## Version Numbering

Given a version number MAJOR.MINOR.PATCH:

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)
