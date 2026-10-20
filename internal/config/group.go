package config

import (
	"fmt"
	"slices"
	"strings"
)

// Group shares Items between its Members: each member's folder links those
// items to groups/<Folder()>/. An account can be in several groups as long
// as no item comes from two of them.
type Group struct {
	Name string `toml:"name"`
	// Dir is the folder under groups/ when it is not Name: after a rename
	// (the folder keeps its name) or for groups split off by the conversion
	// of the old per-account format.
	Dir     string   `toml:"dir,omitempty"`
	Items   []string `toml:"items"`
	Members []string `toml:"members"`
}

// Folder is the group's folder name under groups/.
func (g *Group) Folder() string {
	if g.Dir != "" {
		return g.Dir
	}
	return g.Name
}

func (g *Group) Has(item string) bool { return slices.Contains(g.Items, item) }

func (g *Group) HasMember(name string) bool { return slices.Contains(g.Members, name) }

// reservedGroups are values that mean something else where a group is named.
var reservedGroups = map[string]bool{Local: true, "on": true, "off": true, "none": true}

// ValidateGroupName checks a name for a new group.
func ValidateGroupName(name string) error {
	if reservedGroups[name] {
		return fmt.Errorf("%q cannot be a group name", name)
	}
	return ValidateGroup(name)
}

// ParseItems reads a comma-separated item list such as "mcp,settings" and
// returns the items in display order.
func ParseItems(s string) ([]string, error) {
	var out []string
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !slices.Contains(Items, f) {
			return nil, fmt.Errorf("unknown item %q (items: %s)", f, strings.Join(Items, ", "))
		}
		out = append(out, f)
	}
	return SortItems(out), nil
}

// SortItems puts items in display order without duplicates.
func SortItems(items []string) []string {
	var out []string
	for _, item := range Items {
		if slices.Contains(items, item) {
			out = append(out, item)
		}
	}
	return out
}

// DefaultItems is what group default shares when ccam creates it.
var DefaultItems = []string{"mcp", "settings", "memory", "extensions"}

// DefaultGroups is what a new account joins unless told otherwise: group
// default, or the group that uses its folder after a rename. Before there
// are groups, the first account creates default.
func (c *Config) DefaultGroups() []string {
	if _, ok := c.FindGroup("default"); ok || len(c.Groups) == 0 {
		return []string{"default"}
	}
	// Of the groups using it, the largest is the one the conversion (or a
	// rename) left with the folder.
	if users := c.FolderUsers("default"); len(users) > 0 {
		best := users[0]
		for _, g := range users[1:] {
			if len(g.Items)*len(g.Members) > len(best.Items)*len(best.Members) {
				best = g
			}
		}
		return []string{best.Name}
	}
	return nil
}

// IsAccount reports whether an account is called exactly name (not a slot
// number, as Find also accepts).
func (c *Config) IsAccount(name string) bool {
	return slices.ContainsFunc(c.Profiles, func(p Profile) bool { return p.Name == name })
}

// DropMember takes name out of every group. An account that takes a name
// calls it first: a member of that name was not an account until now, and
// would otherwise hand the new account another's groups.
func (c *Config) DropMember(name string) {
	for i := range c.Groups {
		c.Groups[i].Members = slices.DeleteFunc(c.Groups[i].Members, func(m string) bool { return m == name })
	}
}

// FindGroup returns the group called name.
func (c *Config) FindGroup(name string) (*Group, bool) {
	for i := range c.Groups {
		if c.Groups[i].Name == name {
			return &c.Groups[i], true
		}
	}
	return nil, false
}

// MustFindGroup is FindGroup with a user-facing error.
func (c *Config) MustFindGroup(name string) (*Group, error) {
	if g, ok := c.FindGroup(name); ok {
		return g, nil
	}
	return nil, fmt.Errorf("no group %q (see `ccam group`)", name)
}

// GroupsOf lists the groups account is a member of.
func (c *Config) GroupsOf(account string) []*Group {
	var out []*Group
	for i := range c.Groups {
		if c.Groups[i].HasMember(account) {
			out = append(out, &c.Groups[i])
		}
	}
	return out
}

// GroupFor returns the group account gets item from, or nil.
func (c *Config) GroupFor(account, item string) *Group {
	for _, g := range c.GroupsOf(account) {
		if g.Has(item) {
			return g
		}
	}
	return nil
}

