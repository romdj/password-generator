# Password Generator

A secure, configurable password generator written in Go with advanced features including password strength analysis, policy validation, and flexible configuration options.

## Features

### Core Features
- **Configurable Length**: Set password length (default: 12 characters)
- **Character Types**: Choose from uppercase, lowercase, digits, and symbols
- **Ambiguous Character Exclusion**: Option to exclude confusing characters (0, O, 1, l, I)
- **Multiple Passwords**: Generate multiple passwords in one command
- **Secure**: Uses `crypto/rand` for cryptographically secure randomness

### Advanced Features
- **Password Strength Analysis**: Real-time strength scoring with entropy calculations and feedback
- **Policy Templates**: Pre-built templates for common requirements (basic, corporate, high-security, AWS, Azure, PCI-DSS)
- **Configuration Files**: YAML/JSON config files for default settings
- **Environment Variables**: Override settings with environment variables
- **Password Validation**: Validate existing passwords against policies

## Installation

```bash
git clone <repository-url>
cd password-generator
go build -o pwgen .
```

Or install directly:
```bash
go install github.com/romdj/password-generator@latest
```

## Usage

### Basic Usage

```bash
# Generate a 12-character password (default)
go run main.go

# Build and use the binary
./pwgen

# Generate with strength analysis
./pwgen --strength
```

### Configuration Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--length` | `-l` | 12 | Password length |
| `--upper` | `-u` | true | Include uppercase letters |
| `--lower` | `-L` | true | Include lowercase letters |
| `--digits` | `-d` | true | Include digits |
| `--symbols` | `-s` | false | Include symbols |
| `--no-ambiguous` | `-n` | false | Exclude ambiguous characters |
| `--count` | `-c` | 1 | Number of passwords to generate |
| `--strength` | `-S` | false | Show password strength analysis |
| `--policy` | `-p` | "" | Apply password policy template |

### Special Commands

| Flag | Description |
|------|-------------|
| `--list-policies` | List available password policy templates |
| `--validate "password"` | Validate a password against policy |
| `--save-config path.yaml` | Save example configuration to file |

### Examples

```bash
# Generate a 20-character password with symbols and strength analysis
./pwgen -length 20 -symbols -strength

# Generate 5 passwords without ambiguous characters
./pwgen -count 5 -no-ambiguous

# Generate password following corporate policy
./pwgen -policy corporate -strength

# Validate existing password against policy
./pwgen -validate "MyP@ssw0rd123" -policy corporate

# List all available policies
./pwgen -list-policies
```

## Password Policies

### Available Templates

- **basic**: Minimal security requirements (8+ chars, mixed case, digits)
- **corporate**: Standard corporate policy (12+ chars, all character types)
- **high-security**: Stringent requirements (16+ chars, multiple of each type)
- **aws**: AWS IAM password policy compliance
- **azure**: Azure AD password complexity requirements
- **pci-dss**: PCI DSS compliant passwords

### Policy Features
- Minimum/maximum length requirements
- Character type requirements (uppercase, lowercase, digits, symbols)
- Minimum count for each character type
- Forbidden character patterns
- Entropy requirements
- Ambiguous character exclusion

## Configuration

### Configuration Files

Create a `.pwgen.yaml` file in your home directory or current directory:

```yaml
length: 16
include_upper: true
include_lower: true
include_digits: true
include_symbols: true
exclude_ambiguous: true
count: 1
show_strength: true
policy_template: "corporate"
```

### Environment Variables

Override settings with environment variables:

```bash
export PWGEN_LENGTH=20
export PWGEN_INCLUDE_SYMBOLS=true
export PWGEN_SHOW_STRENGTH=yes
export PWGEN_POLICY_TEMPLATE=corporate
```

### Configuration Priority

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration files
4. Default values (lowest priority)

## Password Strength Analysis

When using `--strength`, the tool provides:

- **Strength Level**: Very Weak, Weak, Fair, Good, Strong, Very Strong
- **Numeric Score**: 0-100 scale
- **Entropy**: Bits of entropy for cryptographic strength
- **Time to Crack**: Estimated time for brute force attacks
- **Feedback**: Specific recommendations for improvement

Example output:
```bash
./pwgen -strength
C0mpl3x!P@ssw0rd [Strong, Score: 85/100, Entropy: 65.2 bits, Time to crack: 2 million years]
```

## Character Sets

- **Lowercase**: `abcdefghijklmnopqrstuvwxyz`
- **Uppercase**: `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
- **Digits**: `0123456789`
- **Symbols**: `!@#$%^&*()_+-=[]{}|;:,.<>?`
- **Ambiguous**: `0O1lI` (excluded when `--no-ambiguous` is used)

## Requirements

- Go 1.25 or higher

## Security

This tool uses Go's `crypto/rand` package to ensure cryptographically secure password generation. All random number generation is suitable for security-sensitive applications.

### Security Features
- Cryptographically secure random number generation
- Pattern detection to avoid predictable passwords
- Entropy calculation for strength assessment
- Policy validation against common attack vectors

## Development

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd password-generator

# Set up development environment (hooks, tools)
go run tools/setup/main.go all

# Run all verification checks
go run tools/setup/main.go verify
```

### Development Commands

```bash
# Testing
go test -v ./...
go run tools/setup/main.go verify  # Full verification suite

# Building
go build -o pwgen .

# Setup and tools
go run tools/setup/main.go hooks   # Install git hooks
go run tools/setup/main.go tools   # Install dev tools
```

### Git Hooks & Quality

This project uses automated git hooks for code quality:

- **Pre-commit**: Formatting, linting, tests, security checks
- **Pre-push**: Full test suite, coverage (>70%), cross-platform builds  
- **Commit-msg**: Conventional commit format validation

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: add new feature
fix: resolve bug issue  
docs: update documentation
```

### Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and trunk-based workflow.