I'm going to give you a file or list of files in the codebase for a golang inventory management micro service. Your sole task is to generate extractive questions and answers in the format that could be used for fine tuning in chatGPT. Everything that follows here before I give you the code will be instructions and guidance on performing this task. I want you to give me the output in pure JSONL with no further commentary. Here is the format of each JSONL line I want you to give me:
```
  {"messages": [{ "role": "system", "content": "You are an experienced golang developer with specific knowledge of https://github.com/VanAnderson/go-micro-example. You can act as an assistant to provide information about this codebase and also assist with technical discussions and creating detailed work tickets for developers. You will give answers specifically in the context of the go-micro-example codebase and should assume all questions pertain to this codebase unless otherwise specified." },{ "role": "user", "content": "Tell me about authentication in go-micro-example" },{ "role": "assistant", "content": " Users are stored in the database along with their hashed password. Users are locally cached using [golang-lru](https://github.com/hashicorp/golang-lru)." }] }
```
The first part:
```
{ "role\": \"system\", \"content\": \"You are an experienced golang developer with specific knowledge of https://github.com/VanAnderson/go-micro-example (or simply go-micro-example). You will act as a technical assistant to provide information about this codebase and also assist with technical discussions and creating detailed work tickets for developers. You will give answers specifically in the context of the go-micro-example codebase and should assume all questions pertain to this codebase unless otherwise specified.\" }
```
Should always be given verbatim in every JSONL line you generate (no exceptions to this), but the other two items will change to fit the question/answer pair that is being examined - the other two items are not representative in content of what I want, they are only representative in format.
*Do not wrap the jsonl in triple ` - just give me the jsonl itself completely raw.*
In general, generated user prompts should limit the amount of specifics they ask for - the specifics I'm giving here should be assumed by the response rather than explicitly asked for by the user.

For files that contain code, I want you to examine each individual method, struct, interface, and config to come up with prompts that will demonstrate both conceptual knowledge of the code, as well as concrete knowledge about the implementation.
Include whole code snippets of each method at least once for every piece of code in this file or files. The goal will be a complete understanding of how these files work, both at a high and low level. Every line of code should be accounted for. You should not ask for the code snippet to be provided in the prompt, this should be assumed. 
Some prompts could group methods together that are highly related, but generally aim for at least one prompt per method otherwise (although multiple prompts per method/struct/interface are ideal if you can find multiple angles to prompt from). Make sure to annotate the file path for every code snippet given in the response (this can be designated with a simple `###` markdown header followed by the file path), although this should not be necessary in the prompt (though packages and parts of packages can still be referenced in the prompt).

Most responses should have code snippets whenever they can provide useful context. Forego code snippets only when dealing with purely conceptual topics. Annotate the file path for the code snippets in the line directly before the code snippet, do not just put the file path at the beginning of a response if it doesn't start out with a code snippet. So every time you give a code snippet it should be:

[descriptive text about whatever]

### path/to/file.go
```golang
[codesnippet]
```

Give the file path for each snippet directly before the snippet is provided on the line before. You should not ask for file paths to be given in the prompt, it should be assumed that the will be provided when they can provide helpful context.

**Only give me the raw JSONL in the requested format, please do not include any commentary. If explicitly instructed to, you may return a blank response. **

# here are the files I want you to analyze:



