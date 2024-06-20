package jsonrpc

import (
	"fmt"

	"github.com/akrovv/warehouse/internal/domain"
	"github.com/akrovv/warehouse/pkg/logger"
)

type warehouseHandler struct {
	service WarehouseService
	logger  logger.Logger
}

func NewWarehouseHandler(service WarehouseService, logger logger.Logger) *warehouseHandler {
	return &warehouseHandler{
		service: service,
		logger:  logger,
	}
}

func (h *warehouseHandler) Create(in []domain.Warehouse, out *[]domain.Warehouse) error {
	var err error
	total := 0
	success := make([]domain.Warehouse, 0, len(in))

	for _, value := range in {
		if err = h.service.Create(&value); err != nil {
			h.logger.Infof("error while creating product %v, error: %s\n", value, err.Error())
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

func (h *warehouseHandler) GetLeftOvers(in domain.GetFromWarehouse, out *[]domain.Product) error {
	products, err := h.service.GetLeftOvers(&in)

	if err != nil {
		return fmt.Errorf("service.GetLeftOvers returned: %w", err)
	}

	*out = products
	return nil
}
