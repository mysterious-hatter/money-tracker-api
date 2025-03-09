package services

import (
	"errors"
)

var (
	ErrSomethingWentWrong error = errors.New("something went wrong")
	ErrUnableToCreate     error = errors.New("unable to create")
	ErrUnableToUpdate     error = errors.New("unable to update")
	ErrUnableToDelete     error = errors.New("unable to delete")
	ErrAccessDenied       error = errors.New("access denied")
)

func checkOwnership(ownerID, userID int64) error {
	if ownerID != userID {
		return ErrAccessDenied
	}
	return nil
}
