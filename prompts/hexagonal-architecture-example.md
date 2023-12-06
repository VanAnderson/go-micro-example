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

For these files, I want you to generate a large list of extractive questions and answers that focus on how the codebase in these files reflects the concepts of a hexagonal architecture.
You should be able to find many different angles to prompt from, and you should be able to generate many questions and answers from each angle.
In particular, the interace of the core package and the db package should be a focus of the prompts, as well as the implementation of the core package. 
The question and response pairs should crystalize aspects of hexagonal architecture and (importantly) how it's realized in this codebase.

Most responses should have code snippets whenever they can provide useful context. Forego code snippets only when dealing with purely conceptual topics. Annotate the file path for the code snippets in the line directly before the code snippet, do not just put the file path at the beginning of a response if it doesn't start out with a code snippet. So every time you give a code snippet it should be:

[descriptive text about whatever]

### path/to/file.go
```golang
[codesnippet]
```

Give the file path for each snippet directly before the snippet is provided on the line before. You should not ask for file paths to be given in the prompt, it should be assumed that the will be provided when they can provide helpful context.

**Only give me the raw JSONL in the requested format, please do not include any commentary. If explicitly instructed to, you may return a blank response. **

# here are the files I want you to analyze:



### repo.go
```
package core

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var ErrNotFound = errors.New("core: record not found")

type Conn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Transaction interface {
	Conn
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UpdateOptions struct {
	Tx Transaction
}

type QueryOptions struct {
	ForUpdate bool
	Tx        Transaction
}

```

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

### model.go
```
// Package inventory is a rudimentary model that represents a fictional inventory tracking system for a factory. A real
// factory would obviously need much more fine grained detail and would probably use a different ubiquitous language.
package inventory

import (
	"time"

	"github.com/pkg/errors"
)

// ProductionRequest is a value object. A request to produce inventory.
type ProductionRequest struct {
	RequestID string `json:"requestID"`
	Quantity  int64  `json:"quantity"`
}

// ProductionEvent is an entity. An addition to inventory through production of a Product.
type ProductionEvent struct {
	ID        uint64    `json:"id"`
	RequestID string    `json:"requestID"`
	Sku       string    `json:"sku"`
	Quantity  int64     `json:"quantity"`
	Created   time.Time `json:"created"`
}

// Product is a value object. A SKU able to be produced by the factory.
type Product struct {
	Sku  string `json:"sku"`
	Upc  string `json:"upc"`
	Name string `json:"name"`
}

// ProductInventory is an entity. It represents current inventory levels for the associated product.
type ProductInventory struct {
	Product
	Available int64 `json:"available"`
}

type ReserveState string

const (
	Open   ReserveState = "Open"
	Closed ReserveState = "Closed"
	None   ReserveState = ""
)

func ParseReserveState(v string) (ReserveState, error) {
	switch v {
	case string(Open):
		return Open, nil
	case string(Closed):
		return Closed, nil
	case string(None):
		return None, nil
	default:
		return None, errors.New("invalid reserve state")
	}
}

type ReservationRequest struct {
	Sku       string `json:"sku"`
	RequestID string `json:"requestId"`
	Requester string `json:"requester"`
	Quantity  int64  `json:"quantity"`
}

// Reservation is an entity. An amount of inventory set aside for a given Customer.
type Reservation struct {
	ID                uint64       `json:"id"`
	RequestID         string       `json:"requestId"`
	Requester         string       `json:"requester"`
	Sku               string       `json:"sku"`
	State             ReserveState `json:"state"`
	ReservedQuantity  int64        `json:"reservedQuantity"`
	RequestedQuantity int64        `json:"requestedQuantity"`
	Created           time.Time    `json:"created"`
}

```

