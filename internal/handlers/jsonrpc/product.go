package jsonrpc

import (
	"fmt"

	"github.com/akrovv/warehouse/internal/domain"
	"github.com/akrovv/warehouse/pkg/logger"
)

type productHandler struct {
	service ProductService
	logger  logger.Logger
}

func NewProductHandler(service ProductService, logger logger.Logger) *productHandler {
	return &productHandler{
		service: service,
		logger:  logger,
	}
}

func (h *productHandler) Create(in []domain.Product, out *[]domain.Product) error {
	var err error
	total := 0
	success := make([]domain.Product, 0, len(in))

	for _, value := range in {
		if err = h.service.Create(&value); err != nil {
			h.logger.Infof("error while creating product %v, error: %w\n", value, err)
			total++
			continue
		}
		success = append(success, value)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = success
	return nil
}

func (h *productHandler) Reserve(in []domain.WarehouseProduct, out *[]domain.WarehouseProduct) error {
	var err error
	total := 0
	reserved := make([]domain.WarehouseProduct, 0, len(in))

	for _, value := range in {
		if err = h.service.Reserve(&value); err != nil {
			h.logger.Infof("can't reserve item: %v, got error: %w", value, err)
			total++
			continue
		}

		value.Status = "reserved"
		reserved = append(reserved, value)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = reserved
	return nil
}

func (h *productHandler) CancelReservation(in []domain.WarehouseProduct, out *[]domain.WarehouseProduct) error {
	var err error
	total := 0
	unreserved := make([]domain.WarehouseProduct, 0, len(in))

	for _, value := range in {
		if err = h.service.CancelReservation(&value); err != nil {
			h.logger.Infof("can't cancel reservation with item: %v, got error: %w", value, err)
			total++
			continue
		}

		value.Status = "canceled"
		unreserved = append(unreserved, value)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = unreserved
	return nil
}

func (h *productHandler) Transfer(in []domain.TransferProduct, out *[]domain.TransferProduct) error {
	var err error
	total := 0
	transfered := make([]domain.TransferProduct, 0, len(in))
	for _, value := range in {
		if err = h.service.Transfer(&value); err != nil {
			h.logger.Infof("can't transfer item: %v, got error: %w", value, err)
			total++
			continue
		}

		transfered = append(transfered, value)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = transfered
	return nil
}

func (h *productHandler) Add(in []domain.AddProduct, out *[]domain.AddProduct) error {
	var err error
	total := 0
	added := make([]domain.AddProduct, 0, len(in))

	for _, value := range in {
		if err = h.service.Add(&value); err != nil {
			h.logger.Infof("can't add item: %v, got error: %w", value, err)
			total++
			continue
		}

		added = append(added, value)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = added
	return nil
}

func (h *productHandler) Delete(in []domain.DeleteProduct, out *[]domain.Product) error {
	var err error
	var product *domain.Product
	total := 0
	deleted := make([]domain.Product, 0, len(in))

	for _, value := range in {
		product, err = h.service.Delete(&value)
		if err != nil {
			h.logger.Infof("can't delete item: %v, got error: %w", value, err)
			total++
			continue
		}

		deleted = append(deleted, *product)
	}

	if total == len(in) {
		return fmt.Errorf("all calls returned: %w", err)
	}

	*out = deleted
	return nil
}
