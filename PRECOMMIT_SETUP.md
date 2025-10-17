# Pre-commit Hook Setup

## Installation

### Option 1: Using pipx (Recommended)
```bash
# Install pipx if not already installed
sudo apt-get install pipx

# Install pre-commit
pipx install pre-commit

# Install the git hooks
pre-commit install
```

### Option 2: Using pip with virtual environment
```bash
# Create virtual environment
python3 -m venv .venv

# Activate virtual environment
source .venv/bin/activate

# Install pre-commit
pip install pre-commit

# Install the git hooks
pre-commit install
```

### Option 3: Manual setup
```bash
# Copy the pre-commit hook manually
cp .pre-commit-config.yaml .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## Running Pre-commit Hooks

### Run on all files
```bash
pre-commit run --all-files
```

### Run on specific files
```bash
pre-commit run --files pkg/**/*.go
```

### Skip pre-commit hooks (not recommended)
```bash
git commit --no-verify
```

## Pre-commit Hooks Configured

The following hooks are configured in `.pre-commit-config.yaml`:

1. **Trailing whitespace** - Removes trailing whitespace
2. **End of file fixer** - Ensures files end with a newline
3. **YAML check** - Validates YAML syntax
4. **Large files check** - Prevents committing large files
5. **Merge conflict check** - Detects merge conflict markers
6. **JSON check** - Validates JSON syntax
7. **TOML check** - Validates TOML syntax
8. **XML check** - Validates XML syntax
9. **Case conflict check** - Checks for case conflicts
10. **Debug statements** - Detects debug statements
11. **Private key detection** - Prevents committing private keys
12. **Mixed line endings** - Fixes mixed line endings
13. **golangci-lint** - Runs Go linting
14. **go fmt** - Formats Go code
15. **go vet** - Runs Go vet
16. **go test** - Runs Go tests
17. **go mod tidy** - Tidies go.mod
18. **go build** - Builds Go code

## Troubleshooting

### Pre-commit not found
If `pre-commit` command is not found after installation, ensure it's in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Hooks failing
If hooks are failing, run them individually to debug:
```bash
pre-commit run <hook-id> --all-files
```

### Updating hooks
To update hooks to the latest versions:
```bash
pre-commit autoupdate
```
