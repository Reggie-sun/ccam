package share

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/jhyunwoo/ccam/internal/config"
	"github.com/jhyunwoo/ccam/internal/fsutil"
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
	dir := f.profile(t, "a", map[string]string{"CLAUDE.md": "x", "skills/s/SKILL.md": "s"})
	on := config.Share{Memory: "default", Extensions: "default"}
	share(t, f, dir, config.Share{}, on, Options{})
	// A new group starts with a copy of what the account saw.
	share(t, f, dir, on, config.Share{Memory: "team", Extensions: "team"}, Options{})
	if target, _ := os.Readlink(filepath.Join(dir, "CLAUDE.md")); target != filepath.Join(f.groups, "team", "CLAUDE.md") {
		t.Fatalf("link %s", target)
	}
	if read(t, filepath.Join(f.groups, "team", "CLAUDE.md")) != "x" || read(t, filepath.Join(dir, "skills", "s", "SKILL.md")) != "s" {
		t.Fatal("team should start with default's content")
	}
	if read(t, filepath.Join(f.groups, "default", "CLAUDE.md")) != "x" {
		t.Fatal("default must stay as it is")
	}
	// A group that has content is linked as it is.
	os.MkdirAll(filepath.Join(f.groups, "pair"), 0o700)
	os.WriteFile(filepath.Join(f.groups, "pair", "CLAUDE.md"), []byte("pair"), 0o600)
	pl := share(t, f, dir, config.Share{Memory: "team"}, config.Share{Memory: "pair"}, Options{})
	if len(pl.Notes) == 0 {
		t.Fatal("switching groups should explain what happens to the old group")
	}
	if read(t, filepath.Join(dir, "CLAUDE.md")) != "pair" {
		t.Fatal("link to pair")
	}
}

