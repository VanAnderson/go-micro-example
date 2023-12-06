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