### repository.go
```
package inventory

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
)

func rollback(ctx context.Context, tx core.Transaction, err error) {
	if tx == nil {
		return
	}
	e := tx.Rollback(ctx)
	if e != nil {
		log.Warn().Err(err).Msg("failed to rollback")
	}
}

type Transactional interface {
	BeginTransaction(ctx context.Context) (core.Transaction, error)
}

type Repository interface {
	ProductionEventRepository
	ReservationRepository
	InventoryRepository
	ProductRepository
}

type ProductionEventRepository interface {
	Transactional
	GetProductionEventByRequestID(ctx context.Context, requestID string, options ...core.QueryOptions) (pe ProductionEvent, err error)

	SaveProductionEvent(ctx context.Context, event *ProductionEvent, options ...core.UpdateOptions) error
}

type ReservationRepository interface {
	Transactional
	GetReservations(ctx context.Context, resOptions GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]Reservation, error)
	GetReservationByRequestID(ctx context.Context, requestId string, options ...core.QueryOptions) (Reservation, error)
	GetReservation(ctx context.Context, ID uint64, options ...core.QueryOptions) (Reservation, error)

	SaveReservation(ctx context.Context, reservation *Reservation, options ...core.UpdateOptions) error
	UpdateReservation(ctx context.Context, ID uint64, state ReserveState, qty int64, options ...core.UpdateOptions) error
}

type InventoryRepository interface {
	Transactional
	GetProductInventory(ctx context.Context, sku string, options ...core.QueryOptions) (pi ProductInventory, err error)
	GetAllProductInventory(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]ProductInventory, error)

	SaveProductInventory(ctx context.Context, productInventory ProductInventory, options ...core.UpdateOptions) error
}

type ProductRepository interface {
	Transactional
	GetProduct(ctx context.Context, sku string, options ...core.QueryOptions) (Product, error)

	SaveProduct(ctx context.Context, product Product, options ...core.UpdateOptions) error
}

type InventoryQueue interface {
	PublishInventory(ctx context.Context, productInventory ProductInventory) error
	PublishReservation(ctx context.Context, reservation Reservation) error
}

```

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

