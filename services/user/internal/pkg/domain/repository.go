package domain

// Repository is the interface for the domain's repository (e.g. some database).
type Repository interface {
	// Create receives a validated input and creates a new User.
	Create(user *ValidatedUser) (*User, error)
	// FindById retrieves an User by a given unique ID.
	FindByID(id string) (*User, error)
	// FindByUsername retrieves an User by a given unique username.
	FindByUsername(username string) (*User, error)
}
