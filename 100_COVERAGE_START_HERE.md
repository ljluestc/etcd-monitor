# 🚀 START HERE: 100% Test Coverage Initiative

Welcome to the etcd-monitor 100% Test Coverage Initiative!

This guide will help you get started with implementing comprehensive test coverage for this project.

---

## Quick Status

```
Current Coverage:  15.3% ████░░░░░░░░░░░░░░░░  Target: 100%
Timeline:          10 weeks (2025-10-19 to 2025-12-28)
Current Phase:     Phase 1 - Foundation (Target: 30%)
Status:            🟢 READY TO START
```

---

## 📚 Essential Reading (in order)

1. **`100_PERCENT_COVERAGE_SESSION_REPORT.md`** ⭐ START HERE
   - Executive summary
   - Current state analysis
   - Quick wins and immediate actions

2. **`100_PERCENT_COVERAGE_TASKMASTER.md`** 📋 PROJECT PLAN
   - Detailed task breakdown
   - Weekly milestones
   - Acceptance criteria
   - Progress tracking

3. **`TESTING_STRATEGY_100_COVERAGE.md`** 🎯 STRATEGY
   - Comprehensive testing approach
   - Infrastructure requirements
   - Embedded etcd setup
   - Best practices

---

## 🏃 Quick Start (< 5 minutes)

### 1. Verify Current State

```bash
# Clone/navigate to project
cd /path/to/etcd-monitor

# Run all tests
go test ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
# Expected: total: (statements) 15.3%

# View HTML report
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

### 2. Run Comprehensive Test Suite

```bash
# Using Python orchestration script
python3 test_comprehensive.py --all

# Or run specific test types
python3 test_comprehensive.py --unit
python3 test_comprehensive.py --integration
python3 test_comprehensive.py --benchmark
python3 test_comprehensive.py --coverage
```

### 3. Set Up Pre-Commit Hooks

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Test hooks
pre-commit run --all-files
```

---

## 🎯 Phase 1: Foundation (Weeks 1-2) - Next Steps

### Priority 1: Complete pkg/k8s (4 hours)
**Current:** 76.2% → **Target:** 100%

```bash
# Create test file
touch pkg/k8s/builder_edge_cases_test.go

# Add tests for:
# - Nil config handling
# - Invalid kubeconfig paths
# - Connection timeouts
# - Authentication failures
```

**Test Template:**
```go
package k8s

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBuilder_NilConfig(t *testing.T) {
    builder := NewSimpleClientBuilder(nil)
    assert.NotNil(t, builder)

    _, err := builder.Config()
    assert.Error(t, err)
}

func TestBuilder_InvalidKubeconfig(t *testing.T) {
    builder := NewSimpleClientBuilder(&Config{
        KubeconfigPath: "/nonexistent/path",
    })

    _, err := builder.Config()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "no such file")
}
```

**Run Tests:**
```bash
go test -v ./pkg/k8s/... -coverprofile=coverage/k8s.out
go tool cover -func=coverage/k8s.out
# Target: 100%
```

---

### Priority 2: Complete pkg/etcdctl (6 hours)
**Current:** 68.3% → **Target:** 100%

```bash
# Create comprehensive test file
touch pkg/etcdctl/wrapper_complete_test.go
```

**What to Test:**
- All KV operations (Put, Get, Delete, Watch)
- Cluster operations (MemberList, Status, AlarmList)
- Error handling
- Context cancellation
- Timeout scenarios

**Key Test:**
```go
func TestWrapper_AllOperations(t *testing.T) {
    // Use mock client from testutil/mocks
    mockClient := mocks.NewMockEtcdClient()
    wrapper := NewWrapper(mockClient, logger)

    // Test Put
    err := wrapper.Put(ctx, "key", "value")
    assert.NoError(t, err)

    // Test Get
    value, err := wrapper.Get(ctx, "key")
    assert.NoError(t, err)
    assert.Equal(t, "value", value)

    // Test Delete
    err = wrapper.Delete(ctx, "key")
    assert.NoError(t, err)
}
```

---

### Priority 3: Enhance pkg/api (8 hours)
**Current:** 28.6% → **Target:** 80%