### service_test.go
```
package inventory_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/db"
	"github.com/sksmith/go-micro-example/db/invrepo"
	"github.com/sksmith/go-micro-example/queue"
	"github.com/sksmith/go-micro-example/testutil"
)

func TestMain(m *testing.M) {
	testutil.ConfigLogging()
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
			name:    "unexpected error beginning transaction",
			product: inventory.Product{Name: "productname", Sku: "productsku", Upc: "productupc"},

			beginTransactionFunc: func(ctx context.Context) (core.Transaction, error) { return nil, errors.New("some unexpected error") },

			wantRepoCallCnt: map[string]int{"SaveProduct": 0, "SaveProductInventory": 0},
			wantTxCallCnt:   map[string]int{"Commit": 0, "Rollback": 0},
			wantErr:         true,
		},
		{
			name:    "unexpected error comitting",
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

func TestGetReservation(t *testing.T) {
	reservations := getReservations()
	tests := []struct {
		name string
		ID   uint64

		getReservationFunc func(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error)

		wantReservation inventory.Reservation
		wantErr         bool
	}{
		{
			name:            "reservation is returned",
			wantReservation: reservations[0],
		},
		{
			name: "reservation is returned",
			getReservationFunc: func(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error) {
				return inventory.Reservation{}, errors.New("some unexpected error")
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getReservationFunc != nil {
			mockRepo.GetReservationFunc = test.getReservationFunc
		} else {
			mockRepo.GetReservationFunc = func(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error) {
				return reservations[0], nil
			}
		}
		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.GetReservation(context.Background(), 0)

			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(res, test.wantReservation) {
				t.Errorf("unexpected reservation got=%v want=%v", res, test.wantReservation)
			}
		})
	}
}

func TestGetReservations(t *testing.T) {
	reservations := getReservations()
	tests := []struct {
		name    string
		options inventory.GetReservationsOptions
		limit   int
		offset  int

		getReservationsFunc func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit int, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error)

		wantReservations []inventory.Reservation
		wantErr          bool
	}{
		{
			name:             "reservations are returned",
			wantReservations: reservations,
		},
		{
			name: "reservation is returned",
			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit int, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{}, errors.New("some unexpected error")
			},
			wantReservations: []inventory.Reservation{},
			wantErr:          true,
		},
	}

	for _, test := range tests {
		mockRepo := invrepo.NewMockRepo()
		if test.getReservationsFunc != nil {
			mockRepo.GetReservationsFunc = test.getReservationsFunc
		} else {
			mockRepo.GetReservationsFunc = func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit int, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return reservations, nil
			}
		}
		mockQueue := queue.NewMockQueue()

		service := inventory.NewService(mockRepo, mockQueue)

		t.Run(test.name, func(t *testing.T) {
			res, err := service.GetReservations(context.Background(), test.options, test.limit, test.offset)

			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(res, test.wantReservations) {
				t.Errorf("unexpected reservations got=%v want=%v", res, test.wantReservations)
			}
		})
	}
}

type reservationUpdate struct {
	ID       uint64
	State    inventory.ReserveState
	Quantity int64
}

func TestFillReserves(t *testing.T) {
	product := inventory.Product{Name: "name", Sku: "sku", Upc: "upc"}
	tests := []struct {
		name                    string
		product                 inventory.Product
		saveProductInventoryErr error
		updateReservationErr    error

		getReservationsFunc      func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit int, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error)
		getProductInventoryFunc  func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error)
		saveProductInventoryFunc func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error
		updateReservationFunc    func(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error

		publishInventoryFunc   func(ctx context.Context, pi inventory.ProductInventory) error
		publishReservationFunc func(ctx context.Context, r inventory.Reservation) error

		beginTransactionFunc func(ctx context.Context) (core.Transaction, error)
		commitFunc           func(ctx context.Context) error

		wantRepoCallCnt      map[string]int
		wantQueueCallCnt     map[string]int
		wantTxCallCnt        map[string]int
		wantSubTxCallCnt     map[string]int
		wantProductInventory inventory.ProductInventory
		wantResUpdates       []reservationUpdate
		wantErr              bool
	}{
		{
			name:    "enough inventory to close reservation",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 10},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},

			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 0,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Closed, Quantity: 10},
			},
			wantQueueCallCnt: map[string]int{"PublishInventory": 1, "PublishReservation": 1},
			wantSubTxCallCnt: map[string]int{"Commit": 1, "Rollback": 0},
			wantTxCallCnt:    map[string]int{"Commit": 1, "Rollback": 0},
		},
		{
			name:    "not enough inventory to close reservation",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 10},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 5}, nil
			},

			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 0,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Open, Quantity: 5},
			},

			wantQueueCallCnt: map[string]int{"PublishInventory": 1, "PublishReservation": 1},
			wantSubTxCallCnt: map[string]int{"Commit": 1, "Rollback": 0},
			wantTxCallCnt:    map[string]int{"Commit": 1, "Rollback": 0},
		},
		{
			name:    "enough inventory to close multiple reservations",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
					{ID: 1, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
					{ID: 2, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},

			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 1,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Closed, Quantity: 3},
				{ID: 1, State: inventory.Closed, Quantity: 3},
				{ID: 2, State: inventory.Closed, Quantity: 3},
			},

			wantQueueCallCnt: map[string]int{"PublishInventory": 3, "PublishReservation": 3},
			wantSubTxCallCnt: map[string]int{"Commit": 3, "Rollback": 0},
			wantTxCallCnt:    map[string]int{"Commit": 1, "Rollback": 0},
		},
		{
			name:    "unexpected error saving inventory",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},
			saveProductInventoryErr: errors.New("some unexpected error"),

			wantErr: true,
			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 7,
			},
			wantResUpdates:   []reservationUpdate{},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantSubTxCallCnt: map[string]int{"Commit": 0, "Rollback": 1},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
		},
		{
			name:    "unexpected error updating reservation",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},
			updateReservationErr: errors.New("some unexpected error"),

			wantErr: true,
			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 7,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Closed, Quantity: 3},
			},
			wantQueueCallCnt: map[string]int{"PublishInventory": 0, "PublishReservation": 0},
			wantSubTxCallCnt: map[string]int{"Commit": 0, "Rollback": 1},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
		},
		{
			name:    "unexpected error publishing inventory",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},
			publishInventoryFunc: func(ctx context.Context, pi inventory.ProductInventory) error {
				return errors.New("some unexpected error")
			},

			wantErr: true,
			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 7,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Closed, Quantity: 3},
			},
			wantQueueCallCnt: map[string]int{"PublishInventory": 1, "PublishReservation": 0},
			wantSubTxCallCnt: map[string]int{"Commit": 1, "Rollback": 1},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
		},
		{
			name:    "unexpected error publishing reservation",
			product: product,

			getReservationsFunc: func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
				return []inventory.Reservation{
					{ID: 0, State: inventory.Open, ReservedQuantity: 0, RequestedQuantity: 3},
				}, nil
			},
			getProductInventoryFunc: func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
				return inventory.ProductInventory{Product: product, Available: 10}, nil
			},
			publishReservationFunc: func(ctx context.Context, r inventory.Reservation) error {
				return errors.New("some unexpected error")
			},

			wantErr: true,
			wantProductInventory: inventory.ProductInventory{
				Product:   product,
				Available: 7,
			},
			wantResUpdates: []reservationUpdate{
				{ID: 0, State: inventory.Closed, Quantity: 3},
			},
			wantQueueCallCnt: map[string]int{"PublishInventory": 1, "PublishReservation": 1},
			wantSubTxCallCnt: map[string]int{"Commit": 1, "Rollback": 1},
			wantTxCallCnt:    map[string]int{"Commit": 0, "Rollback": 1},
		},
	}

	for _, test := range tests {
		if test.name == "unexpected error publishing reservation" {
			fmt.Println("ugh")
		}
		mockTx := db.NewMockTransaction()
		if test.commitFunc != nil {
			mockTx.CommitFunc = test.commitFunc
		}

		mockSubTx := db.NewMockPgxTx()
		mockTx.BeginFunc = func(ctx context.Context) (pgx.Tx, error) {
			return mockSubTx, nil
		}

		mockRepo := invrepo.NewMockRepo()
		if test.beginTransactionFunc != nil {
			mockRepo.BeginTransactionFunc = test.beginTransactionFunc
		} else {
			mockRepo.BeginTransactionFunc = func(ctx context.Context) (core.Transaction, error) {
				return mockTx, nil
			}
		}
		if test.getReservationsFunc != nil {
			mockRepo.GetReservationsFunc = test.getReservationsFunc
		}
		if test.getProductInventoryFunc != nil {
			mockRepo.GetProductInventoryFunc = test.getProductInventoryFunc
		}
		var gotProductInventory inventory.ProductInventory
		mockRepo.SaveProductInventoryFunc = func(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
			gotProductInventory = productInventory
			return test.saveProductInventoryErr
		}

		gotResUpdates := []reservationUpdate{}
		mockRepo.UpdateReservationFunc = func(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error {
			gotResUpdates = append(gotResUpdates, reservationUpdate{ID: ID, State: state, Quantity: qty})
			return test.updateReservationErr
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
			err := service.FillReserves(context.Background(), test.product)

			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(gotProductInventory, test.wantProductInventory) {
				t.Errorf("unexpected product inventory\n got=%+v\nwant=%+v", gotProductInventory, test.wantProductInventory)
			}

			if !reflect.DeepEqual(gotResUpdates, test.wantResUpdates) {
				t.Errorf("unexpected reservation updates\n got=%+v\nwant=%+v", gotResUpdates, test.wantResUpdates)
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
			for f, c := range test.wantSubTxCallCnt {
				mockSubTx.VerifyCount(f, c, t)
			}
		})
	}
}

func TestSubscribeInventory(t *testing.T) {
	mockRepo := invrepo.NewMockRepo()
	mockQueue := queue.NewMockQueue()
	service := inventory.NewService(mockRepo, mockQueue)

	mockRepo.GetProductInventoryFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error) {
		return getProductInventory()[2], nil
	}

	ch := make(chan inventory.ProductInventory)
	id := service.SubscribeInventory(ch)

	go func() {
		service.Produce(context.Background(), getProductInventory()[2].Product, inventory.ProductionRequest{RequestID: "request1", Quantity: 1})
	}()

	want := getProductInventory()[2]
	want.Available++

	select {
	case got := <-ch:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected product got=%v want=%v", got, want)
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("timed out waiting for product inventory from channel")
	}

	service.UnsubscribeInventory(id)

	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("channel should be closed")
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("channel should be closed by now")
	}
}

func TestSubscribeReservations(t *testing.T) {
	mockRepo := invrepo.NewMockRepo()
	mockQueue := queue.NewMockQueue()
	service := inventory.NewService(mockRepo, mockQueue)

	mockRepo.GetReservationsFunc = func(ctx context.Context, resOptions inventory.GetReservationsOptions, limit int, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
		return []inventory.Reservation{getReservations()[3]}, nil
	}
	mockRepo.GetProductInventoryFunc = func(ctx context.Context, sku string, options ...core.QueryOptions) (pi inventory.ProductInventory, err error) {
		pi = getProductInventory()[2]
		pi.Available += 10
		return pi, nil
	}

	ch := make(chan inventory.Reservation)
	id := service.SubscribeReservations(ch)

	go func() {
		service.Produce(context.Background(), getProductInventory()[2].Product, inventory.ProductionRequest{RequestID: "request1", Quantity: 10})
	}()

	want := getReservations()[3]
	want.State = inventory.Closed
	want.ReservedQuantity = want.RequestedQuantity

	select {
	case got := <-ch:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected reservation\n got=%+v\nwant=%+v", got, want)
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("timed out waiting for reservation from channel")
	}

	service.UnsubscribeReservations(id)

	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("channel should be closed")
		}
	case <-time.After(10 * time.Millisecond):
		t.Error("channel should be closed by now")
	}
}

func getProductInventory() []inventory.ProductInventory {
	return []inventory.ProductInventory{
		{Product: inventory.Product{Sku: "sku1", Upc: "upc1", Name: "name1"}, Available: 1},
		{Product: inventory.Product{Sku: "sku2", Upc: "upc2", Name: "name2"}, Available: 10},
		{Product: inventory.Product{Sku: "sku3", Upc: "upc3", Name: "name3"}, Available: 0},
	}
}

func getReservations() []inventory.Reservation {
	return []inventory.Reservation{
		{ID: 0, RequestID: "request1", Requester: "requester1", Sku: "sku1", State: inventory.Closed, ReservedQuantity: 10, RequestedQuantity: 10},
		{ID: 1, RequestID: "request2", Requester: "requester2", Sku: "sku1", State: inventory.Closed, ReservedQuantity: 3, RequestedQuantity: 3},
		{ID: 2, RequestID: "request3", Requester: "requester1", Sku: "sku2", State: inventory.Closed, ReservedQuantity: 10, RequestedQuantity: 10},
		{ID: 3, RequestID: "request4", Requester: "requester1", Sku: "sku3", State: inventory.Open, ReservedQuantity: 2, RequestedQuantity: 10},
	}
}

```

