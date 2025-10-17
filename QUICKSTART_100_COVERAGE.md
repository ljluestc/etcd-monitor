# Quick Start: Achieving 100% Test Coverage

## 🚀 Immediate Actions (Next 30 Minutes)

### Step 1: Run Baseline Assessment (5 min)
```bash
# Clean cache and run tests
go clean -testcache
go test -v -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View summary
go tool cover -func=coverage.out | tail -20
```

### Step 2: Identify Gaps (5 min)
```bash
# Find packages with lowest coverage
go tool cover -func=coverage.out | sort -k3 -n | head -20

# Find uncovered functions
go tool cover -func=coverage.out | grep -v "100.0%" | grep -v "total:"

# Open HTML report in browser
xdg-open coverage.html  # Linux
# open coverage.html    # macOS
```

### Step 3: Priority Testing (20 min)
Focus on these critical packages first:
1. `pkg/monitor` - Core monitoring
2. `pkg/api` - API endpoints
3. `pkg/etcd/client` - Client wrapper
4. `pkg/clusterprovider` - Cluster management
5. `pkg/inspection` - Health checks

## 📋 Today's Goals

- [ ] Run baseline coverage assessment
- [ ] Identify top 5 packages needing tests
- [ ] Create tests for 1 high-priority package
- [ ] Push to GitHub to trigger CI/CD
- [ ] Review coverage on Codecov

## 🎯 This Week's Milestones

### Day 1 (Today)
- Baseline assessment ✓
- Fix 1-2 critical packages
- Target: 50%+ coverage

### Day 2-3
- `pkg/monitor` → 100%
- `pkg/api` → 100%
- Target: 60%+ coverage

### Day 4-5
- `pkg/etcd/client` → 100%
- `pkg/clusterprovider` → 100%
- Target: 75%+ coverage

### Day 6-7
- Remaining packages
- Integration tests
- Target: 85%+ coverage

### Week 2
- Fine-tune coverage
- Edge cases
- Target: 95%+ coverage

## 💡 Quick Test Template

```go
package yourpackage_test

import (
    "context"
    "testing"
    "time"

    "github.com/etcd-monitor/taskmaster/pkg/yourpackage"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestYourFunction_Success(t *testing.T) {
    // Arrange
    ctx := context.Background()

    // Act
    result, err := yourpackage.YourFunction(ctx, "arg")

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
}

func TestYourFunction_Error(t *testing.T) {
    // Test error path
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()

    _, err := yourpackage.YourFunction(ctx, "arg")
    assert.Error(t, err)
}

func TestYourFunction_EdgeCases(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"empty input", "", "", true},
        {"nil input", "", "", true},
        {"valid input", "test", "test", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := yourpackage.YourFunction(context.Background(), tt.input)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## 🔧 Useful Commands

### Run Specific Package Tests
```bash
go test -v -cover ./pkg/monitor/...
```

### Run Specific Test
```bash
go test -v -run TestMonitor_GetHealth ./pkg/monitor/...
```

### Run Tests with Race Detector
```bash
go test -race ./...
```

### Run Integration Tests
```bash
# Start etcd
docker run -d --name etcd-test -p 2379:2379 quay.io/coreos/etcd:v3.5.9

# Run integration tests
go test -tags=integration -v ./integration_tests/...

# Cleanup
docker rm -f etcd-test
```

### Watch Tests (Using entr)
```bash
# Install entr: apt install entr
find . -name "*.go" | entr -c go test -v ./pkg/monitor/...
```

## 📊 Coverage Targets

| Phase | Target | Deadline |
|-------|--------|----------|
| Baseline | Current | Day 1 |
| Initial | 50% | Day 1 |
| Good | 75% | Week 1 |
| Great | 85% | Week 2 |
| Excellent | 95% | Week 3 |
| Perfect | 100% | Week 4 |

## 🎓 Testing Best Practices

### 1. AAA Pattern (Arrange-Act-Assert)
```go
func TestExample(t *testing.T) {
    // Arrange - Setup
    client := setupTestClient()

    // Act - Execute
    result, err := client.DoSomething()

    // Assert - Verify
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 2. Table-Driven Tests
```go
func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Process(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 3. Test Helpers
```go
func setupTestClient(t *testing.T) *Client {
    t.Helper()
    client := NewClient()
    t.Cleanup(func() { client.Close() })
    return client
}
```

### 4. Mock External Dependencies
```go
type MockEtcdClient struct {
    PutFunc func(ctx context.Context, key, val string) error
}

func (m *MockEtcdClient) Put(ctx context.Context, key, val string) error {
    if m.PutFunc != nil {
        return m.PutFunc(ctx, key, val)
    }
    return nil
}
```

## 🐛 Common Issues & Fixes

### Issue: Tests hang
```bash
# Add timeout
go test -timeout 30s ./...
```

### Issue: Flaky tests
```bash
# Run multiple times
go test -count=100 ./pkg/monitor/...

# Run with race detector
go test -race ./...
```

### Issue: Coverage not updating
```bash
# Clean cache
go clean -testcache
rm coverage.out
go test -coverprofile=coverage.out ./...
```

## 📚 Resources

- **Main Guide**: `TESTING_COMPLETE_GUIDE.md`
- **Integration Tests**: `integration_tests/README.md`
- **Pre-commit Setup**: `PRECOMMIT_SETUP.md`
- **Implementation Summary**: `IMPLEMENTATION_SUMMARY.md`

## ✅ Daily Checklist

- [ ] Run tests: `go test ./...`
- [ ] Check coverage: `go tool cover -func=coverage.out`
- [ ] Write tests for 1-2 functions
- [ ] Commit changes
- [ ] Push to trigger CI/CD
- [ ] Review GitHub Actions results
- [ ] Check Codecov report

## 🚀 Let's Go!

Start with:
```bash
# 1. Baseline
go test -coverprofile=coverage.out ./...

# 2. Analyze
go tool cover -html=coverage.out -o coverage.html

# 3. Act
# Open coverage.html and pick the first red (uncovered) function to test!
```

**Remember**: Progress over perfection. Even 1% improvement each day gets you to 100% in ~20 days!
