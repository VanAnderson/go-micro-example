package inventory_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/db"
	"github.com/sksmith/go-micro-example/db/invrepo"
	"github.com/sksmith/go-micro-example/queue"
	"github.com/sksmith/go-micro-example/test"
)

func TestMain(m *testing.M) {
	test.ConfigLogging()
	os.Exit(m.Run())
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name string

		product inventory.Product

		getProductFunc           func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error)
		saveProductFunc          func(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error
		saveProductInventoryFunc func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error

		beginTransactionFunc func(ctx context.Context) (core.Transaction, error)
		commitFunc           func(ctx context.Context) error

		wantRepoCallCnt map[string]int
		wantTxCallCnt   map[string]int
		wantErr         bool
	}{
		{
			name:    "new product and inventory are saved",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			wantRepoCallCnt: map[string]int{"SaveProduct": 1, "SaveProductInventory": 1},
			wantTxCallCnt:   map[string]int{"Commit": 1, "Rollback": 0},
			wantErr:         false,
		},
		{
			name:    "product already exists",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			getProductFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"}, nil
			},

			wantRepoCallCnt: map[string]int{"SaveProduct": 0, "SaveProductInventory": 0},
			wantTxCallCnt:   map[string]int{"Commit": 0, "Rollback": 0},
			wantErr:         false,
		},
		{
			name:    "unexpected error getting product",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			getProductFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return inventory.Product{}, errors.New("some unexpected error")
			},

			wantRepoCallCnt: map[string]int{"SaveProduct": 0, "SaveProductInventory": 0},
			wantTxCallCnt:   map[string]int{"Commit": 0, "Rollback": 0},
			wantErr:         true,
		},
		{
			name:    "unexpected error saving product",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			saveProductFunc: func(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt: map[string]int{"SaveProduct": 1, "SaveProductInventory": 0},
			wantTxCallCnt:   map[string]int{"Commit": 0, "Rollback": 1},
			wantErr:         true,
		},
		{
			name:    "unexpected error saving product inventory",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			saveProductInventoryFunc: func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt: map[string]int{"SaveProduct": 1, "SaveProductInventory": 1},
			wantTxCallCnt:   map[string]int{"Commit": 0, "Rollback": 1},
			wantErr:         true,
		},
		{
			name:    "unexpected error comitting",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			beginTransactionFunc: func(ctx context.Context) (core.Transaction, error) { return nil, errors.New("some unexpected error") },

			wantRepoCallCnt: map[string]int{"SaveProduct": 1, "SaveProductInventory": 1},
			wantTxCallCnt:   map[string]int{"Commit": 1, "Rollback": 1},
			wantErr:         true,
		},
		{
			name:    "unexpected error beginning transaction",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			commitFunc: func(ctx context.Context) error { return errors.New("some unexpected error") },

			wantRepoCallCnt: map[string]int{"SaveProduct": 1, "SaveProductInventory": 1},
			wantTxCallCnt:   map[string]int{"Commit": 1, "Rollback": 1},
			wantErr:         true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getProductFunc != nil {
			mockRepo.GetProductFunc = test.getProductFunc
		} else {
			mockRepo.GetProductFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return inventory.Product{}, core.ErrNotFound
			}
		}
		if test.saveProductFunc != nil {
			mockRepo.SaveProductFunc = test.saveProductFunc
		}
		if test.saveProductInventoryFunc != nil {
			mockRepo.SaveProductInventoryFunc = test.saveProductInventoryFunc
		}

		mockTx := db.NewMockTransaction()
		if test.beginTransactionFunc != nil {
			mockRepo.BeginTransactionFunc = test.beginTransactionFunc
		} else {
			mockRepo.BeginTransactionFunc = func(ctx context.Context) (core.Transaction, error) {
				return mockTx, nil
			}
		}

		if test.commitFunc != nil {
			mockTx.CommitFunc = test.commitFunc
		}

		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			err := service.CreateProduct(context.Background(), test.product)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
			for f, c := range test.wantTxCallCnt {
				mockTx.VerifyCount(f, c, t)
			}
		})
	}
}