### mocks.go
```
package user

import "context"

type MockUserService struct {
	CreateFunc func(ctx context.Context, user CreateUserRequest) (User, error)
	GetFunc    func(ctx context.Context, username string) (User, error)
	DeleteFunc func(ctx context.Context, username string) error
	LoginFunc  func(ctx context.Context, username, password string) (User, error)
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		CreateFunc: func(ctx context.Context, user CreateUserRequest) (User, error) { return User{}, nil },
		GetFunc:    func(ctx context.Context, username string) (User, error) { return User{}, nil },
		DeleteFunc: func(ctx context.Context, username string) error { return nil },
		LoginFunc:  func(ctx context.Context, username, password string) (User, error) { return User{}, nil },
	}
}

func (u *MockUserService) Create(ctx context.Context, user CreateUserRequest) (User, error) {
	return u.CreateFunc(ctx, user)
}

func (u *MockUserService) Get(ctx context.Context, username string) (User, error) {
	return u.GetFunc(ctx, username)
}

func (u *MockUserService) Delete(ctx context.Context, username string) error {
	return u.DeleteFunc(ctx, username)
}

func (u *MockUserService) Login(ctx context.Context, username, password string) (User, error) {
	return u.LoginFunc(ctx, username, password)
}

```

### model.go
```
package user

import "time"

type CreateUserRequest struct {
	Username          string `json:"username,omitempty"`
	IsAdmin           bool   `json:"isAdmin,omitempty"`
	PlainTextPassword string `json:"-"`
}

type User struct {
	Username string
	HashedPassword string
	IsAdmin bool
	Created time.Time
}
```

