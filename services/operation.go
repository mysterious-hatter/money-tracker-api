package services

import (
	"errors"

	"finances-backend/models"
	"finances-backend/storage"
)

var (
	ErrOperationNotFound       error = errors.New("operation not found")
	ErrConnectedWalletNotFound error = errors.New("connected wallet not found")
	ErrNoOperationsFound       error = errors.New("no operations found")
)

type OperationService struct {
	storage storage.Storage
}

func NewOperationService(st storage.Storage) *OperationService {
	os := OperationService{st}
	return &os
}

func (os *OperationService) CreateOperation(operation *models.Operation) (int64, error) {
	id, err := os.storage.CreateOperation(operation)
	if err != nil {
		return 0, ErrUnableToCreate
	}
	return id, nil
}

func (os *OperationService) GetAllOperations(walletID int64) ([]models.Operation, error) {
	operations, err := os.storage.GetOperationsByWalletID(walletID)
	if err != nil {
		return nil, ErrNoOperationsFound
	}
	return operations, nil
}

func (os *OperationService) GetOperationByID(operationID int64, userID int64) (*models.Operation, error) {
	operation, err := os.storage.GetOperationByID(operationID)
	if err != nil {
		return nil, ErrOperationNotFound
	}

	// Check user's ownership of the operation
	err = os.checkOwnershipByConnectedWallet(operation.WalletID, userID)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (os *OperationService) UpdateOperation(operation *models.Operation, userID int64) error {
	// Check user's ownership of the operation
	err := os.checkOwnershipByConnectedWallet(operation.WalletID, userID)
	if err != nil {
		return err
	}

	err = os.storage.UpdateOperation(operation)
	if err != nil {
		return ErrUnableToUpdate
	}

	return nil
}

func (os *OperationService) DeleteOperation(operationID int64, userID int64) error {
	operation, err := os.storage.GetOperationByID(operationID)
	if err != nil {
		return ErrOperationNotFound
	}

	// Check user's ownership of the operation
	err = os.checkOwnershipByConnectedWallet(operation.WalletID, userID)
	if err != nil {
		return err
	}

	err = os.storage.DeleteOperation(operationID)
	if err != nil {
		return ErrUnableToDelete
	}

	return nil
}

func (os *OperationService) checkOwnershipByConnectedWallet(walletID, userID int64) error {
	connectedWallet, err := os.storage.GetWalletByID(walletID)
	if err != nil {
		return ErrConnectedWalletNotFound
	}
	// Check user's ownership of the operation
	if err := checkOwnership(connectedWallet.OwnerID, userID); err != nil {
		return err
	}

	return nil
}