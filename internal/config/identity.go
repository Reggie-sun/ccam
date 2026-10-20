package config

// Identity is who an account's login belongs to, as Claude Code records it
// in .claude.json (oauthAccount). ccam keeps a copy from the last login, to
// notice a login into another account and to keep automatic choices within
// one organization.
type Identity struct {
	AccountUUID string `toml:"account_uuid,omitempty"`
	Email       string `toml:"email,omitempty"`
	OrgUUID     string `toml:"org_uuid,omitempty"`
	OrgName     string `toml:"org_name,omitempty"`
	OrgType     string `toml:"org_type,omitempty"` // claude_pro, claude_max, claude_team, claude_enterprise
}

// Personal is the pool of accounts with a personal plan (Pro, Max).
const Personal = "personal"

// Pool is the set of accounts automatic choices stay in: one per team or
// enterprise organization, and one for every personal plan. "" means
// unknown, which fits any pool.
func (i *Identity) Pool() string {
	if i == nil || i.OrgType == "" {
		return ""
	}
	switch i.OrgType {
	case "claude_team", "claude_enterprise":
		return "org:" + i.OrgUUID
	}
	return Personal
}

// SamePool reports whether a and b may stand in for each other.
func SamePool(a, b *Identity) bool {
	pa, pb := a.Pool(), b.Pool()
	return pa == "" || pb == "" || pa == pb
}

// SameAccount reports whether a and b are one Claude account in one
// organization, which share their usage limits.
func SameAccount(a, b *Identity) bool {
	return a != nil && b != nil && a.AccountUUID != "" && a.AccountUUID == b.AccountUUID && a.OrgUUID == b.OrgUUID
}

// PlanName is the plan as Claude Code names it: Claude Max, Claude Team …
func (i *Identity) PlanName() string {
	if i == nil {
		return ""
	}
	switch i.OrgType {
	case "claude_pro":
		return "Claude Pro"
	case "claude_max":
		return "Claude Max"
	case "claude_team":
		return "Claude Team"
	case "claude_enterprise":
		return "Claude Enterprise"
	}
	return ""
}

// OrgLabel names the pool for messages: the organization, or "personal".
func (i *Identity) OrgLabel() string {
	switch p := i.Pool(); {
	case p == Personal:
		return "personal plans"
	case p != "" && i.OrgName != "":
		return i.OrgName
	case p != "":
		return "its organization"
	}
	return ""
}
