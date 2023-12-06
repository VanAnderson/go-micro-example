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



### reservationapi.go
```
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
)

type ReservationService interface {
	Reserve(ctx context.Context, rr inventory.ReservationRequest) (inventory.Reservation, error)

	GetReservations(ctx context.Context, options inventory.GetReservationsOptions, limit, offset int) ([]inventory.Reservation, error)
	GetReservation(ctx context.Context, ID uint64) (inventory.Reservation, error)

	SubscribeReservations(ch chan<- inventory.Reservation) (id inventory.ReservationsSubID)
	UnsubscribeReservations(id inventory.ReservationsSubID)
}

type ReservationApi struct {
	service ReservationService
}

func NewReservationApi(service ReservationService) *ReservationApi {
	return &ReservationApi{service: service}
}

const (
	CtxKeyReservation CtxKey = "reservation"
)

func (ra *ReservationApi) ConfigureRouter(r chi.Router) {
	r.HandleFunc("/subscribe", ra.Subscribe)

	r.Route("/", func(r chi.Router) {
		r.With(Paginate).Get("/", ra.List)
		r.Put("/", ra.Create)

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(ra.ReservationCtx)
			r.Get("/", ra.Get)
			r.Delete("/", ra.Cancel)
		})
	})
}

func (a *ReservationApi) Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("client requesting subscription")

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Err(err).Msg("failed to establish inventory subscription connection")
		Render(w, r, ErrInternalServer)
	}
	go func() {
		defer conn.Close()

		ch := make(chan inventory.Reservation, 1)

		id := a.service.SubscribeReservations(ch)
		defer func() {
			a.service.UnsubscribeReservations(id)
		}()

		for res := range ch {
			resp := &ReservationResponse{Reservation: res}

			body, err := json.Marshal(resp)
			if err != nil {
				log.Err(err).Interface("clientId", id).Msg("failed to marshal product response")
				continue
			}

			log.Debug().Interface("clientId", id).Interface("reservationResponse", resp).Msg("sending reservation update to client")
			err = wsutil.WriteServerText(conn, body)
			if err != nil {
				log.Err(err).Interface("clientId", id).Msg("failed to write server message, disconnecting client")
				return
			}
		}
	}()
}

func (a *ReservationApi) Get(w http.ResponseWriter, r *http.Request) {
	res := r.Context().Value(CtxKeyReservation).(inventory.Reservation)

	resp := &ReservationResponse{Reservation: res}
	render.Status(r, http.StatusOK)
	Render(w, r, resp)
}

func (a *ReservationApi) Create(w http.ResponseWriter, r *http.Request) {
	data := &ReservationRequest{}
	if err := render.Bind(r, data); err != nil {
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	res, err := a.service.Reserve(r.Context(), *data.ReservationRequest)

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			Render(w, r, ErrNotFound)
		} else {
			log.Error().Err(err).Interface("reservationRequest", data).Msg("failed to reserve")
			Render(w, r, ErrInternalServer)
		}
		return
	}

	resp := &ReservationResponse{Reservation: res}
	render.Status(r, http.StatusCreated)
	Render(w, r, resp)
}

func (a *ReservationApi) Cancel(_ http.ResponseWriter, _ *http.Request) {
	// TODO Not implemented
}

func (a *ReservationApi) ReservationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		IDStr := chi.URLParam(r, "ID")
		if IDStr == "" {
			Render(w, r, ErrInvalidRequest(errors.New("reservation id is required")))
			return
		}

		ID, err := strconv.ParseUint(IDStr, 10, 64)
		if err != nil {
			log.Error().Err(err).Str("ID", IDStr).Msg("invalid reservation id")
			Render(w, r, ErrInvalidRequest(errors.New("invalid reservation id")))
		}

		reservation, err := a.service.GetReservation(r.Context(), ID)

		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				Render(w, r, ErrNotFound)
			} else {
				log.Error().Err(err).Str("id", IDStr).Msg("error acquiring product")
				Render(w, r, ErrInternalServer)
			}
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyReservation, reservation)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *ReservationApi) List(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(CtxKeyLimit).(int)
	offset := r.Context().Value(CtxKeyOffset).(int)

	sku := r.URL.Query().Get("sku")

	state, err := inventory.ParseReserveState(r.URL.Query().Get("state"))
	if err != nil {
		Render(w, r, ErrInvalidRequest(errors.New("invalid state")))
		return
	}

	res, err := a.service.GetReservations(r.Context(), inventory.GetReservationsOptions{Sku: sku, State: state}, limit, offset)

	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			Render(w, r, ErrNotFound)
		} else {
			log.Err(err).Send()
			Render(w, r, ErrInternalServer)
		}
		return
	}

	resList := NewReservationListResponse(res)
	render.Status(r, http.StatusOK)
	RenderList(w, r, resList)
}

```
