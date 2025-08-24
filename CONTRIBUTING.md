# Contributing to Password Generator

We use trunk-based development for this project. This means all development happens directly on the `main` branch with frequent integration and minimal branching.

## Development Workflow

### Trunk-Based Development

1. **Work directly on main**: All commits go directly to the `main` branch
2. **Push frequently**: Commit and push small, incremental changes multiple times per day
3. **Keep changes small**: Each commit should be a cohesive, atomic unit of work
4. **Feature flags for incomplete features**: Use feature flags rather than long-lived branches for work-in-progress features

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd password-generator
   ```

2. **Set up development environment**:
   ```bash
   # Install git hooks and development tools
   go run tools/setup/main.go all
   
   # Or set up components individually:
   go run tools/setup/main.go hooks  # Git hooks only
   go run tools/setup/main.go tools  # Dev tools only
   ```

3. **Verify setup**:
   ```bash
   go run tools/setup/main.go verify
   ```

### Making Changes

1. **Use conventional commits**:
   - All commits must follow conventional commit format
   - Format: `<type>[optional scope]: <description>`
   - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`

2. **Development cycle**:
   ```bash
   # Make your changes
   # Pre-commit hooks will run automatically:
   # - Code formatting (go fmt)
   # - Linting (go vet, staticcheck)
   # - Tests (go test)
   # - Security checks (govulncheck)
   
   git add .
   git commit -m "feat: add new password policy template"
   
   # Pre-push hooks run on push to main:
   # - Full test suite with coverage check (>70%)
   # - Cross-platform build verification
   # - Additional quality checks
   
   git push origin main
   ```

3. **Manual verification**:
   ```bash
   go run tools/setup/main.go verify
   ```

## Code Standards

### Go Best Practices
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write self-documenting code
- Handle errors appropriately
- Use `crypto/rand` for secure random number generation

### Testing
- All code should be tested
- Run tests before pushing: `go test ./...`
- Add tests for new functionality
- Test edge cases and error conditions

### Security
- Never commit secrets or sensitive information
- Use secure random number generation
- Follow security best practices for CLI applications

## CI/CD

The project uses GitHub Actions for continuous integration:

- **CI Pipeline**: Runs on every push to `main`
  - Tests with Go 1.25
  - Builds cross-platform binaries
  - Runs static analysis tools

- **Release Pipeline**: Triggered by version tags
  - Builds release binaries for all platforms
  - Creates GitHub releases with downloadable archives

## Pull Requests vs Direct Commits

Since we use trunk-based development:

- **Small changes**: Commit directly to `main`
- **Large changes**: Consider breaking into smaller commits over time
- **Experimental features**: Use feature flags to hide incomplete work
- **Pull Requests**: Only used for:
  - External contributors who don't have write access
  - Large refactoring that benefits from review
  - Changes to CI/CD workflows

## Release Process

This project uses automated semantic versioning based on conventional commits.

### Automated Releases

1. **Ensure all tests pass and code is ready**:
   ```bash
   go run tools/setup/main.go verify
   ```

2. **Generate release** (requires `standard-version`):
   ```bash
   # Install standard-version if not already installed
   npm install -g standard-version
   
   # Generate changelog and create version tag
   standard-version
   
   # Push changes and tags
   git push --follow-tags
   ```

3. **GitHub Actions automatically**:
   - Builds cross-platform binaries
   - Generates changelog from conventional commits
   - Creates GitHub release with downloadable assets
   - Publishes release notes

### Manual Release (if needed)

```bash
# Create and push version tag manually
git tag v1.0.0
git push origin v1.0.0
```

### Release Types

Based on conventional commits:
- `feat:` → Minor version bump (1.0.0 → 1.1.0)
- `fix:` → Patch version bump (1.0.0 → 1.0.1)
- `feat!:` or `BREAKING CHANGE:` → Major version bump (1.0.0 → 2.0.0)

## Questions or Issues?

- Open an issue for bugs or feature requests
- Discuss larger changes in issues before implementing
- Keep the scope focused on password generation functionality