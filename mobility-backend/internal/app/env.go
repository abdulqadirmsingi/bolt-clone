package app

import (
	"bufio"
	"os"
	"strings"
)

// loadEnvIfExists reads .env in the current directory and sets env vars (only if not already set).
// Lets you run "go run ./cmd/server" without exporting DB_* every time.
func loadEnvIfExists() {
	const envFile = ".env"
	f, err := os.Open(envFile)
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
		i := strings.Index(line, "=")
		if i <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+1:])
		if key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
