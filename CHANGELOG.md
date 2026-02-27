# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows the migration-first workflow:

- One small feature/change per commit
- Immediate push after each commit
- Frequent test verification

## [Unreleased]

### Added

- Added explicit release policy document (`RELEASE_POLICY.md`) with forced bump triggers
- Added public XSS corpus fetch pipeline with deduplicated payload dataset (`benchmarks/corpus`)
- Added local benchmark lab server and one-command evaluation pipeline (`cmd/xsstrike-lab`, `benchmarks/scripts/evaluate_tool.sh`)
- Added baseline benchmark report for public-corpus evaluation (`benchmarks/reports/2026-02-27-baseline-v1.md`)
- DOM source/sink analyzer scaffold (`internal/dom`)
- Scan report now includes DOM analysis summary and findings
- Target URL normalization with `https -> http` fallback
- Richer migrated default payload/function/fuzz rule sets
- HTTP requester abstraction with timeout/delay/proxy and GET/POST support
- Baseline reflected parameter scanner (`internal/scan`)
- CLI integration for scan flow in default URL mode
- File payload mode (`-f/--file`) in fuzzer workflow
- JSON report writer (`--output/--output-json`)
- Added `db` datasets for compatibility (`wafSignatures.json`, `definitions.json`)
- Signature-based WAF detection and scan integration
- Crawl discovery engine with link/form extraction
- Crawl form scanner with optional blind payload injection
- CLI crawl mode with seeds/depth options
- Payload-file bruteforce mode in scan workflow
- Reflection analysis chain (`html parser`, `filter checker`, `checker`, context generator)
- Path-mode injection support in scan and bruteforce workflows
- RetireJS-compatible component vulnerability scanning integrated into crawl
- Active request fuzzing mode for `--fuzzer --url` workflows
- Python behavior parity tests for selected modules
- GitHub Actions CI workflow (`go test`)

### Changed

- Improved benchmark lab WAF simulation to inspect all query values for block triggers

### Planned

- Advanced context scoring and plugin parity improvements

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
