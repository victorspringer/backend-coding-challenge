package domain

// TODO: add comments

type ValidatedUser struct {
	User
	isValidated bool
}

func (vp *ValidatedUser) IsValid() bool {
	return vp.isValidated
}

func NewValidatedUser(user *User) (*ValidatedUser, error) {
	if err := user.validate(); err != nil {
		return nil, err
	}

	return &ValidatedUser{
		User:        *user,
		isValidated: true,
	}, nil
}
