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
