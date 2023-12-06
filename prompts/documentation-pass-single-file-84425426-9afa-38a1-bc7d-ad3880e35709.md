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



### queue.go
```
package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/config"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/streadway/amqp"
)

type InventoryQueue struct {
	cfg         *config.Config
	inventory   chan<- message
	reservation chan<- message
}

func NewInventoryQueue(ctx context.Context, cfg *config.Config) *InventoryQueue {
	invChan := make(chan message)
	resChan := make(chan message)

	iq := &InventoryQueue{
		cfg:         cfg,
		inventory:   invChan,
		reservation: resChan,
	}

	url := getUrl(cfg)

	go func() {
		invExch := cfg.RabbitMQ.Inventory.Exchange.Value
		publish(redial(ctx, url), invExch, invChan)
		ctx.Done()
	}()

	go func() {
		resExch := cfg.RabbitMQ.Reservation.Exchange.Value
		publish(redial(ctx, url), resExch, resChan)
		ctx.Done()
	}()

	return iq
}

func getUrl(cfg *config.Config) string {
	rmq := cfg.RabbitMQ
	return fmt.Sprintf("amqp://%s:%s@%s:%s", rmq.User.Value, rmq.Pass.Value, rmq.Host.Value, rmq.Port.Value)
}

func (i *InventoryQueue) PublishInventory(ctx context.Context, productInventory inventory.ProductInventory) error {
	body, err := json.Marshal(productInventory)
	if err != nil {
		return errors.WithMessage(err, "failed to serialize message for queue")
	}
	i.inventory <- message(body)
	return nil
}

func (i *InventoryQueue) PublishReservation(ctx context.Context, reservation inventory.Reservation) error {
	body, err := json.Marshal(reservation)
	if err != nil {
		return errors.WithMessage(err, "error marshalling reservation to send to queue")
	}
	i.reservation <- message(body)
	return nil
}

type ProductQueue struct {
	cfg        *config.Config
	product    <-chan message
	productDlt chan<- message
}

func NewProductQueue(ctx context.Context, cfg *config.Config, handler ProductHandler) *ProductQueue {
	log.Info().Msg("creating product queue...")

	prodChan := make(chan message)
	prodDltChan := make(chan message)

	pq := &ProductQueue{
		cfg:        cfg,
		product:    prodChan,
		productDlt: prodDltChan,
	}

	url := getUrl(cfg)

	go func() {
		for message := range prodChan {
			product := inventory.Product{}
			err := json.Unmarshal(message, &product)
			if err != nil {
				log.Error().Err(err).Msg("error unmarshalling product, writing to dlt")
				pq.sendToDlt(ctx, message)
			}

			err = handler.CreateProduct(ctx, product)
			if err != nil {
				log.Error().Err(err).Msg("failed to create product, sending to dlt")
				pq.productDlt <- message
			}
		}
	}()

	go func() {
		prodQueue := cfg.RabbitMQ.Product.Queue.Value
		subscribe(redial(ctx, url), prodQueue, prodChan)
		ctx.Done()
	}()

	go func() {
		dltExch := cfg.RabbitMQ.Product.Dlt.Exchange.Value
		publish(redial(ctx, url), dltExch, prodDltChan)
		ctx.Done()
	}()

	return pq
}

type ProductHandler interface {
	CreateProduct(ctx context.Context, product inventory.Product) error
}

func (p *ProductQueue) sendToDlt(ctx context.Context, body []byte) {
	p.productDlt <- message(body)
}

// TODO We should be using one exchange per domain object here.
// exchange binds the publishers to the subscribers
// const exchange = "pubsub"

// message is the application type for a message.  This can contain identity,
// or a reference to the recevier chan for further demuxing.
type message []byte

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// redial continually connects to the URL, exiting the program when no longer possible
func redial(ctx context.Context, url string) chan chan session {
	sessions := make(chan chan session)

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Fatal().Msg("shutting down session factory")
				return
			}

			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatal().Err(err).Str("url", url).Msg("cannot (re)dial")
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatal().Err(err).Msg("cannot create channel")
			}

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Info().Msg("shutting down new session")
				return
			}
		}
	}()

	return sessions
}

// publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func publish(sessions chan chan session, exchange string, messages <-chan message) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan message, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Info().Msg("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Debug().Str("exchange", exchange).Msg("ready to publish messages")

	Publish:
		for {
			var body message
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Info().Uint64("message", confirmed.DeliveryTag).Str("body", string(body)).Msg("nack")
				}
				reading = messages

			case body = <-pending:
				routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
				err := pub.Publish(exchange, routingKey, false, false, amqp.Publishing{
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					_ = pub.Close()
					break Publish
				}

			case body, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

// subscribe consumes deliveries from an exclusive queue from a fanout exchange and sends to the application specific messages chan.
func subscribe(sessions chan chan session, queue string, messages chan<- message) {
	for session := range sessions {
		sub := <-session

		deliveries, err := sub.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			log.Error().Str("queue", queue).Err(err).Msg("cannot consume from")
			return
		}

		log.Info().Str("queue", queue).Msg("listening for messages")

		for msg := range deliveries {
			messages <- message(msg.Body)
			err = sub.Ack(msg.DeliveryTag, false)
			if err != nil {
				log.Error().Err(err).Str("queue", queue).Msg("failed to acknowledge to queue")
			}
		}
	}
}

```