func TestProduce(t *testing.T) {
	product := inventory.Product{Sku: "somesku", Upc: "someupc", Name: "somename"}
	var productInventory *inventory.ProductInventory

	tests := []struct {
		name    string
		request inventory.ProductionRequest

		getProductionEventByRequestIDFunc func(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error)
		saveProductionEventFunc           func(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error
		getProductInventoryFunc           func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error)
		saveProductInventoryFunc          func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error

		publishInventoryFunc   func(ctx context.Context, productInventory inventory.ProductInventory) error
		publishReservationFunc func(ctx context.Context, reservation inventory.Reservation) error

		beginTransactionFunc func(ctx context.Context) (core.Transaction, error)
		commitFunc           func(ctx context.Context) error

		wantRepoCallCnt  map[string]int
		wantQueueCallCnt map[string]int
		wantTxCallCnt    map[string]int
		wantAvailable    int64
		wantErr          bool
	}{
		{
			name:    "inventory is incremented",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 1, "SaveProductInventory": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 1, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 2, "Rollback": 0},
			wantAvailable:    2,
		},
		{
			name:    "cannot produce zero",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 0},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "cannot produce negative",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: -1},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "request id is required",
			request: inventory.ProductionRequest{RequestID: "", Quantity: 1},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "production event already exists",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			getProductionEventByRequestIDFunc: func(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error) {
				return inventory.ProductionEvent{RequestID: "somerequestid", Quantity: 1}, nil
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
		},
		{
			name:    "unexpected error getting production event",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			getProductionEventByRequestIDFunc: func(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error) {
				return inventory.ProductionEvent{}, errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "unexpected error beginning transaction",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			beginTransactionFunc: func(ctx context.Context) (core.Transaction, error) {
				return nil, errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 0, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "unexpected error saving production event",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			saveProductionEventFunc: func(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 1, "SaveProductInventory": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "unexpected error saving product inventory",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			saveProductInventoryFunc: func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 1, "SaveProductInventory": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
			wantAvailable:    1,
			wantErr:          true,
		},
		{
			name:    "unexpected error comitting",
			request: inventory.ProductionRequest{RequestID: "somerequestid", Quantity: 1},

			commitFunc: func(ctx context.Context) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveProductionEvent": 1, "SaveProductInventory": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 1, "Rollback": 1},
			wantAvailable:    2,
			wantErr:          true,
		},
	}

	for _, test := range tests {
		productInventory = &inventory.ProductInventory{Product: product, Available: 1}

		mockTx := db.NewMockTransaction()
		if test.commitFunc != nil {
			mockTx.CommitFunc = test.commitFunc
		}

		mockRepo := invrepo.NewMockRepo()
		if test.beginTransactionFunc != nil {
			mockRepo.BeginTransactionFunc = test.beginTransactionFunc
		} else {
			mockRepo.BeginTransactionFunc = func(ctx context.Context) (core.Transaction, error) {
				return mockTx, nil
			}
		}
		if test.getProductionEventByRequestIDFunc != nil {
			mockRepo.GetProductionEventByRequestIDFunc = test.getProductionEventByRequestIDFunc
		}
		if test.saveProductionEventFunc != nil {
			mockRepo.SaveProductionEventFunc = test.saveProductionEventFunc
		}
		if test.getProductInventoryFunc != nil {
			mockRepo.GetProductInventoryFunc = test.getProductInventoryFunc
		} else {
			mockRepo.GetProductInventoryFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return *productInventory, nil
			}
		}
		if test.saveProductInventoryFunc != nil {
			mockRepo.SaveProductInventoryFunc = test.saveProductInventoryFunc
		} else {
			mockRepo.SaveProductInventoryFunc = func(ctx context.Context, pi inventory.ProductInventory, options ...core.UpdateOptions) error {
				productInventory = &pi
				return nil
			}
		}

		mockQueue := queue.NewMockQueue()
		if test.publishInventoryFunc != nil {
			mockQueue.PublishInventoryFunc = test.publishInventoryFunc
		}
		if test.publishReservationFunc != nil {
			mockQueue.PublishReservationFunc = test.publishReservationFunc
		}

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			err := service.Produce(context.Background(), product, test.request)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if productInventory.Available != test.wantAvailable {
				t.Errorf("unexpected available got=%d want=%d", productInventory.Available, test.wantAvailable)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
			for f, c := range test.wantQueueCallCnt {
				mockQueue.VerifyCount(f, c, t)
			}
			for f, c := range test.wantTxCallCnt {
				mockTx.VerifyCount(f, c, t)
			}
		})
	}
}

