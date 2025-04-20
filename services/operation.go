package services

import (
	"errors"
	"time"

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
	ops := OperationService{st}
	return &ops
}

func (ops *OperationService) CreateOperation(operation *models.Operation, userID int64) (int64, error) {
	// Check user's ownership of the wallet
	err := ops.checkOwnershipByConnectedWallet(operation.WalletID, userID)
	if err != nil {
		return 0, err
	}

	id, err := ops.storage.CreateOperation(operation)
	if err != nil {
		return 0, ErrUnableToCreate
	}
	return id, nil
}

func (ops *OperationService) GetOperationsByWalletID(walletID, userID int64) ([]models.Operation, error) {
	// Check user's ownership of the wallet
	err := ops.checkOwnershipByConnectedWallet(walletID, userID)
	if err != nil {
		return nil, err
	}

	operations, err := ops.storage.GetOperationsByWalletID(walletID)
	if err != nil {
		return nil, ErrNoOperationsFound
	}
	return operations, nil
}

func (ops *OperationService) GetOperationsSinceDateByWalletID(walletID, userID int64, date time.Time) ([]models.Operation, error) {
	// Check user's ownership of the wallet
	err := ops.checkOwnershipByConnectedWallet(walletID, userID)
	if err != nil {
		return nil, err
	}

	operations, err := ops.storage.GetOperationsSinceDateByWalletID(walletID, date)
	if err != nil {
		return nil, ErrNoOperationsFound
	}
	return operations, nil
}

func (ops *OperationService) GetOperationByID(operationID int64, userID int64) (*models.Operation, error) {
	operation, err := ops.storage.GetOperationByID(operationID)
	if err != nil {
		return nil, ErrOperationNotFound
	}

	// Check user's ownership of the operation
	err = ops.checkOwnershipByConnectedWallet(operation.WalletID, userID)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (ops *OperationService) UpdateOperation(operation *models.Operation, userID int64) error {
	// Check user's ownership of the operation and if it exists
	_, err := ops.GetOperationByID(operation.ID, userID)
	if err != nil {
		return err
	}

	err = ops.storage.UpdateOperation(operation)
	if err != nil {
		return ErrUnableToUpdate
	}

	return nil
}

func (ops *OperationService) DeleteOperation(operationID int64, userID int64) error {
	// Check user's ownership of the operation and if it exists
	_, err := ops.GetOperationByID(operationID, userID)
	if err != nil {
		return err
	}

	err = ops.storage.DeleteOperation(operationID)
	if err != nil {
		return ErrUnableToDelete
	}

	return nil
}

func (ops *OperationService) checkOwnershipByConnectedWallet(walletID, userID int64) error {
	connectedWallet, err := ops.storage.GetWalletByID(walletID)
	if err != nil {
		return ErrConnectedWalletNotFound
	}
	// Check user's ownership of the operation
	if err := checkOwnership(connectedWallet.OwnerID, userID); err != nil {
		return err
	}

	return nil
}
