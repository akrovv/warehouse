package postgresql

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/akrovv/warehouse/internal/domain"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type warehouseTestCase struct {
	warehouse domain.Warehouse
	gw        domain.GetFromWarehouse
	query     string
	rows      *sqlmock.Rows
	args      []driver.Value
	returned  driver.Result
	result    error
	isError   bool
}

func TestWarehouseCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewWarehouseStorage(db)
	warehouse := domain.Warehouse{
		Name:         "test-1",
		Availability: true,
	}
	query := `INSERT INTO warehouses \(name, availability\) VALUES \(\$1, \$2\)`
	args := []driver.Value{"test-1", true}

	testCases := []warehouseTestCase{
		{
			warehouse: warehouse,
			query:     query,
			args:      args,
			returned:  sqlmock.NewResult(0, 1),
			result:    nil,
		},
		{
			warehouse: warehouse,
			query:     query,
			args:      args,
			returned:  sqlmock.NewResult(0, 0),
			result:    domain.ErrTest,
		},
		{
			warehouse: warehouse,
			query:     query,
			args:      args,
			returned:  sqlmock.NewErrorResult(domain.ErrTest),
			isError:   true,
		},
		{
			warehouse: warehouse,
			query:     query,
			args:      args,
			returned:  sqlmock.NewResult(0, 0),
			isError:   true,
		},
	}

	for _, tc := range testCases {
		mock.ExpectExec(tc.query).
			WithArgs(tc.args...).
			WillReturnResult(tc.returned).
			WillReturnError(tc.result)

		err = storage.Create(&tc.warehouse)

		if !errors.Is(err, tc.result) {
			if tc.isError && err != nil {
				continue
			}
			t.Errorf("expected: %v, got: %v", tc.result, err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestWarehouseGetLeftOvers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewWarehouseStorage(db)
	gw := domain.GetFromWarehouse{
		WarehouseID: 1,
	}
	query := `SELECT p.name, size, code, available_quantity FROM warehouse_products wp 
	JOIN products p ON wp.product_code = p.code
	JOIN warehouses w ON wp.warehouse_id = w.id`
	rows := sqlmock.NewRows([]string{"name", "size", "code", "available_quantity"}).
		AddRow("test", "test", "test", 10)
	args := []driver.Value{1}
	expectedResult := []domain.Product{
		{
			Name:     "test",
			Size:     "test",
			Code:     "test",
			Quantity: 10,
		},
	}

	testCases := []warehouseTestCase{
		{
			gw:     gw,
			query:  query,
			rows:   rows,
			args:   args,
			result: nil,
		},
		{
			gw:     gw,
			query:  query,
			args:   args,
			result: domain.ErrTest,
		},
		{
			gw:      gw,
			query:   query,
			rows:    sqlmock.NewRows([]string{"name"}).AddRow("test"),
			args:    args,
			result:  nil,
			isError: true,
		},
	}

	for _, tc := range testCases {
		mock.ExpectQuery(tc.query).
			WithArgs(tc.args...).
			WillReturnRows(tc.rows).
			WillReturnError(tc.result)

		products, err := storage.GetLeftOvers(&tc.gw)
		if !errors.Is(err, tc.result) {
			if tc.isError && err != nil {
				continue
			}
			t.Errorf("expected: %v, got: %v", tc.result, err)
		}

		if tc.result == nil && !reflect.DeepEqual(products, expectedResult) {
			t.Fatalf("expected: %v, got: %v", expectedResult, products)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %s", err)
		}
	}
}