// Plans built one after another see what the earlier ones put into the group.
func TestPendingSeesEarlierPlans(t *testing.T) {
	f := newFixture(t)
	a := f.profile(t, "a", map[string]string{"CLAUDE.md": "a", "skills/x/SKILL.md": "x"})
	b := f.profile(t, "b", map[string]string{"CLAUDE.md": "b", "skills/y/SKILL.md": "y"})
	on := config.Share{Memory: "default", Extensions: "default"}
	pa, err := Build(a, f.groups, config.Share{}, on, Options{})
	if err != nil {
		t.Fatal(err)
	}
	pb, err := Build(b, f.groups, config.Share{}, on, Options{Pending: pa.Created()})
	if err != nil {
		t.Fatal(err)
	}
	if len(pb.Conflicts) != 1 || pb.Conflicts[0].Seed != filepath.Join(a, "CLAUDE.md") {
		t.Fatalf("b's CLAUDE.md should conflict with a's: %+v", pb.Conflicts)
	}
	for _, s := range pb.Steps {
		if s.Op == OpAdopt {
			t.Fatalf("b must not adopt %s: a does", s.Entry.Name)
		}
	}
	if err := pa.Apply(); err != nil {
		t.Fatal(err)
	}
	if err := pb.Apply(); err != nil {
		t.Fatal(err)
	}
	if read(t, filepath.Join(b, "CLAUDE.md")) != "a" || read(t, filepath.Join(b, "skills", "x", "SKILL.md")) != "x" || read(t, filepath.Join(a, "skills", "y", "SKILL.md")) != "y" {
		t.Fatal("both should see the merged group")
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
		[]byte(`{"version":2,"plugins":{"one@mk":[{"installPath":"`+jp(a, "plugins/cache/mk/one")+`"}]}}`), 0o600)
	share(t, f, a, config.Share{}, ext, Options{})

	b := f.profile(t, "b", nil)
	os.MkdirAll(filepath.Join(b, "plugins", "cache", "mk", "two"), 0o700)
	os.WriteFile(filepath.Join(b, "plugins", "installed_plugins.json"),
		[]byte(`{"version":2,"plugins":{"two@mk":[{"installPath":"`+jp(b, "plugins/cache/mk/two")+`"}]}}`), 0o600)
	share(t, f, b, config.Share{}, ext, Options{})

	grp := filepath.Join(f.groups, "default", "plugins")
	data := read(t, filepath.Join(grp, "installed_plugins.json"))
	for _, want := range []string{`"one@mk"`, `"two@mk"`, jp(grp, "cache/mk/one"), jp(grp, "cache/mk/two")} {
		if !strings.Contains(data, want) {
			t.Errorf("registry lacks %s:\n%s", want, data)
		}
	}
	if strings.Contains(data, jp(a, "plugins")) || strings.Contains(data, jp(b, "plugins")) {
		t.Errorf("paths still run through an account folder:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(grp, "cache", "mk", "two")); err != nil {
		t.Fatal("b's plugin files not in the group")
	}
	if saved, _ := filepath.Glob(filepath.Join(b, BackupDir, "*", "group", "plugins", "installed_plugins.json")); len(saved) != 1 {
		t.Fatal("group registry must be backed up before merging")
	}
}

// A relative link into dotfiles (GNU stow) still resolves once its folder
// is moved into the group, whether it adopts the group slot or merges in.
func TestRelativeOutsideLinksSurviveSharing(t *testing.T) {
	f := newFixture(t)
	dot := filepath.Join(f.root, "dotfiles", "skill")
	os.MkdirAll(dot, 0o700)
	os.WriteFile(filepath.Join(dot, "SKILL.md"), []byte("stowed"), 0o600)
	ext := config.Share{Extensions: "default"}
	for _, name := range []string{"a", "b"} { // a adopts, b merges
		// An imported ~/.claude sits deeper than the group folder.
		dir := filepath.Join(f.root, "home", name, ".claude")
		os.MkdirAll(filepath.Join(dir, "skills"), 0o700)
		os.Symlink("../../../../dotfiles/skill", filepath.Join(dir, "skills", "stow-"+name))
		share(t, f, dir, config.Share{}, ext, Options{})
		if got := read(t, filepath.Join(f.groups, "default", "skills", "stow-"+name, "SKILL.md")); got != "stowed" {
			t.Fatalf("%s: %q", name, got)
		}
	}
}

// Turning extensions off copies the plugins back; their registry must then
// point at the account's own copy, not at the group.
func TestLeavingRewritesPluginPathsBack(t *testing.T) {
	f := newFixture(t)
	ext := config.Share{Extensions: "default"}
	a := f.profile(t, "a", nil)
	os.MkdirAll(filepath.Join(a, "plugins", "cache", "mk", "one"), 0o700)
	os.WriteFile(filepath.Join(a, "plugins", "installed_plugins.json"),
		[]byte(`{"plugins":{"one@mk":[{"installPath":"`+jp(a, "plugins/cache/mk/one")+`"}]}}`), 0o600)
	share(t, f, a, config.Share{}, ext, Options{})
	share(t, f, a, ext, config.Share{}, Options{})
	data := read(t, filepath.Join(a, "plugins", "installed_plugins.json"))
	if !strings.Contains(data, jp(a, "plugins/cache/mk/one")) || strings.Contains(data, jp(f.groups, "")) {
		t.Fatalf("registry after leaving:\n%s", data)
	}
}

// The same plugin installed in two accounts stays one install record.
func TestJoiningDoesNotDuplicateSharedPlugin(t *testing.T) {
	f := newFixture(t)
	ext := config.Share{Extensions: "default"}
	for _, name := range []string{"a", "b"} {
		dir := f.profile(t, name, nil)
		os.MkdirAll(filepath.Join(dir, "plugins", "cache", "mk", "one"), 0o700)
		os.WriteFile(filepath.Join(dir, "plugins", "installed_plugins.json"),
			[]byte(`{"plugins":{"one@mk":[{"installPath":"`+jp(dir, "plugins/cache/mk/one")+`","version":"1.0"}]}}`), 0o600)
		share(t, f, dir, config.Share{}, ext, Options{})
	}
	var reg struct {
		Plugins map[string][]any `json:"plugins"`
	}
	json.Unmarshal([]byte(read(t, filepath.Join(f.groups, "default", "plugins", "installed_plugins.json"))), &reg)
	if n := len(reg.Plugins["one@mk"]); n != 1 {
		t.Fatalf("%d install records: %v", n, reg.Plugins)
	}
}

// Two plans in the same second must not share a backup folder: the second
// would overwrite (or fail to move over) the first one's backup.
func TestBackupsWithinOneSecondStayApart(t *testing.T) {
	f := newFixture(t)
	f.profile(t, "g", map[string]string{"settings.json": `{"group":1}`})
	share(t, f, filepath.Join(f.root, "profiles", "g"), config.Share{}, config.Share{Settings: "default"}, Options{})
	dir := f.profile(t, "a", map[string]string{"settings.json": `{"mine":1}`})
	on, off := config.Share{Settings: "default"}, config.Share{}
	share(t, f, dir, off, on, Options{})
	share(t, f, dir, on, off, Options{})
	os.WriteFile(filepath.Join(dir, "settings.json"), []byte(`{"mine":2}`), 0o600)
	share(t, f, dir, off, on, Options{})
	saved, _ := filepath.Glob(filepath.Join(dir, BackupDir, "*", "settings.json"))
	var got []string
	for _, s := range saved {
		got = append(got, read(t, s))
	}
	if len(got) != 2 || !strings.Contains(strings.Join(got, ""), `{"mine":1}`) {
		t.Fatalf("backups: %v", got)
	}
}

// jp is filepath.Join(base, rel) as it appears inside a JSON string, where
// Windows backslashes are escaped.
func jp(base, rel string) string {
	b, _ := json.Marshal(filepath.Join(base, filepath.FromSlash(rel)))
	return strings.Trim(string(b), `"`)
}

// A file one earlier plan of a change brings into a folder another earlier
// plan moves in is there, whatever order the Pending map is read in.
func TestPendingFileInsidePendingFolder(t *testing.T) {
	f := newFixture(t)
	a := f.profile(t, "a", nil)
	os.MkdirAll(filepath.Join(a, "plugins"), 0o700) // no plugin list yet
	b := f.profile(t, "b", map[string]string{"plugins/installed_plugins.json": `{}`})
	grp := filepath.Join(f.groups, "g", "plugins")
	opt := Options{Pending: map[string]string{
		grp: filepath.Join(a, "plugins"),
		filepath.Join(grp, "installed_plugins.json"): filepath.Join(b, "plugins", "installed_plugins.json"),
	}}
	for i := range 100 {
		if !opt.exists(filepath.Join(grp, "installed_plugins.json")) {
			t.Fatalf("call %d: the list b brings is not seen", i)
		}
	}
	if opt.exists(filepath.Join(grp, "known_marketplaces.json")) {
		t.Fatal("a file no plan brings is seen")
	}
}

// Files inside a shared folder that the group also has with other content
// stay the group's: the plan says which, and where the account's copies go.
// Here b's plan sees the group's copy only through a's plan (Pending).
func TestMergeNamesTheFilesThatDiffer(t *testing.T) {
	f := newFixture(t)
	hist := config.Share{History: "hist"}
	a := f.profile(t, "a", map[string]string{"projects/-w/memory/MEMORY.md": "memory of a", "projects/-w/s.jsonl": "same"})
	b := f.profile(t, "b", map[string]string{"projects/-w/memory/MEMORY.md": "memory of b", "projects/-w/s.jsonl": "same"})
	pa, err := Build(a, f.groups, config.Share{}, hist, Options{})
	if err != nil {
		t.Fatal(err)
	}
	pb, err := Build(b, f.groups, config.Share{}, hist, Options{Pending: pa.Created()})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	pb.Describe(&buf)
	want := "projects: 1 file differs from the group's copy (" + filepath.Join("-w", "memory", "MEMORY.md") + "); the group's is kept, the account's stays in " + filepath.Join(b, BackupDir, pb.stamp, "projects")
	if !strings.Contains(buf.String(), want) || strings.Contains(buf.String(), "s.jsonl") {
		t.Fatalf("b's plan:\n%s", buf.String())
	}
	if err := pa.Apply(); err != nil {
		t.Fatal(err)
	}
	if err := pb.Apply(); err != nil {
		t.Fatal(err)
	}
	if got := read(t, filepath.Join(b, BackupDir, pb.stamp, "projects", "-w", "memory", "MEMORY.md")); got != "memory of b" {
		t.Fatalf("b's memory in the backup: %q", got)
	}

	// A conversation carried on while not sharing is named too.
	share(t, f, b, hist, config.Share{}, Options{})
	os.WriteFile(filepath.Join(b, "projects", "-w", "s.jsonl"), []byte("same\nmore"), 0o600)
	pl, err := Build(b, f.groups, config.Share{}, hist, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.ContainsFunc(pl.Notes, func(n string) bool { return strings.Contains(n, filepath.Join("-w", "s.jsonl")) }) {
		t.Fatalf("notes: %v", pl.Notes)
	}
}

// Leaving and joining again copies the group's files back and forth; the
// copies that are still the group's byte for byte are not piled up in the
// backup folder.
func TestRoundTripLeavesNoIdenticalBackup(t *testing.T) {
	f := newFixture(t)
	hist := config.Share{History: "hist"}
	b := f.profile(t, "b", map[string]string{"projects/-w/big.jsonl": "same", "projects/-w/memory/MEMORY.md": "m"})
	share(t, f, b, config.Share{}, hist, Options{})
	share(t, f, b, hist, config.Share{}, Options{})
	pl := share(t, f, b, config.Share{}, hist, Options{})
	if len(pl.Notes) != 0 {
		t.Fatalf("notes: %v", pl.Notes)
	}
	var left []string
	filepath.WalkDir(filepath.Join(b, BackupDir), func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			left = append(left, path)
		}
		return nil
	})
	if len(left) != 0 {
		t.Errorf("identical copies backed up: %v", left)
	}
	if read(t, filepath.Join(b, "projects", "-w", "big.jsonl")) != "same" {
		t.Fatal("b lost its conversation")
	}
}

