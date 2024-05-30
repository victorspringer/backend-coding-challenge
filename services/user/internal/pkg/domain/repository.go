package domain

// TODO: add comments

type Repository interface {
	Create(user *ValidatedUser) (*User, error)
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
}
