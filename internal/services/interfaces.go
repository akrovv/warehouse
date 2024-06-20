package services

import "github.com/akrovv/warehouse/internal/domain"

type ProductStorage interface {
	Create(product *domain.Product) error
	Reserve(wp *domain.WarehouseProduct) error
	CancelReservation(wp *domain.WarehouseProduct) error
	Transfer(td *domain.TransferProduct) error
	Add(ad *domain.AddProduct) error
	Delete(dp *domain.DeleteProduct) (*domain.Product, error)
}

type WarehouseStorage interface {
	Create(warehouse *domain.Warehouse) error
	GetLeftOvers(gw *domain.GetFromWarehouse) ([]domain.Product, error)
}
