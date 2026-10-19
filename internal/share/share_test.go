package share

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
)

type fixture struct {
	root, groups string
}

func newFixture(t *testing.T) fixture {
	root := t.TempDir()
	return fixture{root: root, groups: filepath.Join(root, "groups")}
}

func (f fixture) profile(t *testing.T, name string, files map[string]string) string {
	dir := filepath.Join(f.root, "profiles", name)
	for rel, content := range files {
		path := filepath.Join(dir, rel)
		os.MkdirAll(filepath.Dir(path), 0o700)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	os.MkdirAll(dir, 0o700)
	return dir
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func share(t *testing.T, f fixture, dir string, from, to config.Share, opt Options) *Plan {
	t.Helper()
	pl, err := Build(dir, f.groups, from, to, opt)
	if err != nil {
		t.Fatal(err)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	if probs := Check(dir, to, f.groups); len(probs) != 0 {
		t.Fatalf("after apply: %v", probs)
	}
	return pl
}

var allOn = config.Share{Settings: "default", Memory: "default", Extensions: "default", History: "default"}

func TestFirstSharerAdoptsWithoutBackup(t *testing.T) {
	f := newFixture(t)
	dir := f.profile(t, "a", map[string]string{
		"settings.json":       `{"model":"opus"}`,
		"CLAUDE.md":           "# mine",
		"skills/x/SKILL.md":   "skill",
		"projects/p/s1.jsonl": "{}",
		"history.jsonl":       "{\"display\":\"hi\"}\n",
		".credentials.json":   "secret",
		".claude.json":        "{}",
	})
	share(t, f, dir, config.Share{}, allOn, Options{})
	if read(t, filepath.Join(dir, "settings.json")) != `{"model":"opus"}` {
		t.Fatal("settings not readable through link")
	}
	if read(t, filepath.Join(f.groups, "default", "skills", "x", "SKILL.md")) != "skill" {
		t.Fatal("skills not in group")
	}
	for _, private := range []string{".credentials.json", ".claude.json"} {
		info, err := os.Lstat(filepath.Join(dir, private))
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("%s must stay a private regular file", private)
		}
		if _, err := os.Stat(filepath.Join(f.groups, "default", private)); err == nil {
			t.Fatalf("%s leaked into the group", private)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, BackupDir)); err == nil {
		t.Fatal("nothing should need a backup for the first sharer")
	}
}

// TestRoundTripKeepsData: on → off → on must never lose anything.
func TestRoundTripKeepsData(t *testing.T) {
	f := newFixture(t)
	a := f.profile(t, "a", map[string]string{"settings.json": `{"a":1}`, "projects/p/s1.jsonl": "one"})
	b := f.profile(t, "b", map[string]string{"projects/p/s2.jsonl": "two", "projects/q/s3.jsonl": "three"})

	share(t, f, a, config.Share{}, allOn, Options{})
	share(t, f, b, config.Share{}, allOn, Options{})
	for _, s := range []string{"p/s1.jsonl", "p/s2.jsonl", "q/s3.jsonl"} {
		if _, err := os.Stat(filepath.Join(a, "projects", s)); err != nil {
			t.Fatalf("a cannot see %s after b joined", s)
		}
	}

	// Off: b gets its own copy, the group keeps everything.
	share(t, f, b, allOn, config.Share{}, Options{})
	if read(t, filepath.Join(b, "projects", "p", "s1.jsonl")) != "one" {
		t.Fatal("copy back missing data")
	}
	if info, _ := os.Lstat(filepath.Join(b, "projects")); info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("still linked after off")
	}
	os.WriteFile(filepath.Join(b, "projects", "p", "s4.jsonl"), []byte("four"), 0o600)

	// On again: the new local session joins the group.
	share(t, f, b, config.Share{}, allOn, Options{})
	for _, s := range []string{"s1", "s2", "s4"} {
		if _, err := os.Stat(filepath.Join(a, "projects", "p", s+".jsonl")); err != nil {
			t.Fatalf("lost %s", s)
		}
	}
}

func TestFreshOffStartsEmpty(t *testing.T) {
	f := newFixture(t)
	a := f.profile(t, "a", map[string]string{"CLAUDE.md": "shared"})
	on := config.Share{Memory: "default"}
	share(t, f, a, config.Share{}, on, Options{})
	share(t, f, a, on, config.Share{}, Options{Fresh: true})
	if _, err := os.Lstat(filepath.Join(a, "CLAUDE.md")); err == nil {
		t.Fatal("fresh should leave no CLAUDE.md")
	}
	if read(t, filepath.Join(f.groups, "default", "CLAUDE.md")) != "shared" {
		t.Fatal("group content must stay")
	}
}

func TestConflictResolutions(t *testing.T) {
	cases := []struct {
		name   string
		file   string
		group  string
		mine   string
		choice Resolution
		want   string
	}{
		{"keep", "CLAUDE.md", "G", "M", KeepGroup, "G"},
		{"mine", "CLAUDE.md", "G", "M", UseMine, "M"},
		{"lines", "history.jsonl", "1\n2\n", "2\n3\n", MergeLines, "1\n2\n3\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newFixture(t)
			os.MkdirAll(filepath.Join(f.groups, "default"), 0o700)
			os.WriteFile(filepath.Join(f.groups, "default", c.file), []byte(c.group), 0o600)
			dir := f.profile(t, "a", map[string]string{c.file: c.mine})
			on := config.Share{Memory: "default", History: "default"}
			pl, err := Build(dir, f.groups, config.Share{}, on, Options{})
			if err != nil || len(pl.Conflicts) != 1 {
				t.Fatalf("conflicts %v %v", pl, err)
			}
			pl.Conflicts[0].Choice = c.choice
			if err := pl.Apply(); err != nil {
				t.Fatal(err)
			}
			if got := read(t, filepath.Join(dir, c.file)); got != c.want {
				t.Fatalf("got %q want %q", got, c.want)
			}
			backups, _ := filepath.Glob(filepath.Join(dir, BackupDir, "*", c.file))
			if len(backups) != 1 || read(t, backups[0]) != c.mine {
				t.Fatalf("original not backed up: %v", backups)
			}
		})
	}
}

