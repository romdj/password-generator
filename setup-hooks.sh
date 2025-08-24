#!/bin/bash

# Setup script for git hooks
# This script configures the git hooks directory and installs hooks

set -e

echo "🔧 Setting up git hooks for password generator..."

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Not in a git repository. Please run this from the project root."
    exit 1
fi

# Set the git hooks directory
echo "📁 Configuring git hooks directory..."
git config core.hooksPath .githooks

# Verify hooks are executable
echo "🔑 Setting hook permissions..."
chmod +x .githooks/pre-commit .githooks/pre-push .githooks/commit-msg

# Check for required tools
echo "🔍 Checking for required tools..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
else
    echo "✅ Go $(go version | cut -d' ' -f3) found"
fi

# Check for optional tools
if command -v staticcheck &> /dev/null; then
    echo "✅ staticcheck found"
else
    echo "⚠️  staticcheck not found - install with: go install honnef.co/go/tools/cmd/staticcheck@latest"
fi

if command -v govulncheck &> /dev/null; then
    echo "✅ govulncheck found"
else
    echo "⚠️  govulncheck not found - install with: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi

# Check for conventional commit tools (optional)
if command -v standard-version &> /dev/null; then
    echo "✅ standard-version found"
else
    echo "ℹ️  standard-version not found - install with: npm install -g standard-version"
fi

# Test hooks
echo "🧪 Testing hooks..."

# Test commit-msg hook with a sample message
echo "feat: test commit message" | .githooks/commit-msg /dev/stdin
if [ $? -eq 0 ]; then
    echo "✅ commit-msg hook working"
else
    echo "❌ commit-msg hook failed"
    exit 1
fi

echo ""
echo "🎉 Git hooks setup complete!"
echo ""
echo "📋 What's configured:"
echo "  • Pre-commit: Code formatting, linting, tests"
echo "  • Pre-push: Full test suite, coverage check, cross-platform builds"
echo "  • Commit-msg: Conventional commits validation"
echo ""
echo "📝 Conventional commit format:"
echo "  <type>[optional scope]: <description>"
echo ""
echo "  Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build"
echo ""
echo "💡 To make a release:"
echo "  1. Make changes and commit with conventional commits"
echo "  2. Run: standard-version (or npm version)"
echo "  3. Push tags: git push --follow-tags"
echo ""
echo "🔧 To bypass hooks (emergency only):"
echo "  git commit --no-verify"
echo "  git push --no-verify"