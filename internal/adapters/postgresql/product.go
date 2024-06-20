package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/akrovv/warehouse/internal/domain"
)

type productStorage struct {
	db *sql.DB
}

func NewProductStorage(db *sql.DB) *productStorage {
	return &productStorage{
		db: db,
	}
}

func (s *productStorage) Create(product *domain.Product) error {
	res, err := s.db.Exec("INSERT INTO products (name, size, code, quantity) VALUES ($1, $2, $3, $4)",
		product.Name, product.Size, product.Code, product.Quantity)

	if err != nil {
		return fmt.Errorf("db.Exec with command INSERT to products returned: %w", err)
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

func (s *productStorage) Reserve(wp *domain.WarehouseProduct) error {
	res, err := s.db.Exec(`
		UPDATE warehouse_products 
		SET available_quantity = available_quantity - $3, 
			reserved_quantity = reserved_quantity + $3 
		WHERE warehouse_id = $1 AND product_code = $2`,
		wp.WarehouseID, wp.Code, wp.Quantity)

	if err != nil {
		return fmt.Errorf("db.Exec with command UPDATE to warehouse_products returned: %w", err)
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

func (s *productStorage) CancelReservation(wp *domain.WarehouseProduct) error {
	res, err := s.db.Exec(`
		UPDATE warehouse_products 
		SET available_quantity = available_quantity + $3, 
			reserved_quantity = reserved_quantity - $3
		WHERE warehouse_id = $1 AND product_code = $2
		`, wp.WarehouseID, wp.Code, wp.Quantity)

	if err != nil {
		return fmt.Errorf("db.Exec with command UPDATE to warehouse_products returned: %w", err)
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

func (s *productStorage) Transfer(td *domain.TransferProduct) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin() returned: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
		_ = tx.Commit()
	}()

	var quantity uint64

	err = tx.QueryRow(`SELECT available_quantity FROM warehouse_products
						WHERE warehouse_id = $1 AND product_code = $2`,
		td.WarehouseFromID, td.Code).Scan(&quantity)
	if err != nil {
		return fmt.Errorf("db.QueryRow with command SELECT to warehouse_products returned: %w", err)
	}

	if quantity < td.Quantity {
		return fmt.Errorf("not enough quantity: %d, in warehouse: %d. available: %d", td.Quantity, td.WarehouseFromID, quantity)
	}

	res, err := tx.Exec(`UPDATE warehouse_products 
					SET available_quantity = available_quantity - $3
					WHERE warehouse_id = $1 AND product_code = $2 `,
		td.WarehouseFromID, td.Code, td.Quantity)
	if err != nil {
		return fmt.Errorf("db.Exec with command UPDATE to warehouse_products returned: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows.RowsAffected() returned: %w", err)
	}

	if affected == 0 {
		return errors.New("affected 0 rows")
	}

	res, err = tx.Exec(`INSERT INTO warehouse_products (warehouse_id, product_code, available_quantity, reserved_quantity)
					VALUES ($1, $2, $3, $4)
					ON CONFLICT (warehouse_id, product_code) DO UPDATE
					SET available_quantity = warehouse_products.available_quantity + EXCLUDED.available_quantity`,
		td.WarehouseToID, td.Code, td.Quantity, 0)

	if err != nil {
		return fmt.Errorf("db.Exec with command INSERT/UPDATE to warehouse_products returned: %w", err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows.RowsAffected() returned: %w", err)
	}

	if affected == 0 {
		return errors.New("affected 0 rows")
	}

	return nil
}

func (s *productStorage) Add(ad *domain.AddProduct) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin() returned: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	res, err := tx.Exec(`UPDATE products SET quantity = quantity + $1 WHERE code = $2`,
		ad.Quantity, ad.Code)

	if err != nil {
		return fmt.Errorf("db.Exec with command UPDATE to products returned: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows.RowsAffected() returned: %w", err)
	}

	if affected == 0 {
		return errors.New("affected 0 rows")
	}

	res, err = tx.Exec(`
		UPDATE warehouse_products
		SET available_quantity = available_quantity + $3
		WHERE warehouse_id = $1 AND product_code = $2`,
		ad.WarehouseID, ad.Code, ad.Quantity)

	if err != nil {
		return fmt.Errorf("db.Exec with command UPDATE to warehouse_products returned: %w", err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows.RowsAffected() returned: %w", err)
	}

	if affected == 0 {
		return errors.New("affected 0 rows")
	}

	return nil
}

func (s *productStorage) Delete(dp *domain.DeleteProduct) (*domain.Product, error) {
	product := domain.Product{}

	err := s.db.QueryRow(`DELETE FROM products WHERE code = $1
						  RETURNING name, size, code, quantity`,
		dp.Code).
		Scan(&product.Name, &product.Size, &product.Code, &product.Quantity)

	if err != nil {
		return nil, fmt.Errorf("db.Exec with command DELETE to products returned: %w", err)
	}

	return &product, nil
}