### service.go
```
package user

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
	"golang.org/x/crypto/bcrypt"
)

func NewService(repo Repository) *service {
	log.Info().Msg("creating user service...")

	return &service{repo: repo}
}

type service struct {
	repo Repository
}

func (s *service) Get(ctx context.Context, username string) (User, error) {
	return s.repo.Get(ctx, username)
}

func (s *service) Create(ctx context.Context, req CreateUserRequest) (User, error) {
	if !usernameIsValid(req.Username) {
		return User{}, errors.New("invalid username")
	}
	if !passwordIsValid(req.PlainTextPassword) {
		return User{}, errors.New("invalid password")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.PlainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	user := &User{
		Username:       req.Username,
		HashedPassword: string(hash),
		Created:        time.Now(),
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

func usernameIsValid(username string) bool {
	return true
}

func passwordIsValid(password string) bool {
	return true
}

func (s *service) Delete(ctx context.Context, username string) error {
	return s.repo.Delete(ctx, username)
}

func (s *service) Login(ctx context.Context, username, password string) (User, error) {
	u, err := s.repo.Get(ctx, username)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return User{}, err
	}

	return u, nil
}

type Repository interface {
	Create(ctx context.Context, user *User, tx ...core.UpdateOptions) error
	Get(ctx context.Context, username string, tx ...core.QueryOptions) (User, error)
	Delete(ctx context.Context, username string, tx ...core.UpdateOptions) error
}

```

