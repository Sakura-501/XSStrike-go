# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows the migration-first workflow:

- One small feature/change per commit
- Immediate push after each commit
- Frequent test verification

## [Unreleased]

### Planned

- Advanced context scoring and plugin parity improvements

## [1.1.0] - 2026-02-27

### Added

- Reflection analysis chain (`html parser`, `filter checker`, `checker`, context generator`)
- Scan enhancements: DOM findings, WAF detection, path-mode injection, reflection candidate scoring
- Crawl enhancements: link/form extraction, optional blind payload injection, RetireJS component checks
- Active request fuzzing (`--fuzzer --url`) and payload-file bruteforce mode
- Public XSS corpus fetch + dedup pipeline (`benchmarks/corpus`)
- Local benchmark lab + one-command evaluation + baseline report (`benchmarks/reports/2026-02-27-baseline-v1.md`)
- Python compatibility tests and CI workflow (`go test`)
- Release policy document with forced version-bump triggers (`RELEASE_POLICY.md`)

### Changed

- Improved benchmark lab WAF simulation to inspect all query values for block triggers

## [0.1.0-alpha] - 2026-02-27

### Added

- Initial project README with migration roadmap
- Go CLI bootstrap with banner/version output
- Base config defaults and CLI flag parsing
- Utility helpers for headers/URL/parameter parsing
- Minimal payload vector generator and fuzzer mode
- Detailed README usage tutorial and feature explanations

### Notes

- Current phase focuses on foundation and utility migration.
- Behavior parity with Python XSStrike is still in progress.

## Commit Log Snapshot

- `e00ea3d` docs(readme): add detailed usage tutorial and feature descriptions
- `03abde6` feat(fuzzer): port minimal payload generator and fuzzer mode
- `ec6d586` feat(utils): port header and parameter parsing helpers
- `746f459` feat(config): port base defaults and CLI option parsing
- `bf40998` feat(cli): bootstrap Go entrypoint with banner and version flags
- `6efb78b` docs: initialize README with XSStrike-go migration roadmap
