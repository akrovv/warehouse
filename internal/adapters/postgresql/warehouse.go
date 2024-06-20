package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/akrovv/warehouse/internal/domain"
)

type warehouseStorage struct {
	db *sql.DB
}

func NewWarehouseStorage(db *sql.DB) *warehouseStorage {
	return &warehouseStorage{
		db: db,
	}
}

func (s *warehouseStorage) Create(warehouse *domain.Warehouse) error {
	res, err := s.db.Exec("INSERT INTO warehouses (name, availability) VALUES ($1, $2)", warehouse.Name, warehouse.Availability)

	if err != nil {
		return fmt.Errorf("db.Exec with command INSERT to warehouses returned: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows.RowsAffected() returned: %w", err)
	}

	if affected == 0 {
		return errors.New("affected 0 rows")
	}

	return nil
}

func (s *warehouseStorage) GetLeftOvers(gw *domain.GetFromWarehouse) ([]domain.Product, error) {
	product := domain.Product{}
	products := make([]domain.Product, 0, domain.BasicSliceLength)
	rows, err := s.db.Query(`SELECT p.name, size, code, available_quantity FROM warehouse_products wp 
							JOIN products p ON wp.product_code = p.code
							JOIN warehouses w ON wp.warehouse_id = w.id
						  	WHERE availability = true AND warehouse_id = $1 AND available_quantity > 0`,
		gw.WarehouseID)

	if err != nil {
		return nil, fmt.Errorf("db.Query with command SELECT to warehouse_products returned: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&product.Name, &product.Size, &product.Code, &product.Quantity)
		if err != nil {
			return nil, fmt.Errorf("row scan returned: %w", err)
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		return nil, sql.ErrNoRows
	}

	return products, nil
}
