package launcher

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
)

// apiVars are the variables that choose how Claude Code authenticates and
// which provider it talks to. An API account drops every inherited one, so
// a key or provider set in the shell cannot leak in or win over its own.
var apiVars = []string{
	"ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN", "CLAUDE_CODE_OAUTH_TOKEN",
	"CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR", "CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR",
	"ANTHROPIC_BASE_URL", "ANTHROPIC_MODEL", "ANTHROPIC_UNIX_SOCKET",
	"CLAUDE_CODE_USE_BEDROCK", "CLAUDE_CODE_USE_VERTEX", "CLAUDE_CODE_USE_FOUNDRY",
	"CLAUDE_CODE_USE_ANTHROPIC_AWS", "CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD",
	"CLAUDE_CODE_USE_MANTLE", "CLAUDE_CODE_USE_GATEWAY",
}

// WithoutAuth is environ without the variables that choose how Claude Code
// signs in or which provider it uses.
func WithoutAuth(environ []string) []string {
	out := make([]string, 0, len(environ))
	for _, kv := range environ {
		name, _, _ := strings.Cut(kv, "=")
		if !contains(apiVars, name) {
			out = append(out, kv)
		}
	}
	return out
}

// keyTimeout bounds a command: key source (a password manager may prompt).
const keyTimeout = 30 * time.Second

// useAPI points the plan of an API account at its endpoint, with the key
// read now from its source. Only the child's environment ever holds it.
//
// A dry run does not read the key: it only checks that an env: source is
// set, and DryRun shows where the key comes from instead of its value.
func (pl *Plan) useAPI(api *config.API, dry bool) error {
	var key string
	var err error
	if !dry {
		key, err = ReadAPIKey(api.KeyFrom)
	} else if name, ok := strings.CutPrefix(api.KeyFrom, "env:"); ok && strings.TrimSpace(os.Getenv(name)) == "" {
		err = fmt.Errorf("the API key comes from $%s, which is empty", name)
	}
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
	keyVar := "ANTHROPIC_API_KEY"
	if api.Auth == config.Bearer {
		keyVar = "ANTHROPIC_AUTH_TOKEN"
	}
	if dry {
		pl.keyVar, pl.keyFrom = keyVar, api.KeyFrom
	} else {
		env = append(env, keyVar+"="+key)
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
	c := shellCommand(ctx, command)
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
