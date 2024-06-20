package services

import "github.com/akrovv/warehouse/internal/domain"

type warehouseService struct {
	storage WarehouseStorage
}

func NewWarehouseService(storage WarehouseStorage) *warehouseService {
	return &warehouseService{
		storage: storage,
	}
}

func (s *warehouseService) Create(warehouse *domain.Warehouse) error {
	return s.storage.Create(warehouse)
}

func (s *warehouseService) GetLeftOvers(gw *domain.GetFromWarehouse) ([]domain.Product, error) {
	return s.storage.GetLeftOvers(gw)
}
