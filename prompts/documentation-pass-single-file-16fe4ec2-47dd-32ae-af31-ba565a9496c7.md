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



### inventoryapi_test.go
```
package api_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gobwas/ws"
	"github.com/sksmith/go-micro-example/api"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/testutil"

	"github.com/go-chi/chi"
)

func TestInventorySubscribe(t *testing.T) {
	mockSvc := inventory.NewMockInventoryService()

	subscribeCalled := false
	expectedSubId := inventory.InventorySubID("subid1")
	unsubscribeCalled := false

	mockSvc.SubscribeInventoryFunc = func(ch chan<- inventory.ProductInventory) (id inventory.InventorySubID) {
		subscribeCalled = true
		go func() {
			inv := getTestProductInventory()
			for i := 0; i < 3; i++ {
				ch <- inv[i]
			}
			close(ch)
		}()

		return expectedSubId
	}

	mockSvc.UnsubscribeInventoryFunc = func(id inventory.InventorySubID) {
		unsubscribeCalled = true
	}

	invApi := api.NewInventoryApi(mockSvc)
	r := chi.NewRouter()
	invApi.ConfigureRouter(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := strings.Replace(ts.URL, "http", "ws", 1) + "/subscribe"

	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}

	curInv := getTestProductInventory()
	for i := 0; i < 3; i++ {
		got := &inventory.ProductInventory{}
		testutil.ReadWs(conn, got, t)

		if got.Name != curInv[i].Name {
			t.Errorf("unexpected ws response[%d] got=[%s] want=[%s]", i, got.Name, curInv[i].Name)
		}
	}

	if !subscribeCalled {
		t.Errorf("subscribe never called")
	}

	if !unsubscribeCalled {
		t.Errorf("unsubscribe never called")
	}
}

func setupInventoryTestServer() (*httptest.Server, *inventory.MockInventoryService) {
	mockSvc := inventory.NewMockInventoryService()
	invApi := api.NewInventoryApi(mockSvc)
	r := chi.NewRouter()
	invApi.ConfigureRouter(r)
	ts := httptest.NewServer(r)

	return ts, mockSvc
}

func TestInventoryList(t *testing.T) {
	ts, mockInvSvc := setupInventoryTestServer()
	defer ts.Close()

	tests := []struct {
		limit          int
		wantLimit      int
		offset         int
		wantOffset     int
		inventory      []inventory.ProductInventory
		serviceErr     error
		wantInventory  []inventory.ProductInventory
		wantErr        *api.ErrResponse
		wantStatusCode int
	}{
		{
			limit:          -1,
			wantLimit:      50,
			offset:         -1,
			wantOffset:     0,
			inventory:      getTestProductInventory(),
			wantInventory:  getTestProductInventory(),
			serviceErr:     nil,
			wantErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			limit:          5,
			wantLimit:      5,
			offset:         7,
			wantOffset:     7,
			inventory:      getTestProductInventory(),
			wantInventory:  getTestProductInventory(),
			serviceErr:     nil,
			wantErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			limit:          -1,
			wantLimit:      50,
			offset:         -1,
			wantOffset:     0,
			inventory:      []inventory.ProductInventory{},
			wantInventory:  []inventory.ProductInventory{},
			serviceErr:     nil,
			wantErr:        nil,
			wantStatusCode: http.StatusOK,
		},
		{
			limit:          -1,
			wantLimit:      50,
			offset:         -1,
			wantOffset:     0,
			inventory:      []inventory.ProductInventory{},
			wantInventory:  []inventory.ProductInventory{},
			serviceErr:     errors.New("something bad happened"),
			wantErr:        api.ErrInternalServer,
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		gotLimit := -1
		gotOffset := -1
		mockInvSvc.GetAllProductInventoryFunc = func(ctx context.Context, limit int, offset int) ([]inventory.ProductInventory, error) {
			gotLimit = limit
			gotOffset = offset
			return test.inventory, test.serviceErr
		}

		url := ts.URL
		if test.limit > -1 {
			url += fmt.Sprintf("?limit=%d&offset=%d", test.limit, test.offset)
		}

		res, err := http.Get(url)
		if err != nil {
			t.Fatal(err)
		}

		if test.wantErr == nil {
			got := []inventory.ProductInventory{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, test.wantInventory) {
				t.Errorf("inventory\n got:%+v\nwant:%+v\n", got, test.wantInventory)
			}
		} else {
			got := api.ErrResponse{}
			testutil.Unmarshal(res, &got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("errorResponse\n got:%v\nwant:%v\n", got.StatusText, test.wantErr.StatusText)
			}
		}

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=[%d] want=[%d]", res.StatusCode, test.wantStatusCode)
		}

		if gotLimit != test.wantLimit {
			t.Errorf("limit got=[%d] want=[%d]", gotLimit, test.limit)
		}

		if gotOffset != test.wantOffset {
			t.Errorf("offset got=[%d] want=[%d]", gotOffset, test.offset)
		}
	}
}

func TestInventoryCreateProduct(t *testing.T) {
	ts, mockInvSvc := setupInventoryTestServer()
	defer ts.Close()

	tests := []struct {
		request             api.CreateProductRequest
		serviceErr          error
		wantProductResponse *api.ProductResponse
		wantErr             *api.ErrResponse
		wantStatusCode      int
	}{
		{
			request:             createProductRequest("name1", "sku1", "upc1"),
			serviceErr:          nil,
			wantProductResponse: createProductResponse("name1", "sku1", "upc1", 0),
			wantErr:             nil,
			wantStatusCode:      http.StatusCreated,
		},
		{
			request:             createProductRequest("name1", "sku1", "upc1"),
			serviceErr:          errors.New("some unexpected error"),
			wantProductResponse: nil,
			wantErr:             api.ErrInternalServer,
			wantStatusCode:      http.StatusInternalServerError,
		},
		{
			request:             createProductRequest("name1", "sku1", ""),
			serviceErr:          nil,
			wantProductResponse: nil,
			wantErr:             api.ErrInvalidRequest(errors.New("missing required field(s)")),
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			request:             createProductRequest("name1", "", "upc1"),
			serviceErr:          nil,
			wantProductResponse: nil,
			wantErr:             api.ErrInvalidRequest(errors.New("missing required field(s)")),
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			request:             createProductRequest("", "sku1", "upc1"),
			serviceErr:          nil,
			wantProductResponse: nil,
			wantErr:             api.ErrInvalidRequest(errors.New("missing required field(s)")),
			wantStatusCode:      http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		mockInvSvc.CreateProductFunc = func(ctx context.Context, product inventory.Product) error {
			return test.serviceErr
		}

		res := testutil.Put(ts.URL, test.request, t)

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d\nwant=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantErr == nil {
			got := api.ProductResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, *test.wantProductResponse) {
				t.Errorf("product\n got=%+v\nwant=%+v", got, *test.wantProductResponse)
			}
		} else {
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, test.wantErr.StatusText)
			}
			if got.ErrorText != test.wantErr.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, test.wantErr.ErrorText)
			}
		}
	}
}

func TestInventoryCreateProductionEvent(t *testing.T) {
	ts, mockInvSvc := setupInventoryTestServer()
	defer ts.Close()

	tests := []struct {
		getProductFunc              func(ctx context.Context, sku string) (inventory.Product, error)
		produceFunc                 func(ctx context.Context, product inventory.Product, event inventory.ProductionRequest) error
		sku                         string
		request                     *api.CreateProductionEventRequest
		wantProductionEventResponse *api.ProductionEventResponse
		wantErr                     *api.ErrResponse
		wantStatusCode              int
	}{
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return getTestProductInventory()[0].Product, nil
			},
			produceFunc: func(ctx context.Context, product inventory.Product, event inventory.ProductionRequest) error {
				return nil
			},
			sku:                         "testsku1",
			request:                     createProductionEventRequest("abc123", 1),
			wantProductionEventResponse: &api.ProductionEventResponse{},
			wantErr:                     nil,
			wantStatusCode:              http.StatusCreated,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return inventory.Product{}, core.ErrNotFound
			},
			produceFunc:                 nil,
			sku:                         "testsku1",
			request:                     createProductionEventRequest("abc123", 1),
			wantProductionEventResponse: nil,
			wantErr:                     api.ErrNotFound,
			wantStatusCode:              http.StatusNotFound,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return inventory.Product{}, errors.New("some unexpected error")
			},
			produceFunc:                 nil,
			sku:                         "testsku1",
			request:                     createProductionEventRequest("abc123", 1),
			wantProductionEventResponse: nil,
			wantErr:                     api.ErrInternalServer,
			wantStatusCode:              http.StatusInternalServerError,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return getTestProductInventory()[0].Product, nil
			},
			produceFunc: func(ctx context.Context, product inventory.Product, event inventory.ProductionRequest) error {
				return errors.New("some unexpected error")
			},
			sku:                         "testsku1",
			request:                     createProductionEventRequest("abc123", 1),
			wantProductionEventResponse: nil,
			wantErr:                     api.ErrInternalServer,
			wantStatusCode:              http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		mockInvSvc.GetProductFunc = test.getProductFunc
		mockInvSvc.ProduceFunc = test.produceFunc

		url := ts.URL + "/" + test.sku + "/productionEvent"
		res := testutil.Put(url, test.request, t)

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d want=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantErr == nil {
			got := api.ProductionEventResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, *test.wantProductionEventResponse) {
				t.Errorf("product\n got=%+v\nwant=%+v", got, *test.wantProductionEventResponse)
			}
		} else {
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, test.wantErr.StatusText)
			}
			if got.ErrorText != test.wantErr.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, test.wantErr.ErrorText)
			}
		}
	}
}

func TestInventoryGetProductInventory(t *testing.T) {
	ts, mockInvSvc := setupInventoryTestServer()
	defer ts.Close()

	tests := []struct {
		sku                     string
		getProductFunc          func(ctx context.Context, sku string) (inventory.Product, error)
		getProductInventoryFunc func(ctx context.Context, sku string) (inventory.ProductInventory, error)
		wantProductResponse     *api.ProductResponse
		wantErr                 *api.ErrResponse
		wantStatusCode          int
	}{
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return getTestProductInventory()[0].Product, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string) (inventory.ProductInventory, error) {
				return getTestProductInventory()[0], nil
			},
			sku:                 "test1sku",
			wantProductResponse: createProductResponse("test1name", "test1sku", "test1upc", 1),
			wantErr:             nil,
			wantStatusCode:      http.StatusOK,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return inventory.Product{}, core.ErrNotFound
			},
			getProductInventoryFunc: nil,
			sku:                     "test1sku",
			wantProductResponse:     nil,
			wantErr:                 api.ErrNotFound,
			wantStatusCode:          http.StatusNotFound,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return getTestProductInventory()[0].Product, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string) (inventory.ProductInventory, error) {
				return inventory.ProductInventory{}, core.ErrNotFound
			},
			sku:                 "test1sku",
			wantProductResponse: nil,
			wantErr:             api.ErrNotFound,
			wantStatusCode:      http.StatusNotFound,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return inventory.Product{}, errors.New("some unexpected error")
			},
			getProductInventoryFunc: nil,
			sku:                     "test1sku",
			wantProductResponse:     nil,
			wantErr:                 api.ErrInternalServer,
			wantStatusCode:          http.StatusInternalServerError,
		},
		{
			getProductFunc: func(ctx context.Context, sku string) (inventory.Product, error) {
				return getTestProductInventory()[0].Product, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string) (inventory.ProductInventory, error) {
				return inventory.ProductInventory{}, errors.New("some unexpected error")
			},
			sku:                 "test1sku",
			wantProductResponse: nil,
			wantErr:             api.ErrInternalServer,
			wantStatusCode:      http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		mockInvSvc.GetProductFunc = test.getProductFunc
		mockInvSvc.GetProductInventoryFunc = test.getProductInventoryFunc

		res, err := http.Get(ts.URL + "/" + test.sku)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("status code got=%d want=%d", res.StatusCode, test.wantStatusCode)
		}

		if test.wantErr == nil {
			got := api.ProductResponse{}
			testutil.Unmarshal(res, &got, t)

			if !reflect.DeepEqual(got, *test.wantProductResponse) {
				t.Errorf("product\n got=%+v\nwant=%+v", got, *test.wantProductResponse)
			}
		} else {
			got := &api.ErrResponse{}
			testutil.Unmarshal(res, got, t)

			if got.StatusText != test.wantErr.StatusText {
				t.Errorf("status text got=%s want=%s", got.StatusText, test.wantErr.StatusText)
			}
			if got.ErrorText != test.wantErr.ErrorText {
				t.Errorf("error text got=%s want=%s", got.ErrorText, test.wantErr.ErrorText)
			}
		}
	}
}

func createProductionEventRequest(requestID string, quantity int64) *api.CreateProductionEventRequest {
	return &api.CreateProductionEventRequest{
		ProductionRequest: &inventory.ProductionRequest{RequestID: requestID, Quantity: quantity},
	}
}

func createProductRequest(name, sku, upc string) api.CreateProductRequest {
	return api.CreateProductRequest{Product: inventory.Product{Name: name, Sku: sku, Upc: upc}}
}

func createProductResponse(name, sku, upc string, available int64) *api.ProductResponse {
	return &api.ProductResponse{
		ProductInventory: inventory.ProductInventory{
			Available: available,
			Product:   inventory.Product{Name: name, Sku: sku, Upc: upc},
		},
	}
}

func getTestProductInventory() []inventory.ProductInventory {
	return []inventory.ProductInventory{
		{Available: 1, Product: inventory.Product{Sku: "test1sku", Upc: "test1upc", Name: "test1name"}},
		{Available: 2, Product: inventory.Product{Sku: "test2sku", Upc: "test2upc", Name: "test2name"}},
		{Available: 3, Product: inventory.Product{Sku: "test3sku", Upc: "test3upc", Name: "test3name"}},
	}
}

```