// perms is whether file permissions can make a copy fail (not as root, not
// on Windows).
var perms = runtime.GOOS != "windows" && os.Geteuid() != 0

// skillsFixture shares five skills of account a with group default and
// makes the third unreadable, so a copy of the group's skills stops there.
func skillsFixture(t *testing.T) (f fixture, a, locked string) {
	t.Helper()
	if !perms {
		t.Skip("needs file permissions")
	}
	f = newFixture(t)
	files := map[string]string{}
	for i := 1; i <= 5; i++ {
		files[fmt.Sprintf("skills/s%d/SKILL.md", i)] = fmt.Sprint(i)
	}
	a = f.profile(t, "a", files)
	share(t, f, a, config.Share{}, config.Share{Extensions: "default"}, Options{})
	locked = filepath.Join(f.groups, "default", "skills", "s3", "SKILL.md")
	os.Chmod(locked, 0)
	t.Cleanup(func() { os.Chmod(locked, 0o600) })
	return f, a, locked
}

func allSkills(t *testing.T, dir string) {
	t.Helper()
	for i := 1; i <= 5; i++ {
		if got := read(t, filepath.Join(dir, fmt.Sprintf("s%d", i), "SKILL.md")); got != fmt.Sprint(i) {
			t.Fatalf("s%d: %q", i, got)
		}
	}
}