### service_test.go
```
package user_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/user"
	"github.com/sksmith/go-micro-example/db/usrrepo"
)

func TestGet(t *testing.T) {
	usr := user.User{Username: "someuser", HashedPassword: "somehashedpassword", IsAdmin: false, Created: time.Now()}
	tests := []struct {
		name     string
		username string

		getFunc func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error)

		wantUser user.User
		wantErr  bool
	}{
		{
			name:     "user is returned",
			username: "someuser",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantUser: usr,
		},
		{
			name:     "error is returned",
			username: "someuser",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},

			wantErr:  true,
			wantUser: user.User{},
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.getFunc != nil {
			mockRepo.GetFunc = test.getFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Get(context.Background(), test.username)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if !reflect.DeepEqual(got, test.wantUser) {
				t.Errorf("unexpected user\n got=%+v\nwant=%+v", got, test.wantUser)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		request user.CreateUserRequest

		createFunc func(ctx context.Context, user *user.User, tx ...core.UpdateOptions) error

		wantUsername    string
		wantRepoCallCnt map[string]int
		wantErr         bool
	}{
		{
			name:    "user is returned",
			request: user.CreateUserRequest{Username: "someuser", IsAdmin: false, PlainTextPassword: "plaintextpw"},

			wantRepoCallCnt: map[string]int{"Create": 1},
			wantUsername:    "someuser",
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.createFunc != nil {
			mockRepo.CreateFunc = test.createFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Create(context.Background(), test.request)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if got.Username != test.wantUsername {
				t.Errorf("unexpected username got=%+v want=%+v", got.Username, test.wantUsername)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		username string

		deleteFunc func(ctx context.Context, username string, tx ...core.UpdateOptions) error

		wantRepoCallCnt map[string]int
		wantErr         bool
	}{
		{
			name:            "user is deleted",
			username:        "someuser",
			wantRepoCallCnt: map[string]int{"Delete": 1},
		},
		{
			name:     "error is returned",
			username: "someuser",

			deleteFunc: func(ctx context.Context, username string, tx ...core.UpdateOptions) error {
				return errors.New("some unexpected error")
			},

			wantRepoCallCnt: map[string]int{"Delete": 1},
			wantErr:         true,
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.deleteFunc != nil {
			mockRepo.DeleteFunc = test.deleteFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			err := service.Delete(context.Background(), test.username)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			for f, c := range test.wantRepoCallCnt {
				mockRepo.VerifyCount(f, c, t)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	usr := user.User{Username: "someuser", HashedPassword: "$2a$10$t67eB.bOkZGovKD8wqqppO7q.SqWwTS8FUrUx3GAW57GMhkD2Zcwy", IsAdmin: false, Created: time.Now()}
	tests := []struct {
		name     string
		username string
		password string

		getFunc func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error)

		wantUsername string
		wantErr      bool
	}{
		{
			name:     "correct password",
			username: "someuser",
			password: "plaintextpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantUsername: "someuser",
		},
		{
			name:     "wrong password",
			username: "someuser",
			password: "wrongpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return usr, nil
			},

			wantErr:      true,
			wantUsername: "",
		},
		{
			name:     "unexpected error getting user",
			username: "someuser",
			password: "wrongpw",

			getFunc: func(ctx context.Context, username string, options ...core.QueryOptions) (user.User, error) {
				return user.User{}, errors.New("some unexpected error")
			},

			wantErr:      true,
			wantUsername: "",
		},
	}

	for _, test := range tests {
		mockRepo := usrrepo.NewMockRepo()
		if test.getFunc != nil {
			mockRepo.GetFunc = test.getFunc
		}

		service := user.NewService(mockRepo)

		t.Run(test.name, func(t *testing.T) {
			got, err := service.Login(context.Background(), test.username, test.password)
			if test.wantErr && err == nil {
				t.Errorf("expected error, got none")
			} else if !test.wantErr && err != nil {
				t.Errorf("did not want error, got=%v", err)
			}

			if got.Username != test.wantUsername {
				t.Errorf("unexpected username got=%v want=%v", got.Username, test.wantUsername)
			}
		})
	}
}

```

