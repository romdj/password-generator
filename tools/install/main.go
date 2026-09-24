// Command install copies a locally built pwgen binary to a user-accessible bin directory.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fatal("find home directory", err)
	}

	source := "pwgen"
	if len(os.Args) > 1 {
		source = os.Args[1]
	}

	src, err := os.Open(source)
	if err != nil {
		fatal("open "+source, err)
	}
	defer src.Close()

	destinationDir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		fatal("create "+destinationDir, err)
	}

	destination := filepath.Join(destinationDir, "pwgen")
	dst, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fatal("open "+destination, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		fatal("copy binary", err)
	}
	if err := dst.Close(); err != nil {
		fatal("close "+destination, err)
	}

	fmt.Printf("Installed %s\n", destination)
}

func fatal(action string, err error) {
	fmt.Fprintf(os.Stderr, "install: could not %s: %v\n", action, err)
	os.Exit(1)
}
