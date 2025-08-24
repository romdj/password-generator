# Password Generator

A secure, configurable password generator written in Go that uses cryptographically secure random number generation.

## Features

- **Configurable Length**: Set password length (default: 12 characters)
- **Character Types**: Choose from uppercase, lowercase, digits, and symbols
- **Ambiguous Character Exclusion**: Option to exclude confusing characters (0, O, 1, l, I)
- **Multiple Passwords**: Generate multiple passwords in one command
- **Secure**: Uses `crypto/rand` for cryptographically secure randomness

## Installation

```bash
git clone <repository-url>
cd password-generator
go build -o pwgen main.go
```

## Usage

### Basic Usage

```bash
# Generate a 12-character password (default)
go run main.go

# Build and use the binary
./pwgen
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

### Examples

```bash
# Generate a 20-character password with symbols
go run main.go -length 20 -symbols

# Generate 5 passwords without ambiguous characters
go run main.go -count 5 -no-ambiguous

# Generate a password with only lowercase and digits
go run main.go -upper=false -symbols=false

# Short flags version
go run main.go -l 16 -s -c 3
```

## Character Sets

- **Lowercase**: `abcdefghijklmnopqrstuvwxyz`
- **Uppercase**: `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
- **Digits**: `0123456789`
- **Symbols**: `!@#$%^&*()_+-=[]{}|;:,.<>?`
- **Ambiguous**: `0O1lI` (excluded when `--no-ambiguous` is used)

## Requirements

- Go 1.16 or higher

## Security

This tool uses Go's `crypto/rand` package to ensure cryptographically secure password generation. All random number generation is suitable for security-sensitive applications.