
package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"gofr.dev/pkg/gofr"
)

type MockMain struct {
	mock.Mock
}

var Collection *MockMain

func init() {
	Collection = new(MockMain)
}

func (m *MockMain) AddEntryHandler(ctx *gofr.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
func (m *MockMain) UpdateEntryHandler(ctx *gofr.Context, requestBody []byte) (string, error) {
    args := m.Called(ctx, requestBody)
    return args.String(0), args.Error(1)
}

func (m *MockMain) DeleteEntryHandler(ctx *gofr.Context, requestBody []byte) (string, error) {
    args := m.Called(ctx, requestBody)
    return args.String(0), args.Error(1)
}
func TestAddEntryHandler_ValidData(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())

	Collection.On("AddEntryHandler", ctx).Return("someCarID", nil)

	carID, err := Collection.AddEntryHandler(ctx)

	if err != nil {
		t.Errorf("AddEntryHandler failed: %v", err)
	}

	if carID != "someCarID" {
		t.Errorf("Expected carID 'someCarID', got %s", carID)
	}

	Collection.AssertExpectations(t)
}

func TestAddEntryHandler_MissingLicensePlate(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())

	Collection.On("AddEntryHandler", ctx).Return("", errors.New("missing required field 'licensePlate'"))

	_, err := Collection.AddEntryHandler(ctx)

	if err == nil {
		t.Errorf("AddEntryHandler should have failed with missing license plate")
	} else if err.Error() != "missing required field 'licensePlate'" {
		t.Errorf("Unexpected error message: %v", err.Error())
	}

	Collection.AssertExpectations(t)
}

func TestUpdateEntryHandler_ValidData(t *testing.T) {
    ctx := gofr.NewContext(nil, nil, gofr.New())
    requestBody := []byte(`{"id": "9f3996d9-83e5-48da-ae43-f10e44edb51f", "status": "Updated"}`)
    Collection.On("UpdateEntryHandler", ctx, requestBody).Return("UpdatedCarID", nil)
    carID, err := Collection.UpdateEntryHandler(ctx, requestBody)
    if err != nil {
        t.Errorf("UpdateEntryHandler failed: %v", err)
    }
    Collection.AssertExpectations(t)
    if carID != "UpdatedCarID" {
        t.Errorf("Expected carID 'UpdatedCarID', got %s", carID)
    }
}
