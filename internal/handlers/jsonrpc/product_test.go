package jsonrpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/akrovv/warehouse/internal/domain"
	"github.com/akrovv/warehouse/internal/services/mocks"
	"github.com/akrovv/warehouse/pkg/logger"
	"github.com/golang/mock/gomock"
)

type productTestCase struct {
	in           []domain.Product
	out          []domain.Product
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.Product
}

type warehouseProductTestCase struct {
	in           []domain.WarehouseProduct
	out          []domain.WarehouseProduct
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.WarehouseProduct
}

type transferProductTestCase struct {
	in           []domain.TransferProduct
	out          []domain.TransferProduct
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.TransferProduct
}

type addProductTestCase struct {
	in           []domain.AddProduct
	out          []domain.AddProduct
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.AddProduct
}

type deleteProductTestCase struct {
	in           []domain.DeleteProduct
	out          []domain.Product
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.Product
}

func getWarehouseProductTestData(status string) []warehouseProductTestCase {
	in := []domain.WarehouseProduct{
		{
			WarehouseID: 1,
			Code:        "test",
			Quantity:    1,
		},
		{
			WarehouseID: 2,
			Code:        "test",
			Quantity:    2,
		},
	}

	expectWarehouseProduct := [][]domain.WarehouseProduct{
		{
			{
				WarehouseID: 1,
				Code:        "test",
				Quantity:    1,
				Status:      status,
			},
			{
				WarehouseID: 2,
				Code:        "test",
				Quantity:    2,
				Status:      status,
			},
		},
		{
			{
				WarehouseID: 2,
				Code:        "test",
				Quantity:    2,
				Status:      status,
			},
		},
	}

	testCases := []warehouseProductTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectWarehouseProduct[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectWarehouseProduct[1],
		},
		{
			in:          in,
			out:         nil,
			err:         domain.ErrTest,
			repeat:      2,
			repeatError: 2,
		},
	}

	return testCases
}

func TestProductCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := []domain.Product{
		{
			Name:     "test-1",
			Size:     "test-1",
			Code:     "test-1",
			Quantity: 10,
		},
		{
			Name:     "test-2",
			Size:     "test-2",
			Code:     "test-2",
			Quantity: 10,
		},
	}

	expectProduct := [][]domain.Product{
		in,
		{
			in[1],
		},
		nil,
	}

	testCases := []productTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectProduct[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectProduct[1],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  2,
			expectResult: expectProduct[2],
		},
	}

	handler := NewProductHandler(ps, logger)

	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().Create(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().Create(&tc.in[i]).Return(tc.err)
		}

		err = handler.Create(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}

func TestProductReserve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	testCases := getWarehouseProductTestData("reserved")
	handler := NewProductHandler(ps, logger)
	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().Reserve(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().Reserve(&tc.in[i]).Return(tc.err)
		}

		err = handler.Reserve(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}

func TestProductCancelReservation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	testCases := getWarehouseProductTestData("canceled")
	handler := NewProductHandler(ps, logger)
	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().CancelReservation(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().CancelReservation(&tc.in[i]).Return(tc.err)
		}

		err = handler.CancelReservation(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}

func TestProductTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := []domain.TransferProduct{
		{
			WarehouseFromID: 1,
			WarehouseToID:   2,
			Code:            "test",
			Quantity:        10,
		},
		{
			WarehouseFromID: 2,
			WarehouseToID:   3,
			Code:            "test",
			Quantity:        10,
		},
	}

	expectWarehouseProduct := [][]domain.TransferProduct{
		in,
		{
			in[1],
		},
	}

	testCases := []transferProductTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectWarehouseProduct[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectWarehouseProduct[1],
		},
		{
			in:          in,
			out:         nil,
			err:         domain.ErrTest,
			repeat:      2,
			repeatError: 2,
		},
	}

	handler := NewProductHandler(ps, logger)
	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().Transfer(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().Transfer(&tc.in[i]).Return(tc.err)
		}

		err = handler.Transfer(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}

func TestProductAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := []domain.AddProduct{
		{
			Code:        "test",
			Quantity:    10,
			WarehouseID: 1,
		},
		{
			Code:        "test",
			Quantity:    5,
			WarehouseID: 2,
		},
	}

	expectWarehouseProduct := [][]domain.AddProduct{
		in,
		{
			in[1],
		},
	}

	testCases := []addProductTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectWarehouseProduct[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectWarehouseProduct[1],
		},
		{
			in:          in,
			out:         nil,
			err:         domain.ErrTest,
			repeat:      2,
			repeatError: 2,
		},
	}

	handler := NewProductHandler(ps, logger)
	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().Add(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().Add(&tc.in[i]).Return(tc.err)
		}

		err = handler.Add(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}

func TestProductDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := mocks.NewMockProductService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := []domain.DeleteProduct{
		{
			Code: "test-1",
		},
		{
			Code: "test-2",
		},
	}

	product := []domain.Product{
		{
			Name:     "test",
			Size:     "test",
			Code:     "test-1",
			Quantity: 5,
		},
		{
			Name:     "test",
			Size:     "test",
			Code:     "test-2",
			Quantity: 1,
		},
	}

	expectWarehouseProduct := [][]domain.Product{
		product,
		{
			product[1],
		},
	}

	testCases := []deleteProductTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectWarehouseProduct[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectWarehouseProduct[1],
		},
		{
			in:          in,
			out:         nil,
			err:         domain.ErrTest,
			repeat:      2,
			repeatError: 2,
		},
	}

	handler := NewProductHandler(ps, logger)
	for _, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				ps.EXPECT().Delete(&tc.in[i]).Return(&product[i], tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			ps.EXPECT().Delete(&tc.in[i]).Return(&product[i], tc.err)
		}

		err = handler.Delete(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}
