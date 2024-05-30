package domain

// ValidatedUser is used to validate an instance of User data.
type ValidatedUser struct {
	User
	isValidated bool
}

// IsValid returns true if the instance of User is validated.
func (vp *ValidatedUser) IsValid() bool {
	return vp.isValidated
}

// NewValidatedUser returns an instance of ValidatedUser if the given User instance is valid.
func NewValidatedUser(user *User) (*ValidatedUser, error) {
	if err := user.validate(); err != nil {
		return nil, err
	}

	return &ValidatedUser{
		User:        *user,
		isValidated: true,
	}, nil
}
