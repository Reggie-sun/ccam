package launcher

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// apiVars are the variables that choose how Claude Code authenticates and
// which provider it talks to. An API account drops every inherited one, so
// a key or provider set in the shell cannot leak in or win over its own.
var apiVars = []string{
	"ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN", "CLAUDE_CODE_OAUTH_TOKEN",
	"ANTHROPIC_BASE_URL", "ANTHROPIC_MODEL",
	"CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_VERTEX", "CLAUDE_CODE_USE_FOUNDRY",
}

// keyTimeout bounds a command: key source (a password manager may prompt).
const keyTimeout = 30 * time.Second

// useAPI points the plan of an API account at its endpoint, with the key
// read now from its source. Only the child's environment ever holds it.
func (pl *Plan) useAPI(api *config.API) error {
	key, err := ReadAPIKey(api.KeyFrom)
	if err != nil {
		return fmt.Errorf("account %s: %w", pl.Profile.Name, err)
	}
	env := make([]string, 0, len(pl.Env)+3)
	for _, kv := range pl.Env {
		name, _, _ := strings.Cut(kv, "=")
		if !contains(apiVars, name) {
			env = append(env, kv)
		}
	}
	if api.BaseURL != "" {
		env = append(env, "ANTHROPIC_BASE_URL="+api.BaseURL)
	}
	if api.Model != "" {
		env = append(env, "ANTHROPIC_MODEL="+api.Model)
	}
	if api.Auth == config.Bearer {
		env = append(env, "ANTHROPIC_AUTH_TOKEN="+key)
	} else {
		env = append(env, "ANTHROPIC_API_KEY="+key)
	}
	pl.Env = env
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// ReadAPIKey reads a key from env:NAME or command:COMMAND (run by sh, or cmd on
// Windows; its first line of output is the key).
func ReadAPIKey(from string) (string, error) {
	if name, ok := strings.CutPrefix(from, "env:"); ok {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			return v, nil
		}
		return "", fmt.Errorf("the API key comes from $%s, which is empty", name)
	}
	command, ok := strings.CutPrefix(from, "command:")
	if !ok {
		return "", fmt.Errorf("unknown key source %q", from)
	}
	ctx, cancel := context.WithTimeout(context.Background(), keyTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, "sh", "-c", command)
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx, "cmd", "/C", command)
	}
	c.Stdin, c.Stderr = os.Stdin, os.Stderr // a password manager may ask
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("the API key command failed: %w", err)
	}
	key, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\n")
	if key = strings.TrimSpace(key); key == "" {
		return "", fmt.Errorf("the API key command printed nothing")
	}
	return key, nil
}