// A leave that stops while copying the group's content back keeps the
// link, so the plan doctor --fix builds again copies everything; a copy
// that stopped halfway is never taken for the account's own.
func TestStoppedCopyBackKeepsTheLink(t *testing.T) {
	f, a, locked := skillsFixture(t)
	ext := config.Share{Extensions: "default"}
	pl, err := Build(a, f.groups, ext, config.Share{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := pl.Apply(); err == nil {
		t.Fatal("the copy did not stop")
	}
	os.Chmod(locked, 0o600)
	if target, ok := fsutil.LinkTarget(filepath.Join(a, "skills")); !ok || target != filepath.Join(f.groups, "default", "skills") {
		t.Fatalf("skills is no longer the link: %q", target)
	}
	if left, _ := os.ReadDir(filepath.Join(a, BackupDir)); len(left) != 0 {
		t.Fatalf("left in the backup folder: %v", left)
	}
	share(t, f, a, config.Share{}, config.Share{}, Options{})
	allSkills(t, filepath.Join(a, "skills"))
	if fsutil.Exists(filepath.Join(a, BackupDir)) {
		t.Fatal("the copy left a backup folder")
	}
}

// The same for an account moving to a new group, whose folder starts as a
// copy of the old one's.
func TestStoppedCopyIntoAGroupIsCopiedAgain(t *testing.T) {
	f, a, locked := skillsFixture(t)
	team := config.Share{Extensions: "team"}
	pl, err := Build(a, f.groups, config.Share{Extensions: "default"}, team, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := pl.Apply(); err == nil {
		t.Fatal("the copy did not stop")
	}
	os.Chmod(locked, 0o600)
	if fsutil.Exists(filepath.Join(f.groups, "team", "skills")) {
		t.Fatal("a copy that stopped halfway is in the group")
	}
	share(t, f, a, team, team, Options{})
	allSkills(t, filepath.Join(f.groups, "team", "skills"))
	allSkills(t, filepath.Join(a, "skills"))
}

// A plan removes the link at a path as it found it when it was made: what
// took its place since (an editor saving settings.json there) stays.
func TestUnlinkLeavesWhatReplacedTheLink(t *testing.T) {
	f := newFixture(t)
	on := config.Share{Settings: "default"}
	a := f.profile(t, "a", map[string]string{"settings.json": `{"group":1}`})
	share(t, f, a, config.Share{}, on, Options{})
	settings := filepath.Join(a, "settings.json")
	for _, fresh := range []bool{true, false} {
		pl, err := Build(a, f.groups, on, config.Share{}, Options{Fresh: fresh})
		if err != nil {
			t.Fatal(err)
		}
		os.Remove(settings)
		os.WriteFile(settings, []byte(`{"mine":1}`), 0o600)
		if err := pl.Apply(); err == nil {
			t.Errorf("fresh %v: no error", fresh)
		}
		if got := read(t, settings); got != `{"mine":1}` {
			t.Fatalf("fresh %v: the file written meanwhile reads %q", fresh, got)
		}
		os.Remove(settings)
		if err := os.Symlink(filepath.Join(f.groups, "default", "settings.json"), settings); err != nil {
			t.Fatal(err)
		}
	}
}

// A link to the right group folder that is gone (removed by hand or by a
// cleanup tool) leaves Claude Code nowhere to save into: CheckTargets
// reports it and QuickFix makes the folder again. A file link to a file
// not written yet is normal; only a group folder that is gone counts for
// files. Check, which a launch runs, spends no lstat on it.
func TestCheckFindsAMissingGroupFolder(t *testing.T) {
	f := newFixture(t)
	dir := f.profile(t, "a", nil)
	share(t, f, dir, config.Share{}, allOn, Options{}) // CLAUDE.md is not written yet
	for _, name := range []string{"skills", "projects"} {
		if err := os.Remove(filepath.Join(f.groups, "default", name)); err != nil {
			t.Fatal(err)
		}
	}
	if probs := Check(dir, allOn, f.groups); len(probs) != 0 {
		t.Fatalf("Check: %v", probs)
	}
	probs := CheckTargets(dir, allOn, f.groups)
	if len(probs) != 2 || probs[0].Kind != MissingTarget || probs[1].Kind != MissingTarget || !strings.Contains(probs[0].String(), "missing") {
		t.Fatalf("problems %v", probs)
	}
	if rest, err := QuickFix(probs); err != nil || len(rest) != 0 {
		t.Fatalf("quickfix left %v %v", rest, err)
	}
	if probs := CheckTargets(dir, allOn, f.groups); len(probs) != 0 {
		t.Fatalf("after quickfix: %v", probs)
	}

	// The whole group folder gone: the file links lead nowhere either.
	os.RemoveAll(filepath.Join(f.groups, "default"))
	probs = CheckTargets(dir, allOn, f.groups)
	if len(probs) != 14 {
		t.Fatalf("problems %v", probs)
	}
	if _, err := QuickFix(probs); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("x"), 0o600); err != nil {
		t.Fatalf("writing through the link: %v", err)
	}
	if probs := CheckTargets(dir, allOn, f.groups); len(probs) != 0 {
		t.Fatalf("after quickfix: %v", probs)
	}
}

// When the finished copy cannot take the link's place (on Windows a virus
// scanner may hold a file just copied), the link is put back: the account
// keeps seeing the group's content, and the next run copies again.
func TestCopyBackPutsTheLinkBack(t *testing.T) {
	f := newFixture(t)
	ext := config.Share{Extensions: "default"}
	a := f.profile(t, "a", map[string]string{"skills/s/SKILL.md": "s"})
	share(t, f, a, config.Share{}, ext, Options{})
	orig := rename
	rename = func(string, string) error { return os.ErrPermission }
	pl, err := Build(a, f.groups, ext, config.Share{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	err = pl.Apply()
	rename = orig
	if err == nil {
		t.Fatal("no error")
	}
	if target, ok := fsutil.LinkTarget(filepath.Join(a, "skills")); !ok || target != filepath.Join(f.groups, "default", "skills") {
		t.Fatalf("skills is not the link again: %q", target)
	}
	if fsutil.Exists(filepath.Join(a, BackupDir)) {
		t.Fatal("the copy was left behind")
	}
	share(t, f, a, config.Share{}, config.Share{}, Options{})
	if read(t, filepath.Join(a, "skills", "s", "SKILL.md")) != "s" {
		t.Fatal("not copied the second time")
	}
}

// A plan that fills an empty group slot with an account's own link target
// (dotfiles) leaves a file the group got there since the plan was made
// (another account's session writing through its link).
func TestCopyIntoAGroupLeavesWhatAppeared(t *testing.T) {
	f := newFixture(t)
	dotfiles := filepath.Join(f.root, "dotfiles", "CLAUDE.md")
	os.MkdirAll(filepath.Dir(dotfiles), 0o700)
	os.WriteFile(dotfiles, []byte("from dotfiles"), 0o600)
	a := f.profile(t, "a", nil)
	os.Symlink(dotfiles, filepath.Join(a, "CLAUDE.md"))
	pl, err := Build(a, f.groups, config.Share{}, config.Share{Memory: "default"}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	group := filepath.Join(f.groups, "default", "CLAUDE.md")
	os.MkdirAll(filepath.Dir(group), 0o700)
	os.WriteFile(group, []byte("written meanwhile"), 0o600)
	if err := pl.Apply(); err == nil {
		t.Error("no error")
	}
	if got := read(t, group); got != "written meanwhile" {
		t.Fatalf("the group's file reads %q", got)
	}
}

// A group folder can hold a link the user made into the joining account
// (a skill shared by hand before ccam). That link is the account's own
// file, not a copy of it: merging must keep the content, now in the group.
func TestMergeKeepsAFileTheGroupLinksBackTo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("needs a symlink")
	}
	f := newFixture(t)
	b := f.profile(t, "b", map[string]string{"skills/foo/SKILL.md": "the only copy", "skills/bar.md": "bar"})
	grp := filepath.Join(f.groups, "g", "skills")
	if err := os.MkdirAll(grp, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(b, "skills", "foo"), filepath.Join(grp, "foo")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(b, "skills", "bar.md"), filepath.Join(grp, "bar.md")); err != nil {
		t.Fatal(err)
	}
	pl, err := Build(b, f.groups, config.Share{}, config.Share{Extensions: "g"}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := pl.Apply(); err != nil {
		t.Fatal(err)
	}
	for rel, want := range map[string]string{"skills/foo/SKILL.md": "the only copy", "skills/bar.md": "bar"} {
		data, err := os.ReadFile(filepath.Join(b, rel))
		if err != nil || string(data) != want {
			t.Errorf("%s through the account = %q, %v; want %q", rel, data, err, want)
		}
	}
	if _, ok := fsutil.LinkTarget(filepath.Join(grp, "foo")); ok {
		t.Error("the group still holds the link instead of the content")
	}
}

// A copy back that stopped after the link was removed, before the copy was
// renamed into place, leaves the complete copy in the stage: doctor finds
// it, and both its repair and the plan of a resumed change put it back.
func TestCopyBackStoppedAfterTheUnlinkIsFinished(t *testing.T) {
	for _, viaBuild := range []bool{false, true} {
		f := newFixture(t)
		a := f.profile(t, "a", map[string]string{"skills/s/SKILL.md": "mine"})
		share(t, f, a, config.Share{}, config.Share{Extensions: "g"}, Options{})
		stage := filepath.Join(a, BackupDir, "incomplete-skills")
		if err := fsutil.CopyTree(filepath.Join(f.groups, "g", "skills"), stage); err != nil {
			t.Fatal(err)
		}
		os.Remove(filepath.Join(a, "skills")) // the link; stopped here
		var probs []Problem
		for _, pr := range CheckTargets(a, config.Share{}, f.groups) {
			if pr.Entry.Name == "skills" {
				probs = append(probs, pr)
			}
		}
		if len(probs) != 1 || !strings.Contains(probs[0].String(), stage) || !probs[0].Safe() {
			t.Fatalf("problems: %v", probs)
		}
		if viaBuild {
			share(t, f, a, config.Share{}, config.Share{}, Options{})
		} else if rest, err := QuickFix(probs); err != nil || len(rest) != 0 {
			t.Fatalf("quick fix: %v %v", rest, err)
		}
		if got := read(t, filepath.Join(a, "skills", "s", "SKILL.md")); got != "mine" {
			t.Fatalf("skills/s/SKILL.md = %q", got)
		}
		if fsutil.Exists(stage) {
			t.Fatal("the stage is still there")
		}
	}
}

// A launch between the stop and doctor makes the entry anew: the stage is
// still reported, as a warning naming it, and nothing moves it over what
// Claude Code wrote since. A plan that puts a stage in place says so.
func TestStageFoundAfterALaunch(t *testing.T) {
	f := newFixture(t)
	a := f.profile(t, "a", map[string]string{"projects/p/s.jsonl": "history"})
	share(t, f, a, config.Share{}, config.Share{History: "g"}, Options{})
	stage := filepath.Join(a, BackupDir, "incomplete-projects")
	if err := fsutil.CopyTree(filepath.Join(f.groups, "g", "projects"), stage); err != nil {
		t.Fatal(err)
	}
	os.Remove(filepath.Join(a, "projects")) // the link; stopped here
	pl, err := Build(a, f.groups, config.Share{}, config.Share{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	pl.Describe(&buf)
	if strings.Contains(buf.String(), "first content of the group") || !strings.Contains(buf.String(), "incomplete-projects") {
		t.Errorf("plan:\n%s", buf.String())
	}
	os.MkdirAll(filepath.Join(a, "projects", "-new"), 0o700) // the next launch
	var probs []Problem
	for _, pr := range CheckTargets(a, config.Share{}, f.groups) {
		if pr.Entry.Name == "projects" {
			probs = append(probs, pr)
		}
	}
	if len(probs) != 1 || probs[0].Kind != StagedAside || probs[0].Safe() || !strings.Contains(probs[0].String(), stage) {
		t.Fatalf("problems: %v", probs)
	}
	pl, err = Build(a, f.groups, config.Share{}, config.Share{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range pl.Steps {
		if st.Entry.Name == "projects" {
			t.Fatalf("plan step for projects: %+v", st)
		}
	}
	if !fsutil.Exists(stage) || !fsutil.Exists(filepath.Join(a, "projects", "-new")) {
		t.Fatal("something moved")
	}
}

// A group entry that is itself a link to the joining account's entry (one
// a user made in ccam's groups folder) holds nothing of its own: the
// account's content becomes the group's, and the account links to it.
func TestGroupEntryLinkingToTheAccount(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("needs a symlink")
	}
	for _, tc := range []struct {
		entry string
		files map[string]string
		s     config.Share
		read  string
	}{
		{"CLAUDE.md", map[string]string{"CLAUDE.md": "the only copy"}, config.Share{Memory: "g"}, "CLAUDE.md"},
		{"skills", map[string]string{"skills/foo/SKILL.md": "the only copy"}, config.Share{Extensions: "g"}, "skills/foo/SKILL.md"},
	} {
		f := newFixture(t)
		b := f.profile(t, "b", tc.files)
		os.MkdirAll(filepath.Join(f.groups, "g"), 0o700)
		if err := os.Symlink(filepath.Join(b, tc.entry), filepath.Join(f.groups, "g", tc.entry)); err != nil {
			t.Fatal(err)
		}
		share(t, f, b, config.Share{}, tc.s, Options{})
		if got := read(t, filepath.Join(b, filepath.FromSlash(tc.read))); got != "the only copy" {
			t.Fatalf("%s: %q", tc.entry, got)
		}
		if _, link := fsutil.LinkTarget(filepath.Join(f.groups, "g", tc.entry)); link {
			t.Fatalf("%s: the group entry is still a link", tc.entry)
		}
		if _, link := fsutil.LinkTarget(filepath.Join(b, tc.entry)); !link {
			t.Fatalf("%s: the account does not link to the group", tc.entry)
		}
	}
}
