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
package inventory

import (
	"context"

	"github.com/sksmith/go-micro-example/testutil"
)

type MockInventoryService struct {
	ProduceFunc                func(ctx context.Context, product Product, event ProductionRequest) error
	CreateProductFunc          func(ctx context.Context, product Product) error
	GetProductFunc             func(ctx context.Context, sku string) (Product, error)
	GetAllProductInventoryFunc func(ctx context.Context, limit, offset int) ([]ProductInventory, error)
	GetProductInventoryFunc    func(ctx context.Context, sku string) (ProductInventory, error)
	SubscribeInventoryFunc     func(ch chan<- ProductInventory) (id InventorySubID)
	UnsubscribeInventoryFunc   func(id InventorySubID)
	*testutil.CallWatcher
}

func NewMockInventoryService() *MockInventoryService {
	return &MockInventoryService{
		ProduceFunc:       func(ctx context.Context, product Product, event ProductionRequest) error { return nil },
		CreateProductFunc: func(ctx context.Context, product Product) error { return nil },
		GetProductFunc:    func(ctx context.Context, sku string) (Product, error) { return Product{}, nil },
		GetAllProductInventoryFunc: func(ctx context.Context, limit, offset int) ([]ProductInventory, error) {
			return []ProductInventory{}, nil
		},
		GetProductInventoryFunc:  func(ctx context.Context, sku string) (ProductInventory, error) { return ProductInventory{}, nil },
		SubscribeInventoryFunc:   func(ch chan<- ProductInventory) (id InventorySubID) { return "" },
		UnsubscribeInventoryFunc: func(id InventorySubID) {},
		CallWatcher:              testutil.NewCallWatcher(),
	}
}

func (i *MockInventoryService) Produce(ctx context.Context, product Product, event ProductionRequest) error {
	i.AddCall(ctx, product, event)
	return i.ProduceFunc(ctx, product, event)
}

func (i *MockInventoryService) CreateProduct(ctx context.Context, product Product) error {
	i.AddCall(ctx, product)
	return i.CreateProductFunc(ctx, product)
}

func (i *MockInventoryService) GetProduct(ctx context.Context, sku string) (Product, error) {
	i.AddCall(ctx, sku)
	return i.GetProductFunc(ctx, sku)
}

func (i *MockInventoryService) GetAllProductInventory(ctx context.Context, limit, offset int) ([]ProductInventory, error) {
	i.AddCall(ctx, limit, offset)
	return i.GetAllProductInventoryFunc(ctx, limit, offset)
}

func (i *MockInventoryService) GetProductInventory(ctx context.Context, sku string) (ProductInventory, error) {
	i.AddCall(ctx, sku)
	return i.GetProductInventoryFunc(ctx, sku)
}

func (i *MockInventoryService) SubscribeInventory(ch chan<- ProductInventory) (id InventorySubID) {
	i.AddCall(ch)
	return i.SubscribeInventoryFunc(ch)
}

func (i *MockInventoryService) UnsubscribeInventory(id InventorySubID) {
	i.AddCall(id)
	i.UnsubscribeInventoryFunc(id)
}

type MockReservationService struct {
	ReserveFunc func(ctx context.Context, rr ReservationRequest) (Reservation, error)

	GetReservationsFunc func(ctx context.Context, options GetReservationsOptions, limit, offset int) ([]Reservation, error)
	GetReservationFunc  func(ctx context.Context, ID uint64) (Reservation, error)

	SubscribeReservationsFunc   func(ch chan<- Reservation) (id ReservationsSubID)
	UnsubscribeReservationsFunc func(id ReservationsSubID)
	*testutil.CallWatcher
}

func NewMockReservationService() *MockReservationService {
	return &MockReservationService{
		ReserveFunc: func(ctx context.Context, rr ReservationRequest) (Reservation, error) { return Reservation{}, nil },
		GetReservationsFunc: func(ctx context.Context, options GetReservationsOptions, limit, offset int) ([]Reservation, error) {
			return []Reservation{}, nil
		},
		GetReservationFunc:          func(ctx context.Context, ID uint64) (Reservation, error) { return Reservation{}, nil },
		SubscribeReservationsFunc:   func(ch chan<- Reservation) (id ReservationsSubID) { return "" },
		UnsubscribeReservationsFunc: func(id ReservationsSubID) {},
		CallWatcher:                 testutil.NewCallWatcher(),
	}
}

func (r *MockReservationService) Reserve(ctx context.Context, rr ReservationRequest) (Reservation, error) {
	r.CallWatcher.AddCall(ctx, rr)
	return r.ReserveFunc(ctx, rr)
}

func (r *MockReservationService) GetReservations(ctx context.Context, options GetReservationsOptions, limit, offset int) ([]Reservation, error) {
	r.CallWatcher.AddCall(ctx, options, limit, offset)
	return r.GetReservationsFunc(ctx, options, limit, offset)
}

func (r *MockReservationService) GetReservation(ctx context.Context, ID uint64) (Reservation, error) {
	r.CallWatcher.AddCall(ctx, ID)
	return r.GetReservationFunc(ctx, ID)
}

func (r *MockReservationService) SubscribeReservations(ch chan<- Reservation) (id ReservationsSubID) {
	r.CallWatcher.AddCall(ch)
	return r.SubscribeReservationsFunc(ch)
}

func (r *MockReservationService) UnsubscribeReservations(id ReservationsSubID) {
	r.CallWatcher.AddCall(id)
	r.UnsubscribeReservationsFunc(id)
}

```
