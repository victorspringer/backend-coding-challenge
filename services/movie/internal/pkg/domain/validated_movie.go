package domain

// ValidatedMovie is used to validate an instance of Movie data.
type ValidatedMovie struct {
	Movie
	isValidated bool
}

// IsValid returns true if the instance of Movie is validated.
func (vm *ValidatedMovie) IsValid() bool {
	return vm.isValidated
}

// NewValidatedMovie returns an instance of ValidatedMovie if the given Movie instance is valid.
func NewValidatedMovie(movie *Movie) (*ValidatedMovie, error) {
	if err := movie.validate(); err != nil {
		return nil, err
	}

	return &ValidatedMovie{
		Movie:       *movie,
		isValidated: true,
	}, nil
}
