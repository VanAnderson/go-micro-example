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



### service.go
```
package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
)

func NewService(repo Repository, q InventoryQueue) *service {
	log.Info().Msg("creating inventory service...")
	return &service{
		repo:            repo,
		queue:           q,
		inventorySubs:   make(map[InventorySubID]chan<- ProductInventory),
		reservationSubs: make(map[ReservationsSubID]chan<- Reservation),
	}
}

type InventorySubID string
type ReservationsSubID string

type GetReservationsOptions struct {
	Sku   string
	State ReserveState
}

type service struct {
	repo            Repository
	queue           InventoryQueue
	inventorySubs   map[InventorySubID]chan<- ProductInventory
	reservationSubs map[ReservationsSubID]chan<- Reservation
}

func (s *service) CreateProduct(ctx context.Context, product Product) error {
	const funcName = "CreateProduct"

	dbProduct, err := s.repo.GetProduct(ctx, product.Sku)
	if err != nil != errors.Is(err, core.ErrNotFound) {
		return errors.WithStack(err)
	}
	if err == nil {
		log.Debug().Str("func", funcName).Str("sku", dbProduct.Sku).Msg("product already exists")
		return nil
	}

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			rollback(ctx, tx, err)
		}
	}()

	log.Debug().Str("func", funcName).Str("sku", product.Sku).Msg("creating product")
	if err = s.repo.SaveProduct(ctx, product, core.UpdateOptions{Tx: tx}); err != nil {
		return errors.WithStack(err)
	}

	log.Debug().Str("func", funcName).Str("sku", product.Sku).Msg("creating product inventory")
	pi := ProductInventory{Product: product}

	if err = s.repo.SaveProductInventory(ctx, pi, core.UpdateOptions{Tx: tx}); err != nil {
		return errors.WithStack(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) Produce(ctx context.Context, product Product, pr ProductionRequest) error {
	const funcName = "Produce"

	log.Debug().
		Str("func", funcName).
		Str("sku", product.Sku).
		Str("requestId", pr.RequestID).
		Int64("quantity", pr.Quantity).
		Msg("producing inventory")

	if pr.RequestID == "" {
		return errors.New("request id is required")
	}
	if pr.Quantity < 1 {
		return errors.New("quantity must be greater than zero")
	}

	event, err := s.repo.GetProductionEventByRequestID(ctx, pr.RequestID)
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return errors.WithStack(err)
	}

	if event.RequestID != "" {
		log.Debug().Str("func", funcName).Str("requestId", pr.RequestID).Msg("production request already exists")
		return nil
	}

	event = ProductionEvent{
		RequestID: pr.RequestID,
		Sku:       product.Sku,
		Quantity:  pr.Quantity,
		Created:   time.Now(),
	}

	tx, err := s.repo.BeginTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			rollback(ctx, tx, err)
		}
	}()

	if err = s.repo.SaveProductionEvent(ctx, &event, core.UpdateOptions{Tx: tx}); err != nil {
		return errors.WithMessage(err, "failed to save production event")
	}

	productInventory, err := s.repo.GetProductInventory(ctx, product.Sku, core.QueryOptions{Tx: tx, ForUpdate: true})
	if err != nil {
		return errors.WithMessage(err, "failed to get product inventory")
	}

	productInventory.Available += event.Quantity
	if err = s.repo.SaveProductInventory(ctx, productInventory, core.UpdateOptions{Tx: tx}); err != nil {
		return errors.WithMessage(err, "failed to add production to product")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.WithMessage(err, "failed to commit production transaction")
	}

	err = s.publishInventory(ctx, productInventory)
	if err != nil {
		return errors.WithMessage(err, "failed to publish inventory")
	}

	if err = s.FillReserves(ctx, product); err != nil {
		return errors.WithMessage(err, "failed to fill reserves after production")
	}

	return nil
}

func (s *service) Reserve(ctx context.Context, rr ReservationRequest) (Reservation, error) {
	const funcName = "Reserve"

	log.Debug().
		Str("func", funcName).
		Str("requestID", rr.RequestID).
		Str("sku", rr.Sku).
		Str("requester", rr.Requester).
		Int64("quantity", rr.Quantity).
		Msg("reserving inventory")

	if err := validateReservationRequest(rr); err != nil {
		return Reservation{}, err
	}

	tx, err := s.repo.BeginTransaction(ctx)
	defer func() {
		if err != nil {
			rollback(ctx, tx, err)
		}
	}()
	if err != nil {
		return Reservation{}, err
	}

	pr, err := s.repo.GetProduct(ctx, rr.Sku, core.QueryOptions{Tx: tx, ForUpdate: true})
	if err != nil {
		log.Error().Err(err).Str("requestId", rr.RequestID).Msg("failed to get product")
		return Reservation{}, errors.WithStack(err)
	}

	res, err := s.repo.GetReservationByRequestID(ctx, rr.RequestID, core.QueryOptions{Tx: tx, ForUpdate: true})
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		log.Error().Err(err).Str("requestId", rr.RequestID).Msg("failed to get reservation request")
		return Reservation{}, errors.WithStack(err)
	}
	if res.RequestID != "" {
		log.Debug().Str("func", funcName).Str("requestId", rr.RequestID).Msg("reservation already exists, returning it")
		rollback(ctx, tx, err)
		return res, nil
	}

	res = Reservation{
		RequestID:         rr.RequestID,
		Requester:         rr.Requester,
		Sku:               rr.Sku,
		State:             Open,
		RequestedQuantity: rr.Quantity,
		Created:           time.Now(),
	}

	if err = s.repo.SaveReservation(ctx, &res, core.UpdateOptions{Tx: tx}); err != nil {
		return Reservation{}, errors.WithStack(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return Reservation{}, errors.WithStack(err)
	}

	if err = s.FillReserves(ctx, pr); err != nil {
		return Reservation{}, errors.WithStack(err)
	}

	return res, nil
}

func validateReservationRequest(rr ReservationRequest) error {
	if rr.RequestID == "" {
		return errors.New("request id is required")
	}
	if rr.Requester == "" {
		return errors.New("requester is required")
	}
	if rr.Sku == "" {
		return errors.New("sku is requred")
	}
	if rr.Quantity < 1 {
		return errors.New("quantity is required")
	}
	return nil
}

func (s *service) GetAllProductInventory(ctx context.Context, limit, offset int) ([]ProductInventory, error) {
	return s.repo.GetAllProductInventory(ctx, limit, offset)
}

func (s *service) GetProduct(ctx context.Context, sku string) (Product, error) {
	const funcName = "GetProduct"

	log.Debug().Str("func", funcName).Str("sku", sku).Msg("getting product")

	product, err := s.repo.GetProduct(ctx, sku)
	if err != nil {
		return product, errors.WithStack(err)
	}
	return product, nil
}

func (s *service) GetProductInventory(ctx context.Context, sku string) (ProductInventory, error) {
	const funcName = "GetProductInventory"

	log.Debug().Str("func", funcName).Str("sku", sku).Msg("getting product inventory")

	product, err := s.repo.GetProductInventory(ctx, sku)
	if err != nil {
		return product, errors.WithStack(err)
	}
	return product, nil
}

func (s *service) GetReservation(ctx context.Context, ID uint64) (Reservation, error) {
	const funcName = "GetReservation"

	log.Debug().Str("func", funcName).Uint64("id", ID).Msg("getting reservation")

	rsv, err := s.repo.GetReservation(ctx, ID)
	if err != nil {
		return rsv, errors.WithStack(err)
	}
	return rsv, nil
}

func (s *service) GetReservations(ctx context.Context, options GetReservationsOptions, limit, offset int) ([]Reservation, error) {
	const funcName = "GetProductInventory"

	log.Debug().
		Str("func", funcName).
		Str("sku", options.Sku).
		Str("state", string(options.State)).
		Msg("getting reservations")

	rsv, err := s.repo.GetReservations(ctx, options, limit, offset)
	if err != nil {
		return rsv, errors.WithStack(err)
	}
	return rsv, nil
}

func (s *service) SubscribeInventory(ch chan<- ProductInventory) (id InventorySubID) {
	id = InventorySubID(uuid.NewString())
	s.inventorySubs[id] = ch
	log.Debug().Interface("clientId", id).Msg("subscribing to inventory")
	return id
}

func (s *service) UnsubscribeInventory(id InventorySubID) {
	log.Debug().Interface("clientId", id).Msg("unsubscribing from inventory")
	close(s.inventorySubs[id])
	delete(s.inventorySubs, id)
}

func (s *service) SubscribeReservations(ch chan<- Reservation) (id ReservationsSubID) {
	id = ReservationsSubID(uuid.NewString())
	s.reservationSubs[id] = ch
	log.Debug().Interface("clientId", id).Msg("subscribing to reservations")
	return id
}

func (s *service) UnsubscribeReservations(id ReservationsSubID) {
	log.Debug().Interface("clientId", id).Msg("unsubscribing from reservations")
	close(s.reservationSubs[id])
	delete(s.reservationSubs, id)
}

func (s *service) FillReserves(ctx context.Context, product Product) error {
	const funcName = "fillReserves"

	tx, err := s.repo.BeginTransaction(ctx)
	defer func() {
		if err != nil {
			rollback(ctx, tx, err)
		}
	}()
	if err != nil {
		return errors.WithStack(err)
	}

	openReservations, err := s.repo.GetReservations(ctx, GetReservationsOptions{Sku: product.Sku, State: Open}, 100, 0, core.QueryOptions{Tx: tx, ForUpdate: true})
	if err != nil {
		return errors.WithStack(err)
	}

	productInventory, err := s.repo.GetProductInventory(ctx, product.Sku, core.QueryOptions{Tx: tx, ForUpdate: true})
	if err != nil {
		return errors.WithStack(err)
	}

	for _, reservation := range openReservations {
		var subtx pgx.Tx
		subtx, err = tx.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				rollback(ctx, subtx, err)
			}
		}()
		reservation := reservation

		log.Trace().
			Str("func", funcName).
			Str("sku", product.Sku).
			Str("reservation.RequestID", reservation.RequestID).
			Int64("productInventory.Available", productInventory.Available).
			Msg("fulfilling reservation")

		if productInventory.Available == 0 {
			break
		}

		reserveAmount := reservation.RequestedQuantity - reservation.ReservedQuantity
		if reserveAmount > productInventory.Available {
			reserveAmount = productInventory.Available
		}
		productInventory.Available -= reserveAmount
		reservation.ReservedQuantity += reserveAmount

		if reservation.ReservedQuantity == reservation.RequestedQuantity {
			reservation.State = Closed
		}

		log.Debug().
			Str("func", funcName).
			Str("sku", product.Sku).
			Str("reservation.RequestID", reservation.RequestID).
			Msg("saving product inventory")

		err = s.repo.SaveProductInventory(ctx, productInventory, core.UpdateOptions{Tx: tx})
		if err != nil {
			return errors.WithStack(err)
		}

		log.Debug().
			Str("func", funcName).
			Str("sku", product.Sku).
			Str("reservation.RequestID", reservation.RequestID).
			Str("state", string(reservation.State)).
			Msg("updating reservation")

		err = s.repo.UpdateReservation(ctx, reservation.ID, reservation.State, reservation.ReservedQuantity, core.UpdateOptions{Tx: tx})
		if err != nil {
			return errors.WithStack(err)
		}

		if err = subtx.Commit(ctx); err != nil {
			return errors.WithStack(err)
		}

		err = s.publishInventory(ctx, productInventory)
		if err != nil {
			return errors.WithStack(err)
		}

		err = s.publishReservation(ctx, reservation)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) publishInventory(ctx context.Context, pi ProductInventory) error {
	err := s.queue.PublishInventory(ctx, pi)
	if err != nil {
		return errors.WithMessage(err, "failed to publish inventory to queue")
	}
	go s.notifyInventorySubscribers(pi)
	return nil
}

func (s *service) publishReservation(ctx context.Context, r Reservation) error {
	err := s.queue.PublishReservation(ctx, r)
	if err != nil {
		return errors.WithMessage(err, "failed to publish reservation to queue")
	}
	go s.notifyReservationSubscribers(r)
	return nil
}

func (s *service) notifyInventorySubscribers(pi ProductInventory) {
	for id, ch := range s.inventorySubs {
		log.Debug().Interface("clientId", id).Interface("productInventory", pi).Msg("notifying subscriber of inventory update")
		ch <- pi
	}
}

func (s *service) notifyReservationSubscribers(r Reservation) {
	for id, ch := range s.reservationSubs {
		log.Debug().Interface("clientId", id).Interface("productInventory", r).Msg("notifying subscriber of reservation update")
		ch <- r
	}
}

```
