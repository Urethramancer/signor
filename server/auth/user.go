package auth

// User structure for the main system.
type User struct {
	// ID is a unique identifier in the database.
	ID int64 `json:"id"`
	// Name is the username for logins and the default display name.
	Name string `json:"name"`
	// Email is used for verification and password retrieval.
	Email string `json:"email"`
	// Fullnames is the real name of the user in western order (first name
	// followed by middle names, ending in last names and everything else between).
	Fullnames []string `json:"fullnames,omitempty"`
	// Title if applicable.
	Title string `json:"title,omitempty"`
	// Phone number is optional.
	Phone string `json:"phone,omitempty"`
	// Address is optional.
	Address []string `json:"userid,omitempty"`
}