```bash
# Create new test files
touch pkg/api/server_lifecycle_test.go
touch pkg/api/handlers_complete_test.go
touch pkg/api/concurrency_test.go
```

**Critical Tests:**
- Server Start/Stop lifecycle
- All HTTP endpoint handlers
- Concurrent request handling
- Error responses (4xx, 5xx)
- Metrics endpoint validation

---

### Priority 4: Enhance pkg/etcd (12 hours)
**Current:** 27.8% → **Target:** 80%

```bash
# Create specialized test files
touch pkg/etcd/client_tls_test.go
touch pkg/etcd/client_auth_test.go
touch pkg/etcd/client_retry_test.go
touch pkg/etcd/health_complete_test.go
```

**Focus Areas:**
- TLS configuration (certs, keys, CA)
- Authentication (username/password, client certs)
- Connection retry logic
- K8s secret-based configuration
- Health monitoring
- Statistics collection

---

## 🛠️ Development Workflow

### Daily Workflow

```bash
# 1. Pull latest changes
git pull origin main

# 2. Create feature branch
git checkout -b test/improve-pkg-k8s-coverage

# 3. Write tests
# (edit test files)

# 4. Run tests locally
go test -v ./pkg/k8s/... -coverprofile=coverage/k8s.out

# 5. Check coverage
go tool cover -func=coverage/k8s.out

# 6. Run all tests
go test ./...

# 7. Run with race detection
go test -race ./...

# 8. Commit (pre-commit hooks will run)
git add .
git commit -m "test: improve pkg/k8s coverage to 100%"

# 9. Push and create PR
git push origin test/improve-pkg-k8s-coverage
```

### PR Checklist

Before creating a PR, ensure:
- ✅ All tests pass (`go test ./...`)
- ✅ No race conditions (`go test -race ./...`)
- ✅ Coverage increased or stayed same
- ✅ New code has > 80% coverage
- ✅ golangci-lint passes
- ✅ Pre-commit hooks pass
- ✅ Documentation updated if needed

---

## 📊 Tracking Progress

### Check Current Coverage

```bash
# Overall coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# Per-package coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "pkg/k8s"
go tool cover -func=coverage.out | grep "pkg/api"
# etc.
```

### Generate Coverage Report

```bash
# HTML report
go tool cover -html=coverage.out -o coverage/coverage.html

# Text report
go tool cover -func=coverage.out | tee coverage/coverage-summary.txt

# JSON report (for automation)
go test -json ./... > coverage/test-results.json
```

### Update Taskmaster

After completing a task, update `100_PERCENT_COVERAGE_TASKMASTER.md`:

```markdown
| Task ID | Package | Current | Target | Status | Owner | Due Date |
|---------|---------|---------|--------|--------|-------|----------|
| P1-001 | pkg/k8s | 76.2% | 100% | ✅ DONE | Your Name | 2025-10-23 |
```

---

## 🧪 Test Infrastructure

### Using Mock etcd Client

```go
import "github.com/etcd-monitor/taskmaster/testutil/mocks"

func TestWithMockClient(t *testing.T) {
    // Create mock client
    mockClient := mocks.NewMockEtcdClient()

    // Configure responses
    mockClient.SetGetResponse("key", "value")

    // Use in tests
    value, err := mockClient.Get(ctx, "key")
    assert.NoError(t, err)
    assert.Equal(t, "value", value)

    // Verify operations
    counts := mockClient.GetOperationCounts()
    assert.Equal(t, 1, counts["Get"])
}
```

### Creating Test Fixtures

```bash
# Create fixture directory
mkdir -p testutil/fixtures/etcd
mkdir -p testutil/fixtures/k8s
mkdir -p testutil/fixtures/metrics

# Add sample data
cat > testutil/fixtures/etcd/sample_config.yaml <<EOF
endpoints:
  - localhost:2379
dialTimeout: 5s
username: test
password: test
EOF
```

---

## 🔍 Common Testing Patterns

### Pattern 1: Table-Driven Tests