func TestReserve(t *testing.T) {
	tests := []struct {
		name    string
		request inventory.ReservationRequest

		getProductFunc                func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error)
		getReservationByRequestIDFunc func(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error)
		saveReservationFunc           func(ctx context.Context, reservation *inventory.Reservation, options ...core.UpdateOptions) error

		beginTransactionFunc func(ctx context.Context) (core.Transaction, error)
		commitFunc           func(ctx context.Context) error

		wantRepoCallCnt  map[string]int
		wantQueueCallCnt map[string]int
		wantTxCallCnt    map[string]int
		wantState        inventory.ReserveState
		wantErr          bool
	}{
		{
			name:    "reservation is created",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 2, "Rollback": 0},
			wantState:        inventory.Open,
		},
		{
			name:            "reservation request id is required",
			request:         inventory.ReservationRequest{Sku: "somesku", Requester: "somerequester", Quantity: 1},
			wantRepoCallCnt: map[string]int{"SaveReservation": 0},
			wantErr:         true,
		},
		{
			name:            "reservation sku is required",
			request:         inventory.ReservationRequest{RequestID: "somerequestid", Requester: "somerequester", Quantity: 1},
			wantRepoCallCnt: map[string]int{"SaveReservation": 0},
			wantErr:         true,
		},
		{
			name:            "reservation requester is required",
			request:         inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Quantity: 1},
			wantRepoCallCnt: map[string]int{"SaveReservation": 0},
			wantErr:         true,
		},
		{
			name:            "reservation quantity must be greater than zero",
			request:         inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 0},
			wantRepoCallCnt: map[string]int{"SaveReservation": 0},
			wantErr:         true,
		},
		{
			name:            "reservation quantity must not be negative",
			request:         inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: -1},
			wantRepoCallCnt: map[string]int{"SaveReservation": 0},
			wantErr:         true,
		},
		{
			name:    "unexpected error beginning transaction",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			beginTransactionFunc: func(ctx context.Context) (core.Transaction, error) {
				return nil, errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 0},
			wantErr:          true,
		},
		{
			name:    "unexpected error getting product",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			getProductFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return inventory.Product{}, errors.New("unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
			wantErr:          true,
		},
		{
			name:    "reservation request has already been processed",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			getReservationByRequestIDFunc: func(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error) {
				return inventory.Reservation{RequestID: "somerequestid"}, nil
			},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 0},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
			wantErr:          false,
		},
		{
			name:    "unexpected error saving reservation",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			saveReservationFunc: func(ctx context.Context, reservation *inventory.Reservation, options ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
			wantErr:          true,
		},
		{
			name:    "unexpected error comitting",
			request: inventory.ReservationRequest{RequestID: "somerequestid", Sku: "somesku", Requester: "somerequester", Quantity: 1},

			commitFunc: func(ctx context.Context) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt:  map[string]int{"SaveReservation": 1},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantTxCallCnt:    map[string]int{"Commit": 1, "Rollback": 1},
			wantErr:          true,
		},
	}

	for _, test := range tests {
		mockTx := db.NewMockTransaction()
		if test.commitFunc != nil {
			mockTx.CommitFunc = test.commitFunc
		}

		mockRepo := invrepo.NewMockRepo()
		if test.beginTransactionFunc != nil {
			mockRepo.BeginTransactionFunc = test.beginTransactionFunc
		} else {
			mockRepo.BeginTransactionFunc = func(ctx context.Context) (core.Transaction, error) {
				return mockTx, nil
			}
		}
		if test.getProductFunc != nil {
			mockRepo.GetProductFunc = test.getProductFunc
		}
		if test.getReservationByRequestIDFunc != nil {
			mockRepo.GetReservationByRequestIDFunc = test.getReservationByRequestIDFunc
		} else {
			mockRepo.GetReservationByRequestIDFunc = func(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error) {
				return inventory.Reservation{}, core.ErrNotFound
			}
		}
		if test.saveReservationFunc != nil {
			mockRepo.SaveReservationFunc = test.saveReservationFunc
		}

		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.Reserve(context.Background(), test.request)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if res.State != test.wantState {
				t.Errorf("unexpected state got=%s want=%s", res.State, test.wantState)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
			for f, c := range test.wantQueueCallCnt {
				mockQueue.VerifyCount(f, c, t)
			}
			for f, c := range test.wantTxCallCnt {
				mockTx.VerifyCount(f, c, t)
			}
		})
	}
}

