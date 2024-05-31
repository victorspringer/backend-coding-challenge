package domain

// ValidatedRating is used to validate an instance of Rating data.
type ValidatedRating struct {
	Rating
	isValidated bool
}

// IsValid returns true if the instance of Rating is validated.
func (vr *ValidatedRating) IsValid() bool {
	return vr.isValidated
}

// NewValidatedRating returns an instance of ValidatedRating if the given Rating instance is valid.
func NewValidatedRating(rating *Rating) (*ValidatedRating, error) {
	if err := rating.validate(); err != nil {
		return nil, err
	}

	return &ValidatedRating{
		Rating:      *rating,
		isValidated: true,
	}, nil
}
