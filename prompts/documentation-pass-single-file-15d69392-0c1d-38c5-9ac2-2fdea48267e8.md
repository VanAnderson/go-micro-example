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

Here are three templates for interface, struct, and function documentation. If and only if this file defines a struct, interface, methods on a struct/interface, or individual function, I want you to generate extractive questions and answers where a user asks for documentation of a particular entity in the context of a package, then reply with the proper documentation with one of these templates.
You should only use the function documentation if the function is not a method on a struct or interface.
**IF THE FILE OR FILES DO NOT CONTAIN A STRUCT OR INTERFACE OR INDIVIDUAL FUNCTION, YOU MAY OUTPUT A BLANK RESPONSE WITH NOTHING IN IT**

```
## Struct Documentation

### `TypeName`

**Description:**

`TypeName` represents a ... It is used to ...

**Fields:**

- `FieldName`: The ... It is used for ...

**Constructors:**

#### `NewTypeName(param1 ParamType1, param2 ParamType2) *TypeName`

Creates a new instance of `TypeName`. It initializes the struct with the given parameters.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- A pointer to a newly allocated `TypeName` with fields initialized.

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Performs a specific operation on `TypeName`.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.

## Interface Documentation

### `InterfaceName`

**Description:**

`InterfaceName` defines the behavior for ... It is implemented by types that need to ...

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Performs an action. It is used to ...

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.

### Implementations

#### `ImplementingTypeName`

**Description:**

`ImplementingTypeName` is a struct that implements `InterfaceName`. It has the following properties and behaviors...

**Methods:**

#### `MethodName(param1 ParamType1, param2 ParamType2) (returnType, error)`

Implements the `InterfaceName.MethodName` for `ImplementingTypeName`.

**Parameters:**

- `param1` (ParamType1): The first parameter used for ...
- `param2` (ParamType2): The second parameter used for ...

**Returns:**

- `returnType`: Description of the return value.
- `error`: Any error that occurred during the operation, or `nil` if no error occurred.
## Function Documentation

### `FunctionName`

**Description:**

`FunctionName` performs a specific operation or calculation. It is used to ...

**Parameters:**

- `param1` (ParamType1): Description of the first parameter and its purpose.
- `param2` (ParamType2): Description of the second parameter and its purpose.
- ... (additional parameters as needed)

**Returns:**

- `returnType`: Description of the return value and what it represents.
- `error`: Description of the error returned, if applicable, or `nil` if no error occurred.

**Example:**

```go
result, err := FunctionName(param1, param2)
if err != nil {
    // Handle error
}
// Use result
```

**Notes:**

- Additional notes or considerations regarding the function's behavior, side effects, or usage in certain contexts.
```

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