func TestGetAllProductInventory(t *testing.T) {
	productInv := getProductInventory()
	tests := []struct {
		name   string
		limit  int
		offset int

		getAllProductInventoryFunc func(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error)

		wantProductInventory []inventory.ProductInventory
		wantErr              bool
	}{
		{
			name:                 "product is returned",
			wantProductInventory: productInv,
		},
		{
			name: "error is returned",
			getAllProductInventoryFunc: func(ctx context.Context, limit, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error) {
				return []inventory.ProductInventory{}, errors.New("some unexpected error")
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getAllProductInventoryFunc != nil {
			mockRepo.GetAllProductInventoryFunc = test.getAllProductInventoryFunc
		} else {
			mockRepo.GetAllProductInventoryFunc = func(ctx context.Context, limit, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error) {
				return productInv, nil
			}
		}
		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.GetAllProductInventory(context.Background(), test.limit, test.offset)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if len(res) != len(test.wantProductInventory) {
				t.Errorf("unexpected product inventory got=%v want=%v", res, test.wantProductInventory)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	productInv := getProductInventory()
	tests := []struct {
		name   string
		limit  int
		offset int

		getProductFunc func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error)

		wantProduct inventory.Product
		wantErr     bool
	}{
		{
			name:        "product is returned",
			wantProduct: productInv[0].Product,
		},
		{
			name: "error is returned",
			getProductFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return inventory.Product{}, errors.New("some unexpected error")
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getProductFunc != nil {
			mockRepo.GetProductFunc = test.getProductFunc
		} else {
			mockRepo.GetProductFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
				return productInv[0].Product, nil
			}
		}
		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.GetProduct(context.Background(), "sku1")
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(res, test.wantProduct) {
				t.Errorf("unexpected product inventory got=%v want=%v", res, test.wantProduct)
			}
		})
	}
}

func TestGetProductInventory(t *testing.T) {
	productInv := getProductInventory()
	tests := []struct {
		name   string
		limit  int
		offset int

		getProductInventoryFunc func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error)

		wantProductInv inventory.ProductInventory
		wantErr        bool
	}{
		{
			name:           "product is returned",
			wantProductInv: productInv[0],
		},
		{
			name: "error is returned",
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{}, errors.New("some unexpected error")
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getProductInventoryFunc != nil {
			mockRepo.GetProductInventoryFunc = test.getProductInventoryFunc
		} else {
			mockRepo.GetProductInventoryFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error) {
				return productInv[0], nil
			}
		}
		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.GetProductInventory(context.Background(), "sku1")

			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(res, test.wantProductInv) {
				t.Errorf("unexpected product inventory got=%v want=%v", res, test.wantProductInv)
			}
		})
	}
}

func getProductInventory() []inventory.ProductInventory {
	return []inventory.ProductInventory{
		{Product: inventory.Product{Sku: "sku1", Upc: "upc1", Name: "name1"}, Available: 1},
		{Product: inventory.Product{Sku: "sku2", Upc: "upc2", Name: "name2"}, Available: 10},
		{Product: inventory.Product{Sku: "sku3", Upc: "upc3", Name: "name3"}, Available: 0},
	}
}