### README.md
```
# Go Micro Example

![Linter](https://github.com/sksmith/note-server/actions/workflows/lint.yml/badge.svg)
![Security](https://github.com/sksmith/note-server/actions/workflows/sec.yml/badge.svg)
![Test](https://github.com/sksmith/note-server/actions/workflows/test.yml/badge.svg)

This is an inventory management microservice for an online retailer. I structured the project using a hexagonal style abstracting
away business logic from dependencies like the RESTful API, the Postgres database, and RabbitMQ message queue.

## Structure

The Go community generally likes application directory structures to be as simple as possible which is
totally admirable and applicable for a small simple microservice. I could probably have kept everything
for this project in a single directory and focused on making sure it met twelve factors. But I'm a big
fan of [Domain Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html), and how it gels so
nicely with [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) and I wanted
to see how a Go microservice might look structured using them.

The starting point of the application is under the [cmd](cmd/main.go) directory. The "domain"
core of the application where all business logic should reside is under the [core](core)
directory. The other directories listed there are each of the external dependencies for the project.

![structure diagram](inventory.jpg)

## Running the Application Locally

This project requires that you have Docker, Go and Make installed locally. If you do, you can start
the application first by starting the docker-compose file, then start the application using the
supplied Makefile.

```shell
docker-compose -f ./scripts/docker-compose.yml up -d
make run
```

If you want to create a deployable executable and run it:

```shell
make build
./bin/inventory
```

### Run Docker Compose

```shell
docker-compose up
```

## Application Features

### RESTful API

This application uses the wonderful [go-chi](https://github.com/go-chi/chi) for routing
[beautifuly documentation](https://github.com/go-chi/chi/blob/master/_examples/rest/main.go) served as the main 
inspiration for how to structure the API. Seriously, I was so impressed.

In Java I like to generate the controller layer using Open API so that the contract and implementation always match 
exactly. I couldn't quite find an equivalent solution I liked.

Truth be told, if I were doing inter-microservice communication I would strongly consider using gRPC rather than a 
RESTful API.

### Authentication

Many of the endpoints in this project are protected by using a [simple authentication middleware](api/middleware.go). If 
you're interested in hitting them you can use basic auth admin:admin. Users are stored in the database along with their 
hashed password. Users are locally cached using [golang-lru](https://github.com/hashicorp/golang-lru). In a production 
setting if I actually wanted caching I'd either use a remote cache like Redis, or a distributed local cache like 
groupcache to prevent stale or out of sync data.

### Metrics

This application outputs prometheus metrics using middleware I plugged into the go-chi router. If you're running
locally check them out at [http://localhost:8080/metrics](http://localhost:8080/metrics). Every URL automatically
gets a hit count and a latency metric added. You can find the configurations [here](api/middleware.go).

### Logging

I ended up going with [zerolog](https://github.com/rs/zerolog) for logging in this project. I really like its API and 
their benchmarks look really great too! You can get structured logging or nice human-readable logging by
[changing some configs](config.yml#L10)

### Configuration

This project uses [viper](https://github.com/spf13/viper) for handling externalized configurations. At the moment it only reads from the local config.yml but the plan is to make it compatible with [Spring Cloud Config](https://cloud.spring.io/spring-cloud-config), and [etcd](https://etcd.io).

### Testing

I chose not to go with any of the test frameworks when putting this project together. I felt like using interfaces and 
injecting dependencies would be enough to allow me to mock what I need to. There's a fair bit of boilerplate code 
required to mock, say, the inventory repository but not having to pull in and learn yet another dependency for testing 
seemed like a fair tradeoff.

The testing in this project is pretty bare-bones and mostly just proof-of-concept. If you want to see some tests, 
though, they're in [api](api). I personally prefer more integration tests that test an application front-to-back for 
features rather than tons and tons of tightly-coupled unit tests.

### Database Migrations

I'm using the [migrate](https://github.com/golang-migrate/migrate) project to manage database migrations.

```shell
migrate create -ext sql -dir db/migrations -seq create_products_table