```go
func TestFunction_TableDriven(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "TEST", false},
        {"empty input", "", "", true},
        {"nil input", nil, "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### Pattern 2: Subtests for Related Scenarios

```go
func TestServer(t *testing.T) {
    t.Run("Lifecycle", func(t *testing.T) {
        server := NewServer(config)

        t.Run("Start", func(t *testing.T) {
            err := server.Start()
            assert.NoError(t, err)
        })

        t.Run("Stop", func(t *testing.T) {
            err := server.Stop()
            assert.NoError(t, err)
        })
    })
}
```

### Pattern 3: Mock-Based Testing

```go
func TestWithMock(t *testing.T) {
    mock := &MockInterface{}
    mock.On("Method", "arg").Return("result", nil)

    result, err := mock.Method("arg")

    assert.NoError(t, err)
    assert.Equal(t, "result", result)
    mock.AssertExpectations(t)
}
```

---

## 🚨 Troubleshooting

### Tests Failing Locally

```bash
# Clean and rebuild
go clean -testcache
go test ./...

# Run specific test
go test -v -run TestName ./pkg/package/...

# Run with more verbosity
go test -v -x ./...
```

### Coverage Not Updating

```bash
# Remove old coverage files
rm coverage.out coverage/*.out

# Regenerate
go test -coverprofile=coverage.out ./...
```

### Pre-Commit Hooks Failing

```bash
# See what failed
pre-commit run --all-files --verbose

# Skip hooks temporarily (NOT recommended)
git commit --no-verify

# Update hooks
pre-commit autoupdate
```

### CI/CD Failures

1. Check GitHub Actions logs
2. Reproduce locally:
   ```bash
   go test -v -race -coverprofile=coverage.out ./...
   ```
3. Fix issues
4. Push again

---

## 📞 Getting Help

### Resources

- **Strategy Doc:** `TESTING_STRATEGY_100_COVERAGE.md`
- **Taskmaster:** `100_PERCENT_COVERAGE_TASKMASTER.md`
- **Session Report:** `100_PERCENT_COVERAGE_SESSION_REPORT.md`
- **Existing Tests:** Look at `pkg/signals/*_test.go` for examples

### Questions?

- Review existing test files for patterns
- Check the strategy document for approach
- Look at taskmaster for specific task details

---

## 🎯 Success Criteria

### Phase 1 Success (Week 2)
- [ ] Overall coverage >= 30%
- [ ] pkg/k8s at 100%
- [ ] pkg/etcdctl at 100%
- [ ] pkg/api at 80%
- [ ] pkg/etcd at 80%
- [ ] All tests passing
- [ ] CI/CD green

### Your First Contribution
- [ ] Pick a task from Phase 1
- [ ] Write comprehensive tests
- [ ] Achieve target coverage
- [ ] All tests pass
- [ ] Create PR with:
  - Test coverage improvement shown
  - Clear description
  - Before/after coverage numbers

---

## 🚀 Let's Get Started!

**Recommended First Task:**

```bash
# Start with pkg/k8s (easiest, quickest win)
git checkout -b test/pkg-k8s-100-coverage
code pkg/k8s/builder_edge_cases_test.go

# Follow the test template above
# Run tests frequently
# Aim for 100% coverage
# Create PR when done
```

**Estimated Time:** 4 hours
**Impact:** ~1-2% overall coverage increase
**Difficulty:** ⭐⭐☆☆☆ Easy

---

## 📈 Weekly Goals

### Week 1 (2025-10-20 to 2025-10-26)
**Target:** 22% coverage (+6.7%)

**Tasks:**
- ✅ Complete pkg/k8s (4 hours)
- ✅ Complete pkg/etcdctl (6 hours)

### Week 2 (2025-10-27 to 2025-11-02)
**Target:** 30% coverage (+8%)

**Tasks:**
- ✅ Enhance pkg/api (8 hours)
- ✅ Enhance pkg/etcd (12 hours)

---

## 🎓 Learning Path

1. **Day 1:** Read session report, understand current state
2. **Day 2:** Study existing test patterns in pkg/signals
3. **Day 3:** Start with pkg/k8s tests (easiest)
4. **Day 4-5:** Continue pkg/k8s, aim for 100%
5. **Week 2:** Move to pkg/etcdctl
6. **Week 3+:** Progress through taskmaster phases

---

**Status:** 🟢 READY TO START
**Next Action:** Pick a Phase 1 task and start coding!
**Questions?** Review the documentation files or check existing test patterns.

Good luck! 🚀
