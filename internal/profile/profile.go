// Package profile manages the per-account CLAUDE_CONFIG_DIR folders.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
	"github.com/jhyunwoo/ccam/internal/paths"
)

// Dir is the CLAUDE_CONFIG_DIR of a profile.
func Dir(p paths.Paths, prof *config.Profile) string {
	if prof.Dir != "" {
		return prof.Dir
	}
	return p.ProfileDir(prof.Name)
}

// Add registers a new profile in cfg and creates its folder (0700).
// It does not save cfg.
func Add(p paths.Paths, cfg *config.Config, name string, share config.Share) (*config.Profile, error) {
	if err := config.ValidateName(name); err != nil {
		return nil, err
	}
	if _, ok := cfg.Find(name); ok {
		return nil, fmt.Errorf("이미 있는 계정 이름입니다: %q", name)
	}
	dir := p.ProfileDir(name)
	if fsutil.Exists(dir) {
		return nil, fmt.Errorf("폴더가 이미 있습니다: %s (`ccam import` 를 쓰거나 폴더를 정리하세요)", dir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: cfg.NextSlot(), Share: share})
	if len(cfg.Profiles) == 1 {
		cfg.Default = name
	}
	return &cfg.Profiles[len(cfg.Profiles)-1], nil
}

// Trash moves a managed profile folder into the trash. Imported folders
// (registered in place) are left untouched; only the registration goes away.
func Trash(p paths.Paths, prof *config.Profile) (string, error) {
	if prof.Dir != "" {
		return "", nil
	}
	src := p.ProfileDir(prof.Name)
	if !fsutil.Exists(src) {
		return "", nil
	}
	dst := filepath.Join(p.TrashDir(), prof.Name+"-"+fsutil.Stamp(time.Now()))
	return dst, fsutil.Move(src, dst)
}

// Rename renames a profile and its managed folder.
func Rename(p paths.Paths, cfg *config.Config, oldName, newName string) error {
	prof, err := cfg.MustFind(oldName)
	if err != nil {
		return err
	}
	if err := config.ValidateName(newName); err != nil {
		return err
	}
	if _, ok := cfg.Find(newName); ok {
		return fmt.Errorf("이미 있는 계정 이름입니다: %q", newName)
	}
	if prof.Dir == "" {
		if err := os.Rename(p.ProfileDir(oldName), p.ProfileDir(newName)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		sep := string(filepath.Separator)
		if _, err := RewritePluginPaths(filepath.Join(p.ProfileDir(newName), "plugins"),
			p.ProfileDir(oldName)+sep, p.ProfileDir(newName)+sep); err != nil {
			return err
		}
	}
	prof.Name = newName
	if cfg.Default == oldName {
		cfg.Default = newName
	}
	return nil
}

// ReadEmail pulls oauthAccount.emailAddress out of <dir>/.claude.json.
// The file can be several MB, so only the one key is decoded.
func ReadEmail(dir string) (string, error) {
	var acct struct {
		EmailAddress string `json:"emailAddress"`
	}
	if _, err := DecodeFileKey(filepath.Join(dir, ".claude.json"), "oauthAccount", &acct); err != nil {
		return "", err
	}
	return acct.EmailAddress, nil
}

// DecodeFileKey decodes one top-level key of a JSON file into v.
func DecodeFileKey(path, key string, v any) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	raw, ok, err := TopLevelValue(data, key)
	if err != nil || !ok {
		return false, err
	}
	return true, json.Unmarshal(raw, v)
}