func TestMergeDirKeepsConflictingFileInBackup(t *testing.T) {
	f := newFixture(t)
	os.MkdirAll(filepath.Join(f.groups, "default", "projects", "p"), 0o700)
	os.WriteFile(filepath.Join(f.groups, "default", "projects", "p", "s.jsonl"), []byte("group"), 0o600)
	dir := f.profile(t, "a", map[string]string{"projects/p/s.jsonl": "mine", "projects/p/new.jsonl": "new"})
	share(t, f, dir, config.Share{}, config.Share{History: "default"}, Options{})
	if read(t, filepath.Join(dir, "projects", "p", "new.jsonl")) != "new" {
		t.Fatal("new file not merged")
	}
	if read(t, filepath.Join(dir, "projects", "p", "s.jsonl")) != "group" {
		t.Fatal("group copy should win")
	}
	backups, _ := filepath.Glob(filepath.Join(dir, BackupDir, "*", "projects", "p", "s.jsonl"))
	if len(backups) != 1 || read(t, backups[0]) != "mine" {
		t.Fatal("my conflicting file must be in the backup")
	}
}

func TestSwitchGroupRelinks(t *testing.T) {
	f := newFixture(t)
	dir := f.profile(t, "a", map[string]string{"CLAUDE.md": "x"})
	share(t, f, dir, config.Share{}, config.Share{Memory: "default"}, Options{})
	pl := share(t, f, dir, config.Share{Memory: "default"}, config.Share{Memory: "team"}, Options{})
	if len(pl.Notes) == 0 {
		t.Fatal("switching groups should explain what happens to the old group")
	}
	if target, _ := os.Readlink(filepath.Join(dir, "CLAUDE.md")); target != filepath.Join(f.groups, "team", "CLAUDE.md") {
		t.Fatalf("link %s", target)
	}
}

func TestCheckAndQuickFix(t *testing.T) {
	f := newFixture(t)
	dir := f.profile(t, "a", map[string]string{"CLAUDE.md": "real file"})
	on := config.Share{Memory: "default", Extensions: "default"}
	probs := Check(dir, on, f.groups)
	kinds := map[ProblemKind]int{}
	for _, p := range probs {
		kinds[p.Kind]++
	}
	if kinds[NotLinked] != 1 || kinds[MissingLink] != 6 {
		t.Fatalf("problems %v", probs)
	}
	rest, err := QuickFix(probs)
	if err != nil || len(rest) != 1 || rest[0].Kind != NotLinked {
		t.Fatalf("quickfix left %v %v", rest, err)
	}
	if info, err := os.Stat(filepath.Join(dir, "skills")); err != nil || !info.IsDir() {
		t.Fatal("skills link should resolve to a dir")
	}
	// Local item still linked → StaleLink.
	stale := Check(dir, config.Share{Memory: "default"}, f.groups)
	if len(stale) != 7 { // 1 NotLinked + 6 extensions stale
		t.Fatalf("stale %v", stale)
	}
}

