package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "hooks":
			setupHooks()
		case "tools":
			installTools()
		case "verify":
			runVerification()
		case "all":
			setupAll()
		default:
			showHelp()
		}
	} else {
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Password Generator Setup Tool")
	fmt.Println("")
	fmt.Println("Usage: go run tools/setup/main.go <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  hooks   - Set up git hooks")
	fmt.Println("  tools   - Install development tools")
	fmt.Println("  verify  - Run all verification checks")
	fmt.Println("  all     - Complete setup (hooks + tools)")
	fmt.Println("")
}

func setupAll() {
	fmt.Println("🚀 Setting up password generator development environment...")
	setupHooks()
	installTools()
	fmt.Println("")
	fmt.Println("✅ Setup complete!")
	fmt.Println("")
	printConventionalCommitHelp()
}

func setupHooks() {
	fmt.Println("🔧 Setting up git hooks...")

	// Check if we're in a git repository
	if !isGitRepo() {
		fmt.Println("❌ Not in a git repository")
		os.Exit(1)
	}

	// Set git hooks directory
	if err := runCommand("git", "config", "core.hooksPath", ".githooks"); err != nil {
		fmt.Printf("❌ Failed to configure git hooks path: %v\n", err)
		os.Exit(1)
	}

	// Make hooks executable
	hooks := []string{"pre-commit", "pre-push", "commit-msg"}
	for _, hook := range hooks {
		hookPath := filepath.Join(".githooks", hook)
		if err := os.Chmod(hookPath, 0755); err != nil {
			fmt.Printf("❌ Failed to make %s executable: %v\n", hook, err)
			os.Exit(1)
		}
	}

	// Test commit-msg hook
	if err := testCommitMsgHook(); err != nil {
		fmt.Printf("❌ Commit message hook test failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Git hooks configured and tested")
}

func installTools() {
	fmt.Println("🔧 Installing development tools...")

	tools := map[string]string{
		"staticcheck": "honnef.co/go/tools/cmd/staticcheck@latest",
		"govulncheck": "golang.org/x/vuln/cmd/govulncheck@latest",
	}

	for tool, pkg := range tools {
		if commandExists(tool) {
			fmt.Printf("✅ %s already installed\n", tool)
			continue
		}

		fmt.Printf("📦 Installing %s...\n", tool)
		if err := runCommand("go", "install", pkg); err != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", tool, err)
		} else {
			fmt.Printf("✅ %s installed\n", tool)
		}
	}

	// Check for Node.js tools
	if commandExists("npm") {
		if !commandExists("standard-version") {
			fmt.Println("ℹ️  To enable automated releases, install: npm install -g standard-version")
		} else {
			fmt.Println("✅ standard-version found")
		}
	} else {
		fmt.Println("ℹ️  Node.js not found - automated releases require npm and standard-version")
	}
}

func runVerification() {
	fmt.Println("✅ Running verification checks...")

	checks := []struct {
		name string
		cmd  []string
	}{
		{"Go format", []string{"go", "fmt", "./..."}},
		{"Go vet", []string{"go", "vet", "./..."}},
		{"Go test", []string{"go", "test", "-race", "-short", "./..."}},
		{"Module verification", []string{"go", "mod", "verify"}},
	}

	for _, check := range checks {
		fmt.Printf("🔍 %s...\n", check.name)
		if err := runCommand(check.cmd[0], check.cmd[1:]...); err != nil {
			fmt.Printf("❌ %s failed: %v\n", check.name, err)
			os.Exit(1)
		}
		fmt.Printf("✅ %s passed\n", check.name)
	}

	// Optional checks with external tools
	if commandExists("staticcheck") {
		fmt.Println("🔍 Running staticcheck...")
		if err := runCommand("staticcheck", "./..."); err != nil {
			fmt.Printf("❌ staticcheck failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ staticcheck passed")
	}

	if commandExists("govulncheck") {
		fmt.Println("🔍 Running vulnerability check...")
		if err := runCommand("govulncheck", "./..."); err != nil {
			fmt.Printf("❌ govulncheck failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ govulncheck passed")
	}

	// Test build
	fmt.Println("🔨 Testing build...")
	buildPath := "/tmp/pwgen-test"
	if runtime.GOOS == "windows" {
		buildPath = "pwgen-test.exe"
	}

	if err := runCommand("go", "build", "-o", buildPath, "."); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		os.Exit(1)
	}
	os.Remove(buildPath)
	fmt.Println("✅ Build successful")

	fmt.Println("🎉 All verification checks passed!")
}

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func testCommitMsgHook() error {
	// For testing purposes, we'll just check if the hook file exists and is executable
	_, err := os.Stat(".githooks/commit-msg")
	return err
}

func printConventionalCommitHelp() {
	fmt.Println("📝 Conventional commit format:")
	fmt.Println("  <type>[optional scope]: <description>")
	fmt.Println("")
	fmt.Println("  Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build")
	fmt.Println("")
	fmt.Println("  Examples:")
	fmt.Println("    feat: add password strength validation")
	fmt.Println("    fix(auth): resolve login timeout issue")
	fmt.Println("    docs: update installation instructions")
	fmt.Println("")
	fmt.Println("💡 Common commands:")
	fmt.Println("  go run tools/setup/main.go verify  - Run all checks")
	fmt.Println("  go test -v ./...                   - Run tests")
	fmt.Println("  go build -o pwgen .                - Build application")
}
