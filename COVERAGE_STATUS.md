# Coverage Status

- Overall coverage: see coverage/coverage-summary.txt (example snapshot shows ~9.6% with current failing tests).
- Key failing areas to address:
  - API: CORS preflight (OPTIONS 405, missing headers)
  - etcd TLS tests: invalid PEM fixtures leading to parse errors and nil deref
  - Patterns: timeouts due to real etcd client usage in unit tests
- Zero/low coverage packages to raise:
  - pkg/patterns, pkg/inspection, pkg/featureprovider, pkg/examples, pkg/clusterprovider

## Commands

```
# Run all unit tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Summary
go tool cover -func=coverage.out | tee coverage_summary.txt

# Integration tests (tagged)
go test -tags=integration ./integration_tests/...
```

## Recent Additions

- Implemented pkg/algorithms with 35+ tests (sorting/search) including fuzz and benchmarks.
- CI and pre-commit configs present and enforcing coverage gates.

## Next Targets

1) Fix API CORS middleware tests
2) Replace TLS fixtures with valid PEMs for etcd tests
3) Mock etcd in pkg/patterns to avoid timeouts and raise coverage