// resolved maps each account to what its groups share; it fails on a group
// that cannot be used or on an account that would get an item twice.
func (c *Config) resolved() (map[string]Share, error) {
	out := make(map[string]Share, len(c.Profiles))
	for _, p := range c.Profiles {
		out[p.Name] = Share{}
	}
	from := map[[2]string]string{}    // account, item → group name
	folders := map[[2]string]string{} // folder, item → group name
	names := map[string]bool{}
	for i := range c.Groups {
		g := &c.Groups[i]
		if err := ValidateGroupName(g.Name); err != nil {
			return nil, err
		}
		if names[g.Name] {
			return nil, fmt.Errorf("group %s is defined twice", g.Name)
		}
		names[g.Name] = true
		// "local" is what Share stores for an item that is not shared.
		if g.Dir != "" && (!isName(g.Dir) || g.Dir == Local) {
			return nil, fmt.Errorf("group %s has an invalid folder %q", g.Name, g.Dir)
		}
		for _, item := range g.Items {
			if !slices.Contains(Items, item) {
				return nil, fmt.Errorf("group %s has an unknown item %q (items: %s)", g.Name, item, strings.Join(Items, ", "))
			}
			// Two groups sharing one item through one folder would quietly be
			// one group.
			k := [2]string{g.Folder(), item}
			if prev, dup := folders[k]; dup && prev != g.Name {
				return nil, fmt.Errorf("groups %s and %s both share %s through the folder groups/%s; give one of them its own folder", prev, g.Name, item, g.Folder())
			}
			folders[k] = g.Name
		}
		for _, m := range g.Members {
			s, ok := out[m]
			if !ok {
				continue // not an account (any more); doctor reports it
			}
			for _, item := range g.Items {
				k := [2]string{m, item}
				if prev, dup := from[k]; dup && prev != g.Name {
					return nil, fmt.Errorf("account %s gets %s from both group %s and group %s; an item can come from one group only", m, item, prev, g.Name)
				}
				from[k] = g.Name
				*s.field(item) = g.Folder()
			}
			out[m] = s
		}
	}
	return out, nil
}

// Resolve sets every account's Share from the groups.
func (c *Config) Resolve() error {
	m, err := c.resolved()
	if err != nil {
		return err
	}
	for i := range c.Profiles {
		c.Profiles[i].Share = m[c.Profiles[i].Name]
	}
	return nil
}

// ConvertShares replaces the groups with ones built from every account's
// Share; Load uses it to read the old per-account format. Within a folder,
// items shared by the same accounts become one group. The largest keeps the
// folder's name; the rest are named <folder>-<item> or <folder>-<n> and
// keep using the folder (Dir), so nothing moves on disk.
func (c *Config) ConvertShares() {
	var dirs []string
	taken := map[string]bool{}
	for _, p := range c.Profiles {
		for _, item := range Items {
			if d := p.Share.Get(item); d != Local && !taken[d] {
				taken[d] = true
				dirs = append(dirs, d)
			}
		}
	}
	var groups []Group
	for _, dir := range dirs {
		var parts []Group
		index := map[string]int{}
		for _, item := range Items {
			var members []string
			for _, p := range c.Profiles {
				if p.Share.Get(item) == dir {
					members = append(members, p.Name)
				}
			}
			if len(members) == 0 {
				continue
			}
			k := strings.Join(members, ",")
			if i, ok := index[k]; ok {
				parts[i].Items = append(parts[i].Items, item)
				continue
			}
			index[k] = len(parts)
			parts = append(parts, Group{Items: []string{item}, Members: members})
		}
		main := 0
		for i, g := range parts {
			if len(g.Members)*len(g.Items) > len(parts[main].Members)*len(parts[main].Items) {
				main = i
			}
		}
		n := 2
		for i := range parts {
			g := &parts[i]
			if i == main && !reservedGroups[dir] {
				g.Name = dir
				continue
			}
			g.Dir = dir
			base := dir + "-" + g.Items[0]
			if len(g.Items) > 1 {
				base = fmt.Sprintf("%s-%d", dir, n)
				n++
			}
			g.Name = uniqueName(base, taken)
		}
		groups = append(groups, parts...)
	}
	c.Groups = groups
}

// uniqueName shortens base to a valid length and adds -2, -3, … until no
// name in taken has it, then records it.
func uniqueName(base string, taken map[string]bool) string {
	name := base[:min(len(base), 32)]
	for i := 2; taken[name]; i++ {
		sfx := fmt.Sprintf("-%d", i)
		name = base[:min(len(base), 32-len(sfx))] + sfx
	}
	taken[name] = true
	return name
}

// AddGroup creates a group; members must be accounts.
func (c *Config) AddGroup(name string, items, members []string) error {
	if err := ValidateGroupName(name); err != nil {
		return err
	}
	if _, ok := c.FindGroup(name); ok {
		return fmt.Errorf("a group named %q already exists", name)
	}
	for _, g := range c.Groups {
		if g.Folder() == name {
			return fmt.Errorf("group %s uses the folder groups/%s; choose another name", g.Name, name)
		}
	}
	for _, m := range members {
		if _, err := c.MustFind(m); err != nil {
			return err
		}
	}
	c.Groups = append(c.Groups, Group{Name: name, Items: SortItems(items), Members: memberNames(c, members)})
	return c.Resolve()
}

