// Package profile manages the per-account CLAUDE_CONFIG_DIR folders.
package profile

import (
	"cmp"
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
	return p.ProfileDir(cmp.Or(prof.Folder, prof.Name))
}

// folderOwner is the account whose managed folder is dir, or nil.
func folderOwner(p paths.Paths, cfg *config.Config, dir string) *config.Profile {
	for i := range cfg.Profiles {
		if q := &cfg.Profiles[i]; q.Dir == "" && Dir(p, q) == dir {
			return q
		}
	}
	return nil
}

// Add registers a new profile in cfg and creates its folder (0700).
// It does not save cfg.
func Add(p paths.Paths, cfg *config.Config, name string) (*config.Profile, error) {
	if err := config.ValidateName(name); err != nil {
		return nil, err
	}
	if _, ok := cfg.Find(name); ok {
		return nil, fmt.Errorf("an account named %q already exists", name)
	}
	dir := p.ProfileDir(name)
	if q := folderOwner(p, cfg, dir); q != nil {
		return nil, fmt.Errorf("folder %s is account %s's (it kept the folder when it was renamed); choose another name", dir, q.Name)
	}
	if fsutil.Exists(dir) {
		return nil, fmt.Errorf("folder already exists: %s (use `ccam import`, or clean up the folder)", dir)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	cfg.DropMember(name)
	cfg.Profiles = append(cfg.Profiles, config.Profile{Name: name, Slot: cfg.NextSlot()})
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
	src := Dir(p, prof)
	if !fsutil.Exists(src) {
		return "", nil
	}
	// Two removals of one name within a second must not share a folder.
	dst := fsutil.UniqueDir(filepath.Join(p.TrashDir(), prof.Name+"-"+fsutil.Stamp(time.Now())))
	return dst, fsutil.Move(src, dst)
}

// Rename renames a profile and its managed folder. With keepFolder the
// folder keeps its name (Profile.Folder): a login store named after the
// folder's path would not follow it.
func Rename(p paths.Paths, cfg *config.Config, oldName, newName string, keepFolder bool) error {
	prof, err := cfg.MustFind(oldName)
	if err != nil {
		return err
	}
	if err := config.ValidateName(newName); err != nil {
		return err
	}
	if _, ok := cfg.Find(newName); ok {
		return fmt.Errorf("an account named %q already exists", newName)
	}
	cfg.DropMember(newName)
	switch oldDir, newDir := Dir(p, prof), p.ProfileDir(newName); {
	case prof.Dir != "":
	case keepFolder:
		prof.Folder = filepath.Base(oldDir)
		if prof.Folder == newName {
			prof.Folder = ""
		}
	default:
		if q := folderOwner(p, cfg, newDir); q != nil && q != prof {
			return fmt.Errorf("folder %s is account %s's (it kept the folder when it was renamed); choose another name", newDir, q.Name)
		}
		// The caller does not save cfg on error, so the folder must keep
		// the name cfg still has.
		if err := MoveFolder(oldDir, newDir); err != nil {
			return err
		}
		prof.Folder = ""
	}
	prof.Name = newName
	if cfg.Default == oldName {
		cfg.Default = newName
	}
	for i := range cfg.Folders {
		if cfg.Folders[i].Account == oldName {
			cfg.Folders[i].Account = newName
		}
	}
	for i := range cfg.Groups {
		for j, m := range cfg.Groups[i].Members {
			if m == oldName {
				cfg.Groups[i].Members[j] = newName
			}
		}
	}
	return nil
}

// MoveFolder moves a managed account folder from oldDir to newDir and points
// the plugin paths recorded through it at newDir. A missing folder is no
// error; when the paths cannot be rewritten, the folder goes back.
func MoveFolder(oldDir, newDir string) error {
	if err := os.Rename(oldDir, newDir); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	sep := string(filepath.Separator)
	if _, err := RewritePluginPaths(filepath.Join(newDir, "plugins"), oldDir+sep, newDir+sep); err != nil {
		if rerr := os.Rename(newDir, oldDir); rerr != nil && !errors.Is(rerr, fs.ErrNotExist) {
			return errors.Join(err, rerr)
		}
		return err
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

// ReadIdentity reads who dir is logged in as from its .claude.json. A
// folder without a login gives a nil identity and no error.
func ReadIdentity(dir string) (*config.Identity, error) {
	var acct struct {
		AccountUUID      string `json:"accountUuid"`
		EmailAddress     string `json:"emailAddress"`
		OrganizationUUID string `json:"organizationUuid"`
		OrganizationName string `json:"organizationName"`
		OrganizationType string `json:"organizationType"`
	}
	ok, err := DecodeFileKey(filepath.Join(dir, ".claude.json"), "oauthAccount", &acct)
	if errors.Is(err, fs.ErrNotExist) || err == nil && (!ok || acct.AccountUUID == "") {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config.Identity{AccountUUID: acct.AccountUUID, Email: acct.EmailAddress, OrgUUID: acct.OrganizationUUID,
		OrgName: acct.OrganizationName, OrgType: acct.OrganizationType}, nil
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
