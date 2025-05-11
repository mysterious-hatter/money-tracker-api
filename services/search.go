package services

import (
	"finances-backend/models"
	"finances-backend/storage"
)

type SearchService struct {
	storage storage.Storage
}

func NewSearchService(st storage.Storage) *SearchService {
	ss := SearchService{st}
	return &ss
}

func (ss *SearchService) SearchOperations(userId int64, name string, walletId int64, date models.DateOnly,
	place string, categoryId int64, sortBy string) ([]models.Operation, error) {
	if sortBy != "date" && sortBy != "sum" && sortBy != "" {
		return nil, ErrUnsupportedFilter
	}
	// Check user's ownership of the wallet
	err := checkOwnershipByConnectedWallet(walletId, userId, ss.storage)
	if err != nil {
		return nil, err
	}

	operations, err := ss.storage.SearchOperations(name, walletId, date, place, categoryId, sortBy)
	if err != nil {
		return nil, err
	}

	return operations, nil
}
