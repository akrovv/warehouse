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

type warehouseTestCase struct {
	in           []domain.Warehouse
	out          []domain.Warehouse
	err          error
	repeat       uint8
	repeatError  uint8
	expectResult []domain.Warehouse
}

type getLeftOversTestCase struct {
	in           domain.GetFromWarehouse
	out          []domain.Product
	err          error
	expectResult []domain.Product
}

func TestWarehouseCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wh := mocks.NewMockWarehouseService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := []domain.Warehouse{
		{
			Name:         "test-1",
			Availability: true,
		},
		{
			Name:         "test-2",
			Availability: false,
		},
	}

	expectWarehouse := [][]domain.Warehouse{
		in,
		{
			in[1],
		},
		nil,
	}

	testCases := []warehouseTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			repeat:       2,
			expectResult: expectWarehouse[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  1,
			expectResult: expectWarehouse[1],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			repeat:       2,
			repeatError:  2,
			expectResult: expectWarehouse[2],
		},
	}

	handler := NewWarehouseHandler(wh, logger)

	for index, tc := range testCases {
		for i := 0; i < int(tc.repeat); i++ {
			if tc.repeatError > 0 {
				wh.EXPECT().Create(&tc.in[i]).Return(tc.err)
				tc.repeatError--
				continue
			} else {
				tc.err = nil
			}
			wh.EXPECT().Create(&tc.in[i]).Return(tc.err)
		}

		err = handler.Create(tc.in, &tc.out)
		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, expectWarehouse[index]) {
			t.Fatalf("expected: %v, got: %v", expectWarehouse[index], tc.out)
		}
	}
}

func TestWarehouseGetLeftOvers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wh := mocks.NewMockWarehouseService(ctrl)
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("can't create logger: %s", err)
	}

	in := domain.GetFromWarehouse{
		WarehouseID: 1,
	}
	product := domain.Product{
		Name:     "test",
		Size:     "test",
		Code:     "test",
		Quantity: 10,
	}
	expectedResults := [][]domain.Product{
		{
			product,
		},
		{},
	}
	testCases := []getLeftOversTestCase{
		{
			in:           in,
			out:          nil,
			err:          nil,
			expectResult: expectedResults[0],
		},
		{
			in:           in,
			out:          nil,
			err:          domain.ErrTest,
			expectResult: nil,
		},
	}

	handler := NewWarehouseHandler(wh, logger)

	for _, tc := range testCases {
		wh.EXPECT().GetLeftOvers(&tc.in).Return(tc.expectResult, tc.err)

		err = handler.GetLeftOvers(tc.in, &tc.out)

		if !errors.Is(err, tc.err) {
			t.Fatalf("expected error: %v, got: %v", tc.err, err)
		}

		if !reflect.DeepEqual(tc.out, tc.expectResult) {
			t.Fatalf("expected: %v, got: %v", tc.expectResult, tc.out)
		}
	}
}
