package rf

import (
	"embed"
	"os"
	"strings"
)

//go:embed *.env*
var embededEnvFile embed.FS

var noop = newEnv()

func newEnv() bool {
	file, err := embededEnvFile.ReadFile(".env")
	if err != nil {
		file, err = embededEnvFile.ReadFile(".env.example")
		if err != nil {
			return false
		}
	}

	origEnvMap := make(map[string]bool)
	for _, line := range os.Environ() {
		key, _, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		origEnvMap[key] = true
	}

	for _, line := range strings.Split(string(file), "\n") {
		key, val, found := strings.Cut(line, "=")
		if origEnvMap[key] || !found {
			continue
		}
		os.Setenv(key, strings.ReplaceAll(val, "\"", ""))
	}
	return true
}
