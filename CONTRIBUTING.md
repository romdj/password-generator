# Contributing to Password Generator

We use trunk-based development for this project. This means all development happens directly on the `main` branch with frequent integration and minimal branching.

## Development Workflow

### Trunk-Based Development

1. **Work directly on main**: All commits go directly to the `main` branch
2. **Push frequently**: Commit and push small, incremental changes multiple times per day
3. **Keep changes small**: Each commit should be a cohesive, atomic unit of work
4. **Feature flags for incomplete features**: Use feature flags rather than long-lived branches for work-in-progress features

### Making Changes

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd password-generator
   ```

2. **Make your changes**:
   - Keep commits small and focused
   - Write clear commit messages explaining the business context
   - Test your changes locally before pushing

3. **Test your changes**:
   ```bash
   go test ./...
   go build -o pwgen main.go
   ./pwgen --help
   ./pwgen -length 16 -symbols -count 3
   ```

4. **Commit and push directly to main**:
   ```bash
   git add .
   git commit -m "Add feature X: clear description of change"
   git push origin main
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

1. Ensure all tests pass
2. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically create a release with binaries

## Questions or Issues?

- Open an issue for bugs or feature requests
- Discuss larger changes in issues before implementing
- Keep the scope focused on password generation functionality