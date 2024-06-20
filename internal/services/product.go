package services

import "github.com/akrovv/warehouse/internal/domain"

type productService struct {
	storage ProductStorage
}

func NewProductService(storage ProductStorage) *productService {
	return &productService{
		storage: storage,
	}
}

func (s *productService) Create(product *domain.Product) error {
	return s.storage.Create(product)
}

func (s *productService) Reserve(wp *domain.WarehouseProduct) error {
	return s.storage.Reserve(wp)
}

func (s *productService) CancelReservation(wp *domain.WarehouseProduct) error {
	return s.storage.CancelReservation(wp)
}

func (s *productService) Transfer(td *domain.TransferProduct) error {
	return s.storage.Transfer(td)
}

func (s *productService) Add(ad *domain.AddProduct) error {
	return s.storage.Add(ad)
}

func (s *productService) Delete(dp *domain.DeleteProduct) (*domain.Product, error) {
	return s.storage.Delete(dp)
}
