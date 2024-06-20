package postgresql

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"

	"github.com/akrovv/warehouse/internal/domain"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type productTestCase struct {
	product  domain.Product
	wp       domain.WarehouseProduct
	dp       domain.DeleteProduct
	query    string
	rows     *sqlmock.Rows
	args     []driver.Value
	returned driver.Result
	result   error
	isError  bool
}

type productTransactTestCase struct {
	ad                     domain.AddProduct
	td                     domain.TransferProduct
	expectedQuery          string
	expectedSelect         string
	expectedSelectError    error
	selectRows             *sqlmock.Rows
	expectedExecQuery      string
	expectedExecQueryError error
	execArgs               []driver.Value
	execResult             driver.Result
	execError              error
	expectCommit           bool
	expectError            bool
	beginError             error
}

func TestProductCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewProductStorage(db)
	product := domain.Product{
		Name:     "test-1",
		Size:     "test-1",
		Code:     "test-1",
		Quantity: 10,
	}
	query := `INSERT INTO products \(name, size, code, quantity\) VALUES \(\$1, \$2, \$3, \$4\)`
	args := []driver.Value{"test-1", "test-1", "test-1", 10}

	testCases := []productTestCase{
		{
			product:  product,
			query:    query,
			args:     args,
			result:   nil,
			returned: sqlmock.NewResult(0, 1),
		},
		{
			product:  product,
			query:    query,
			args:     args,
			result:   domain.ErrTest,
			returned: sqlmock.NewResult(0, 0),
		},
		{
			product:  product,
			query:    query,
			args:     args,
			returned: sqlmock.NewErrorResult(domain.ErrTest),
			isError:  true,
		},
		{
			product:  product,
			query:    query,
			args:     args,
			returned: sqlmock.NewResult(0, 0),
			isError:  true,
		},
	}

	for _, tc := range testCases {
		mock.ExpectExec(tc.query).
			WithArgs(tc.args...).
			WillReturnResult(tc.returned).
			WillReturnError(tc.result)

		err = storage.Create(&tc.product)

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

func TestProductReserve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewProductStorage(db)
	wp := domain.WarehouseProduct{
		WarehouseID: 10,
		Code:        "test-1",
		Quantity:    10,
	}

	query := `
	UPDATE warehouse_products
	SET available_quantity = available_quantity - \$3,
		reserved_quantity = reserved_quantity \+ \$3
	WHERE warehouse_id = \$1 AND product_code = \$2`

	args := []driver.Value{10, "test-1", 10}
	testCases := []productTestCase{
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewResult(0, 1),
			result:   nil,
		},
		{
			wp:      wp,
			query:   query,
			args:    args,
			result:  domain.ErrTest,
			isError: true,
		},
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewResult(0, 0),
			isError:  true,
		},
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewErrorResult(domain.ErrTest),
			isError:  true,
		},
	}

	for _, tc := range testCases {
		mock.ExpectExec(tc.query).
			WithArgs(tc.args...).
			WillReturnResult(tc.returned).
			WillReturnError(tc.result)

		err = storage.Reserve(&tc.wp)

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

func TestProductCancelReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewProductStorage(db)
	wp := domain.WarehouseProduct{
		WarehouseID: 10,
		Code:        "test-1",
		Quantity:    10,
	}

	query := `
	UPDATE warehouse_products
	SET available_quantity = available_quantity \+ \$3,
		reserved_quantity = reserved_quantity - \$3
	WHERE warehouse_id = \$1 AND product_code = \$2`
	args := []driver.Value{10, "test-1", 10}
	testCases := []productTestCase{
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewResult(0, 1),
			result:   nil,
		},
		{
			wp:      wp,
			query:   query,
			args:    args,
			result:  domain.ErrTest,
			isError: true,
		},
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewResult(0, 0),
			isError:  true,
		},
		{
			wp:       wp,
			query:    query,
			args:     args,
			returned: sqlmock.NewErrorResult(domain.ErrTest),
			isError:  true,
		},
	}

	for _, tc := range testCases {
		mock.ExpectExec(tc.query).
			WithArgs(tc.args...).
			WillReturnResult(tc.returned).
			WillReturnError(tc.result)

		err = storage.CancelReservation(&tc.wp)

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

func TestProductTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := productStorage{db: db}
	td := domain.TransferProduct{
		WarehouseFromID: 1,
		WarehouseToID:   2,
		Code:            "test-1",
		Quantity:        5,
	}

	expectedSelect := "SELECT available_quantity FROM warehouse_products"
	expectedQuery := "INSERT INTO warehouse_products"
	expectedExecQuery := "UPDATE warehouse_products"

	testCases := []productTransactTestCase{
		{
			td:                td,
			expectedExecQuery: expectedExecQuery,
			expectedSelect:    expectedSelect,
			selectRows:        sqlmock.NewRows([]string{"available_quantity"}).AddRow(10),
			expectedQuery:     expectedQuery,
			execArgs:          []driver.Value{2, "test-1", 5, 0},
			execResult:        sqlmock.NewResult(0, 1),
			expectCommit:      true,
		},
		{
			td:          td,
			expectError: true,
			beginError:  domain.ErrTest,
		},
		{
			td:                  td,
			expectedSelect:      expectedSelect,
			expectError:         true,
			expectedSelectError: domain.ErrTest,
		},
		{
			td:                     td,
			expectedSelect:         expectedSelect,
			expectedExecQuery:      expectedExecQuery,
			selectRows:             sqlmock.NewRows([]string{"available_quantity"}).AddRow(10),
			expectError:            true,
			expectedExecQueryError: domain.ErrTest,
		},
		{
			td:             td,
			expectedSelect: expectedSelect,
			selectRows:     sqlmock.NewRows([]string{"available_quantity"}).AddRow(0),
			expectError:    true,
		},
		{
			td:                td,
			expectedExecQuery: expectedExecQuery,
			expectedSelect:    expectedSelect,
			selectRows:        sqlmock.NewRows([]string{"available_quantity"}).AddRow(10),
			expectedQuery:     expectedQuery,
			execArgs:          []driver.Value{2, "test-1", 5, 0},
			execResult:        sqlmock.NewResult(0, 1),
			execError:         domain.ErrTest,
			expectError:       true,
		},
	}

	for _, tc := range testCases {
		if tc.expectedSelect != "" {
			mock.ExpectBegin().WillReturnError(tc.beginError)

			mock.ExpectQuery(tc.expectedSelect).
				WithArgs(tc.td.WarehouseFromID, tc.td.Code).
				WillReturnRows(tc.selectRows).
				WillReturnError(tc.expectedSelectError)
		}

		if tc.expectedExecQuery != "" {
			mock.ExpectExec(tc.expectedExecQuery).
				WithArgs(tc.td.WarehouseFromID, tc.td.Code, tc.td.Quantity).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.expectedExecQueryError)
		}

		if tc.expectedQuery != "" {
			mock.ExpectExec(tc.expectedQuery).
				WithArgs(tc.execArgs...).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)
		}

		if tc.expectCommit {
			mock.ExpectCommit()
		}

		err = storage.Transfer(&tc.td)
		if (err != nil) != tc.expectError {
			t.Errorf("unexpected error: %v", err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestProductAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := productStorage{db: db}
	ad := domain.AddProduct{
		Code:        "test-1",
		Quantity:    10,
		WarehouseID: 1,
	}

	expectedQuery := "UPDATE warehouse_products"
	expectedExecQuery := "UPDATE products"

	testCases := []productTransactTestCase{
		{
			ad:                ad,
			expectedExecQuery: expectedExecQuery,
			expectedQuery:     expectedQuery,
			execArgs:          []driver.Value{1, "test-1", 10},
			execResult:        sqlmock.NewResult(0, 1),
			expectCommit:      true,
		},
		{
			ad:                     ad,
			expectedExecQuery:      expectedExecQuery,
			expectError:            true,
			expectedExecQueryError: domain.ErrTest,
		},
		{
			ad:                ad,
			expectedExecQuery: expectedExecQuery,
			expectedQuery:     expectedQuery,
			expectError:       true,
			execError:         domain.ErrTest,
			execArgs:          []driver.Value{1, "test-1", 10},
			execResult:        sqlmock.NewResult(0, 1),
		},
	}

	for _, tc := range testCases {
		if tc.expectedExecQuery != "" {
			mock.ExpectBegin().WillReturnError(tc.beginError)

			mock.ExpectExec(tc.expectedExecQuery).
				WithArgs(tc.ad.Quantity, tc.ad.Code).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.expectedExecQueryError)
		}

		if tc.expectedQuery != "" {
			mock.ExpectExec(tc.expectedQuery).
				WithArgs(tc.execArgs...).
				WillReturnResult(tc.execResult).
				WillReturnError(tc.execError)
		}

		if tc.expectCommit {
			mock.ExpectCommit()
		}

		err = storage.Add(&tc.ad)
		if (err != nil) != tc.expectError {
			t.Errorf("unexpected error: %v", err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestProductDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock: %s", err)
	}
	defer db.Close()

	storage := NewProductStorage(db)
	dp := domain.DeleteProduct{
		Code: "test",
	}

	query := `DELETE FROM products WHERE code = \$1
			  RETURNING name, size, code, quantity`
	rows := sqlmock.NewRows([]string{"name", "size", "code", "quantity"}).
		AddRow("test", "test", "test", 10)
	args := []driver.Value{"test"}
	expectedResult := &domain.Product{
		Name:     "test",
		Size:     "test",
		Code:     "test",
		Quantity: 10,
	}

	testCases := []productTestCase{
		{
			dp:     dp,
			query:  query,
			rows:   rows,
			args:   args,
			result: nil,
		},
		{
			dp:     dp,
			query:  query,
			args:   args,
			result: domain.ErrTest,
		},
	}

	for _, tc := range testCases {
		mock.ExpectQuery(tc.query).
			WithArgs(tc.args...).
			WillReturnRows(tc.rows).
			WillReturnError(tc.result)

		product, err := storage.Delete(&tc.dp)

		if !errors.Is(err, tc.result) {
			t.Errorf("expected: %v, got: %v", tc.result, err)
		}

		if tc.result == nil && !reflect.DeepEqual(product, expectedResult) {
			t.Fatalf("expected: %v, got: %v", expectedResult, product)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %s", err)
		}
	}
}
