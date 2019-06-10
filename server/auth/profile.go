package auth

// Profile is used per domain.
type Profile struct {
	// UserID in the user table.
	UserID int64 `json:"userid"`
	// Domain (site) this profile is for.
	Domain string `json:"domain"`
	// Username displayed on the site.
	Username string `json:"username"`
}
