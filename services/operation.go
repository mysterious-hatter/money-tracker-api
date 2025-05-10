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
	ErrUnsupportedFilter       error = errors.New("unsupported filter")
)

type OperationService struct {
	storage storage.Storage
}

func NewOperationService(st storage.Storage) *OperationService {
	ops := OperationService{st}
	return &ops
}

func (ops *OperationService) CreateOperation(operation *models.Operation, userId int64) (int64, error) {
	// Check user's ownership of the wallet
	err := ops.checkOwnershipByConnectedWallet(operation.WalletId, userId)
	if err != nil {
		return 0, err
	}

	id, err := ops.storage.CreateOperation(operation)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// walletId and userId are obligatory, sinceDate and sortBy are optional
func (ops *OperationService) GetOperations(userId, walletId int64, sinceDate models.DateOnly, sortBy string) ([]models.Operation, error) {
	if sortBy != "date" || sortBy != "sum" {
		return nil, ErrUnsupportedFilter
	}
	// Check user's ownership of the wallet
	err := ops.checkOwnershipByConnectedWallet(walletId, userId)
	if err != nil {
		return nil, err
	}

	operations, err := ops.storage.GetOperations(walletId, sinceDate, sortBy)
	return operations, err
}

func (ops *OperationService) GetOperationById(operationId int64, userId int64) (*models.Operation, error) {
	operation, err := ops.storage.GetOperationById(operationId)
	if err != nil {
		return nil, ErrOperationNotFound
	}

	// Check user's ownership of the operation
	err = ops.checkOwnershipByConnectedWallet(operation.WalletId, userId)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (ops *OperationService) UpdateOperation(operation *models.Operation, userId int64) error {
	// Check user's ownership of the operation and if it exists
	_, err := ops.GetOperationById(operation.Id, userId)
	if err != nil {
		return err
	}

	err = ops.storage.UpdateOperation(operation)
	if err != nil {
		return err
	}

	return nil
}

func (ops *OperationService) DeleteOperation(operationId int64, userId int64) error {
	// Check user's ownership of the operation and if it exists
	_, err := ops.GetOperationById(operationId, userId)
	if err != nil {
		return err
	}

	err = ops.storage.DeleteOperation(operationId)
	if err != nil {
		return err
	}

	return nil
}

func (ops *OperationService) checkOwnershipByConnectedWallet(walletId, userId int64) error {
	connectedWallet, err := ops.storage.GetWalletById(walletId)
	if err != nil {
		return ErrConnectedWalletNotFound
	}
	// Check user's ownership of the operation
	if err := checkOwnership(connectedWallet.OwnerId, userId); err != nil {
		return err
	}

	return nil
}
