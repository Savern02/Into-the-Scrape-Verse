package config

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

// loadDotEnv reads KEY=VALUE lines from a file into the process environment.
//
// Real environment variables always win over the file, so a value exported in
// your shell or injected by a deploy platform overrides .env without you having
// to edit anything. That ordering is what makes the same binary work locally
// and in production.
//
// A missing file is not an error -- in production there is no .env, just real
// environment variables.
func loadDotEnv(path string) {
	if p := os.Getenv("ENV_FILE"); p != "" {
		path = p
	}

	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("could not read env file", "path", path, "err", err)
		}
		return
	}
	defer f.Close()

	loaded := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Tolerate "export KEY=VALUE" so the file can also be sourced by bash.
		line = strings.TrimPrefix(line, "export ")

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if key == "" {
			continue
		}

		// Strip one layer of matching quotes. Unquoted values keep everything
		// after the first '=' intact, which matters: a Postgres URL contains
		// ':' and '/' and sometimes '=' in the password.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if _, exists := os.LookupEnv(key); exists {
			continue // real environment wins
		}
		if err := os.Setenv(key, val); err != nil {
			slog.Warn("could not set env var", "key", key, "err", err)
			continue
		}
		loaded++
	}
	if err := sc.Err(); err != nil {
		slog.Warn("error reading env file", "path", path, "err", err)
		return
	}
	slog.Info("loaded env file", "path", path, "vars", loaded)
}