migrate -database postgres://postgres:postgres@localhost:5432/smfg-db?sslmode=disable -path db/migrations up

migrate -source file://db/migrations -database postgres://localhost:5432/database down
```

## 12 Factors

One of the goals of this service was to ensure all [12 principals](https://12factor.net/) of a 12-factor app are adhered 
to. This was a nice way to make sure the app I built offered most of what you need out of a Spring Boot application.

### I. Codebase

The application is stored in my git repository.

### II. Dependencies

Go handles this for us through its dependency management system (yay!)

### III. Config

See the [configuration section](#Configuration) section above.

### IV. Backing Services

The application connects to all external dependencies (in this case, RabbitMQ, and Postgres) via URLs which it gets from 
remote configuration.

### V. Build, release, run

The application can easily be plugged into any CI/CD pipeline. This is mostly thanks to Go making this easy through 
great command line tools.

### VI. Processes

This app is not *strictly* stateless. There is a cache in the user repository. This was a design choice I made in the 
interest of seeing what setting up a local cache in go might look like. In a more real-world application you would 
probably want an external cache (like Redis), or a distributed cache (like 
[Group Cache](https://github.com/golang/groupcache) - which is really cool!)

This app is otherwise stateless and threadsafe.

### VII. Port Binding

The application binds to a supplied port on startup.

### VIII. Concurrency

Other than maintaining an instance-based cache (see Process above), the application will scale horizontally without 
issue. The database dependency would need to scale vertically unless you started using sharding, or a distributed data 
store like [Cosmos DB](https://docs.microsoft.com/en-us/azure/cosmos-db/distribute-data-globally).

### IX. Disposability

One of the wonderful things about Go is how *fast* it starts up. This application can start up and shut down in a 
fraction of the time that similar Spring Boot microservices. In addition, they use a much smaller footprint. This is 
perfect for services that need to be highly elastic on demand.

### X. Dev/Prod Parity

Docker makes standing up a prod-like environment on your local environment a breeze. This application has
[a docker-compose file](scripts/docker-compose.yml) that starts up a local instance of rabbit and postgres. This 
obviously doesn't account for ensuring your dev and stage environments are up to snuff but at least that's a good start 
for local development.

### XI. Logs

Logs in the application are written to the stdout allowing for logscrapers like 
[logstash](https://www.elastic.co/logstash) to consume and parse the logs. Through configuration the logs can output as 
plain text for ease of reading during local development and then switched after deployment into json structured logs for 
automatic parsing.

### XII. Admin Processes

Database migration is automated in the project using [migrate](https://github.com/golang-migrate/migrate).

## TODO

- [ ] Recreate architecture diagram
- [ ] Add godoc
- [ ] Return 204 no content if data already exists
- [ ] Cleanup TODOs
```
