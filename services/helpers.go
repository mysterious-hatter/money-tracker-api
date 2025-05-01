package services

import (
	"errors"
)

var (
	ErrSomethingWentWrong error = errors.New("something went wrong")
	ErrAccessDenied       error = errors.New("access denied")
	// Following errors may be reduced, as they don't give any useful information
	ErrUnableToCreate     error = errors.New("unable to create")
	ErrUnableToUpdate     error = errors.New("unable to update")
	ErrUnableToDelete     error = errors.New("unable to delete")
)

func checkOwnership(ownerId, userId int64) error {
	if ownerId != userId {
		return ErrAccessDenied
	}
	return nil
}