func TestMergeJSON(t *testing.T) {
	base := []byte(`{"model":"opus","permissions":{"allow":["Read","Edit"]},"env":{"A":"1"}}`)
	mine := []byte(`{"model":"sonnet","permissions":{"allow":["Edit","Bash(git *)"],"deny":["WebFetch"]},"theme":"dark"}`)
	out, err := MergeJSONBytes(base, mine)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	json.Unmarshal(out, &got)
	want := map[string]any{
		"model":       "opus",
		"theme":       "dark",
		"env":         map[string]any{"A": "1"},
		"permissions": map[string]any{"allow": []any{"Read", "Edit", "Bash(git *)"}, "deny": []any{"WebFetch"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s", out)
	}
	if _, err := MergeJSONBytes([]byte("{"), []byte("{}")); err == nil {
		t.Fatal("bad json accepted")
	}
}

func TestDescribe(t *testing.T) {
	f := newFixture(t)
	dir := f.profile(t, "a", map[string]string{"CLAUDE.md": "x"})
	pl, _ := Build(dir, f.groups, config.Share{}, config.Share{Memory: "default"}, Options{})
	var buf bytes.Buffer
	pl.Describe(&buf)
	if !bytes.Contains(buf.Bytes(), []byte("move")) || !bytes.Contains(buf.Bytes(), []byte("link")) {
		t.Fatalf("describe:\n%s", buf.String())
	}
}

func TestOverwritingGroupFileBacksItUp(t *testing.T) {
	for _, choice := range []Resolution{UseMine, MergeJSON} {
		f := newFixture(t)
		os.MkdirAll(filepath.Join(f.groups, "default"), 0o700)
		os.WriteFile(filepath.Join(f.groups, "default", "settings.json"), []byte(`{"a":1}`), 0o600)
		dir := f.profile(t, "b", map[string]string{"settings.json": `{"b":2}`})
		pl, _ := Build(dir, f.groups, config.Share{}, config.Share{Settings: "default"}, Options{})
		pl.Conflicts[0].Choice = choice
		if err := pl.Apply(); err != nil {
			t.Fatal(err)
		}
		saved, _ := filepath.Glob(filepath.Join(dir, BackupDir, "*", "group", "settings.json"))
		if len(saved) != 1 || read(t, saved[0]) != `{"a":1}` {
			t.Fatalf("%v: group copy not backed up: %v", choice, saved)
		}
	}
}

func TestForeignLinkContentIsAdopted(t *testing.T) {
	f := newFixture(t)
	dotfiles := filepath.Join(f.root, "dotfiles")
	os.MkdirAll(dotfiles, 0o700)
	os.WriteFile(filepath.Join(dotfiles, "CLAUDE.md"), []byte("from dotfiles"), 0o600)
	dir := f.profile(t, "a", nil)
	os.Symlink(filepath.Join(dotfiles, "CLAUDE.md"), filepath.Join(dir, "CLAUDE.md"))
	share(t, f, dir, config.Share{}, config.Share{Memory: "default"}, Options{})
	if read(t, filepath.Join(dir, "CLAUDE.md")) != "from dotfiles" {
		t.Fatal("dotfiles content must become the group's")
	}
	if read(t, filepath.Join(dotfiles, "CLAUDE.md")) != "from dotfiles" {
		t.Fatal("dotfiles file must stay")
	}
	if links, _ := filepath.Glob(filepath.Join(dir, BackupDir, "*", "CLAUDE.md")); len(links) != 1 {
		t.Fatal("the old link should be kept in the backup")
	}
}

// A second account joining extension sharing keeps its plugins: the plugin
// registries are merged and its paths point at the group folder.
func TestJoiningMergesPluginRegistry(t *testing.T) {
	f := newFixture(t)
	ext := config.Share{Extensions: "default"}
	a := f.profile(t, "a", nil)
	os.MkdirAll(filepath.Join(a, "plugins", "cache", "mk", "one"), 0o700)
	os.WriteFile(filepath.Join(a, "plugins", "installed_plugins.json"),
		[]byte(`{"version":2,"plugins":{"one@mk":[{"installPath":"`+a+`/plugins/cache/mk/one"}]}}`), 0o600)
	share(t, f, a, config.Share{}, ext, Options{})

	b := f.profile(t, "b", nil)
	os.MkdirAll(filepath.Join(b, "plugins", "cache", "mk", "two"), 0o700)
	os.WriteFile(filepath.Join(b, "plugins", "installed_plugins.json"),
		[]byte(`{"version":2,"plugins":{"two@mk":[{"installPath":"`+b+`/plugins/cache/mk/two"}]}}`), 0o600)
	share(t, f, b, config.Share{}, ext, Options{})

	grp := filepath.Join(f.groups, "default", "plugins")
	data := read(t, filepath.Join(grp, "installed_plugins.json"))
	for _, want := range []string{`"one@mk"`, `"two@mk"`, grp + "/cache/mk/one", grp + "/cache/mk/two"} {
		if !strings.Contains(data, want) {
			t.Errorf("registry lacks %s:\n%s", want, data)
		}
	}
	if strings.Contains(data, a+"/plugins") || strings.Contains(data, b+"/plugins") {
		t.Errorf("paths still run through an account folder:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(grp, "cache", "mk", "two")); err != nil {
		t.Fatal("b's plugin files not in the group")
	}
	if saved, _ := filepath.Glob(filepath.Join(b, BackupDir, "*", "group", "plugins", "installed_plugins.json")); len(saved) != 1 {
		t.Fatal("group registry must be backed up before merging")
	}
}
