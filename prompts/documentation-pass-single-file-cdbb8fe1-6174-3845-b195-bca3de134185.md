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



### inventoryapi.go
```
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
)

type InventoryService interface {
	Produce(ctx context.Context, product inventory.Product, event inventory.ProductionRequest) error
	CreateProduct(ctx context.Context, product inventory.Product) error

	GetProduct(ctx context.Context, sku string) (inventory.Product, error)
	GetAllProductInventory(ctx context.Context, limit, offset int) ([]inventory.ProductInventory, error)
	GetProductInventory(ctx context.Context, sku string) (inventory.ProductInventory, error)

	SubscribeInventory(ch chan<- inventory.ProductInventory) (id inventory.InventorySubID)
	UnsubscribeInventory(id inventory.InventorySubID)
}

type InventoryApi struct {
	service InventoryService
}

func NewInventoryApi(service InventoryService) *InventoryApi {
	return &InventoryApi{service: service}
}

const (
	CtxKeyProduct CtxKey = "product"
)

func (a *InventoryApi) ConfigureRouter(r chi.Router) {
	r.HandleFunc("/subscribe", a.Subscribe)

	r.Route("/", func(r chi.Router) {
		r.With(Paginate).Get("/", a.List)
		r.Put("/", a.CreateProduct)

		r.Route("/{sku}", func(r chi.Router) {
			r.Use(a.ProductCtx)
			r.Put("/productionEvent", a.CreateProductionEvent)
			r.Get("/", a.GetProductInventory)
		})
	})
}

// Subscribe provides consumes real-time inventory updates and sends them
// to the client via websocket connection.
//
// Note: This isn't exactly realistic because in the real world, this application
// would need to be able to scale. If it were scaled, clients would only get updates
// that occurred in their connected instance.
func (a *InventoryApi) Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("client requesting subscription")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Err(err).Msg("failed to establish inventory subscription connection")
		Render(w, r, ErrInternalServer)
	}
	go func() {
		defer conn.Close()

		ch := make(chan inventory.ProductInventory, 1)

		id := a.service.SubscribeInventory(ch)
		defer func() {
			a.service.UnsubscribeInventory(id)
		}()

		for inv := range ch {
			resp := &ProductResponse{ProductInventory: inv}
			body, err := json.Marshal(resp)
			if err != nil {
				log.Err(err).Interface("clientId", id).Msg("failed to marshal product response")
				continue
			}

			log.Debug().Interface("clientId", id).Interface("productResponse", resp).Msg("sending inventory update to client")
			err = wsutil.WriteServerText(conn, body)
			if err != nil {
				log.Err(err).Interface("clientId", id).Msg("failed to write server message, disconnecting client")
				return
			}
		}
	}()
}

func (a *InventoryApi) List(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(CtxKeyLimit).(int)
	offset := r.Context().Value(CtxKeyOffset).(int)

	products, err := a.service.GetAllProductInventory(r.Context(), limit, offset)
	if err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}

	RenderList(w, r, NewProductListResponse(products))
}

func (a *InventoryApi) CreateProduct(w http.ResponseWriter, r *http.Request) {
	data := &CreateProductRequest{}
	if err := render.Bind(r, data); err != nil {
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.service.CreateProduct(r.Context(), data.Product); err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	Render(w, r, NewProductResponse(inventory.ProductInventory{Product: data.Product}))
}

func (a *InventoryApi) ProductCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var product inventory.Product
		var err error

		sku := chi.URLParam(r, "sku")
		if sku == "" {
			Render(w, r, ErrInvalidRequest(errors.New("sku is required")))
			return
		}

		product, err = a.service.GetProduct(r.Context(), sku)

		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				Render(w, r, ErrNotFound)
			} else {
				log.Error().Err(err).Str("sku", sku).Msg("error acquiring product")
				Render(w, r, ErrInternalServer)
			}
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyProduct, product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *InventoryApi) CreateProductionEvent(w http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(CtxKeyProduct).(inventory.Product)

	data := &CreateProductionEventRequest{}
	if err := render.Bind(r, data); err != nil {
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.service.Produce(r.Context(), product, *data.ProductionRequest); err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	Render(w, r, &ProductionEventResponse{})
}

func (a *InventoryApi) GetProductInventory(w http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(CtxKeyProduct).(inventory.Product)

	res, err := a.service.GetProductInventory(r.Context(), product.Sku)

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			Render(w, r, ErrNotFound)
		} else {
			log.Err(err).Send()
			Render(w, r, ErrInternalServer)
		}
		return
	}

	resp := &ProductResponse{ProductInventory: res}
	render.Status(r, http.StatusOK)
	Render(w, r, resp)
}

```
