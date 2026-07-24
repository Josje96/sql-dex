package main

import (
	"bufio"
	"os"
	"strings"
)

// loadDotEnv reads simple KEY=VALUE lines from a .env file and sets any that
// aren't already present in the environment. Missing file is not an error —
// the tutor just stays disabled without a key. This avoids pulling in a
// dependency just to read one file.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`) // drop optional quotes
		if key != "" {
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, val)
			}
		}
	}
}