### mocks.go
```
package invrepo

import (
	"context"

	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/db"
	"github.com/sksmith/go-micro-example/testutil"
)

type MockRepo struct {
	GetProductionEventByRequestIDFunc func(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error)
	SaveProductionEventFunc           func(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error

	GetReservationFunc            func(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error)
	GetReservationsFunc           func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error)
	GetReservationByRequestIDFunc func(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error)
	UpdateReservationFunc         func(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error
	SaveReservationFunc           func(ctx context.Context, reservation *inventory.Reservation, options ...core.UpdateOptions) error

	GetProductFunc  func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error)
	SaveProductFunc func(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error

	GetProductInventoryFunc    func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error)
	GetAllProductInventoryFunc func(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error)
	SaveProductInventoryFunc   func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error

	BeginTransactionFunc func(ctx context.Context) (core.Transaction, error)

	*testutil.CallWatcher
}

func (r *MockRepo) SaveProductionEvent(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error {
	r.AddCall(ctx, event, options)
	return r.SaveProductionEventFunc(ctx, event, options...)
}

func (r *MockRepo) UpdateReservation(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error {
	r.AddCall(ctx, ID, state, options)
	return r.UpdateReservationFunc(ctx, ID, state, qty, options...)
}

func (r *MockRepo) GetProductionEventByRequestID(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error) {
	r.AddCall(ctx, requestID, options)
	return r.GetProductionEventByRequestIDFunc(ctx, requestID, options...)
}

func (r *MockRepo) SaveReservation(ctx context.Context, reservation *inventory.Reservation, options ...core.UpdateOptions) error {
	r.AddCall(ctx, reservation, options)
	return r.SaveReservationFunc(ctx, reservation, options...)
}

func (r *MockRepo) GetReservation(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error) {
	r.AddCall(ctx, ID, options)
	return r.GetReservationFunc(ctx, ID, options...)
}

func (r *MockRepo) GetReservations(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
	r.AddCall(ctx, resOptions, limit, offset, options)
	return r.GetReservationsFunc(ctx, resOptions, limit, offset, options...)
}

func (r *MockRepo) SaveProduct(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error {
	r.AddCall(ctx, product, options)
	return r.SaveProductFunc(ctx, product, options...)
}

func (r *MockRepo) GetProduct(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
	r.AddCall(ctx, sku, options)
	return r.GetProductFunc(ctx, sku, options...)
}

func (r *MockRepo) GetProductInventory(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error) {
	r.AddCall(ctx, sku, options)
	return r.GetProductInventoryFunc(ctx, sku, options...)
}

func (r *MockRepo) SaveProductInventory(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
	r.AddCall(ctx, productInventory, options)
	return r.SaveProductInventoryFunc(ctx, productInventory, options...)
}

func (r *MockRepo) GetAllProductInventory(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error) {
	r.AddCall(ctx, limit, offset, options)
	return r.GetAllProductInventoryFunc(ctx, limit, offset, options...)
}

func (r *MockRepo) BeginTransaction(ctx context.Context) (core.Transaction, error) {
	r.AddCall(ctx)
	return r.BeginTransactionFunc(ctx)
}

func (r *MockRepo) GetReservationByRequestID(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error) {
	r.AddCall(ctx, requestId, options)
	return r.GetReservationByRequestIDFunc(ctx, requestId, options...)
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		SaveProductionEventFunc: func(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error {
			return nil
		},
		GetProductionEventByRequestIDFunc: func(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error) {
			return inventory.ProductionEvent{}, nil
		},
		SaveReservationFunc: func(ctx context.Context, reservation *inventory.Reservation, options ...core.UpdateOptions) error {
			return nil
		},
		GetReservationFunc: func(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error) {
			return inventory.Reservation{}, nil
		},
		GetReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
			return nil, nil
		},
		SaveProductFunc: func(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error { return nil },
		GetProductFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
			return inventory.Product{}, nil
		},
		GetAllProductInventoryFunc: func(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error) {
			return nil, nil
		},
		BeginTransactionFunc: func(ctx context.Context) (core.Transaction, error) { return db.NewMockTransaction(), nil },
		GetReservationByRequestIDFunc: func(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error) {
			return inventory.Reservation{}, nil
		},
		UpdateReservationFunc: func(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error {
			return nil
		},
		GetProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error) {
			return inventory.ProductInventory{}, nil
		},
		SaveProductInventoryFunc: func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
			return nil
		},
		CallWatcher: testutil.NewCallWatcher(),
	}
}

```