// memberNames turns account names or numbers into names, without duplicates.
func memberNames(c *Config, sels []string) []string {
	var out []string
	for _, s := range sels {
		if p, ok := c.Find(s); ok && !slices.Contains(out, p.Name) {
			out = append(out, p.Name)
		}
	}
	return out
}

// SetGroupItems changes what a group shares.
func (c *Config) SetGroupItems(name string, items []string) error {
	g, err := c.MustFindGroup(name)
	if err != nil {
		return err
	}
	g.Items = SortItems(items)
	return c.Resolve()
}

// JoinGroup adds account to a group. An item the account already gets from
// another group is refused, unless replace: then the account leaves those
// groups first. It returns the groups left.
func (c *Config) JoinGroup(account, name string, replace bool) ([]string, error) {
	p, err := c.MustFind(account)
	if err != nil {
		return nil, err
	}
	g, err := c.MustFindGroup(name)
	if err != nil {
		return nil, err
	}
	if g.HasMember(p.Name) {
		return nil, nil
	}
	var left []string
	for _, other := range c.GroupsOf(p.Name) {
		var both []string
		for _, item := range g.Items {
			if other.Has(item) {
				both = append(both, item)
			}
		}
		if len(both) == 0 {
			continue
		}
		if !replace {
			return nil, fmt.Errorf("%s already gets %s from group %s, and an item can come from one group only.\n"+
				"Use --replace to leave %s (it shares %s), or take %s out of it with `ccam group set %s --remove %s`",
				p.Name, strings.Join(both, ", "), other.Name, other.Name, strings.Join(other.Items, ", "), both[0], other.Name, both[0])
		}
		left = append(left, other.Name)
		other.Members = slices.DeleteFunc(other.Members, func(m string) bool { return m == p.Name })
	}
	g.Members = append(g.Members, p.Name)
	return left, c.Resolve()
}

// LeaveGroup takes account out of a group.
func (c *Config) LeaveGroup(account, name string) error {
	p, err := c.MustFind(account)
	if err != nil {
		return err
	}
	g, err := c.MustFindGroup(name)
	if err != nil {
		return err
	}
	if !g.HasMember(p.Name) {
		return fmt.Errorf("%s is not in group %s", p.Name, g.Name)
	}
	g.Members = slices.DeleteFunc(g.Members, func(m string) bool { return m == p.Name })
	return c.Resolve()
}

// RenameGroup renames a group in the config only: its folder keeps its name.
func (c *Config) RenameGroup(oldName, newName string) error {
	g, err := c.MustFindGroup(oldName)
	if err != nil {
		return err
	}
	if err := ValidateGroupName(newName); err != nil {
		return err
	}
	if _, ok := c.FindGroup(newName); ok {
		return fmt.Errorf("a group named %q already exists", newName)
	}
	g.Dir = g.Folder()
	g.Name = newName
	if g.Dir == g.Name {
		g.Dir = ""
	}
	return c.Resolve()
}

// RemoveGroup drops a group from the config.
func (c *Config) RemoveGroup(name string) error {
	if _, err := c.MustFindGroup(name); err != nil {
		return err
	}
	c.Groups = slices.DeleteFunc(c.Groups, func(g Group) bool { return g.Name == name })
	return c.Resolve()
}

// FolderUsers lists the groups that use folder dir.
func (c *Config) FolderUsers(dir string) []*Group {
	var out []*Group
	for i := range c.Groups {
		if c.Groups[i].Folder() == dir {
			out = append(out, &c.Groups[i])
		}
	}
	return out
}

// Clone copies the config deeply enough that edits to the copy's groups and
// accounts leave c as it is.
func (c *Config) Clone() *Config {
	n := *c
	n.Profiles = slices.Clone(c.Profiles)
	n.Folders = slices.Clone(c.Folders)
	n.Groups = make([]Group, len(c.Groups))
	for i, g := range c.Groups {
		g.Items = slices.Clone(g.Items)
		g.Members = slices.Clone(g.Members)
		n.Groups[i] = g
	}
	for i := range n.Profiles {
		if a := n.Profiles[i].API; a != nil {
			cp := *a
			n.Profiles[i].API = &cp
		}
		if l := n.Profiles[i].Login; l != nil {
			cp := *l
			n.Profiles[i].Login = &cp
		}
	}
	return &n
}
