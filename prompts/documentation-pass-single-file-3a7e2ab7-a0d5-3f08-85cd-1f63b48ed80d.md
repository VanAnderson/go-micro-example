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



### integration_test.go
```
package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sksmith/go-micro-example/api"
	"github.com/sksmith/go-micro-example/config"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/core/user"
	"github.com/sksmith/go-micro-example/db"
	"github.com/sksmith/go-micro-example/db/invrepo"
	"github.com/sksmith/go-micro-example/db/usrrepo"
	"github.com/sksmith/go-micro-example/queue"
	"github.com/sksmith/go-micro-example/testutil"
)

var cfg *config.Config

func TestMain(m *testing.M) {

	log.Info().Msg("configuring logging...")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx := context.Background()
	cfg = config.Load("config_test")

	level, err := zerolog.ParseLevel(cfg.Log.Level.Value)
	if err != nil {
		log.Fatal().Err(err)
	}
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cfg.Print()

	dbPool, err := db.ConnectDb(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	iq := queue.NewInventoryQueue(ctx, cfg)

	ir := invrepo.NewPostgresRepo(dbPool)

	invService := inventory.NewService(ir, iq)

	ur := usrrepo.NewPostgresRepo(dbPool)

	userService := user.NewService(ur)

	r := api.ConfigureRouter(cfg, invService, invService, userService)

	_ = queue.NewProductQueue(ctx, cfg, invService)

	go func() {
		log.Fatal().Err(http.ListenAndServe(":"+cfg.Port.Value, r))
	}()

	waitForReady()
	os.Exit(m.Run())
}

func waitForReady() {
	for {
		res, err := http.Get(host() + "/health")
		if err == nil && res.StatusCode == 200 {
			break
		}
		log.Info().Msg("application not ready, sleeping")
		time.Sleep(1 * time.Second)
	}
}

func TestCreateProduct(t *testing.T) {
	cases := []struct {
		name    string
		product inventory.Product

		wantSku        string
		wantStatusCode int
	}{
		{
			name:           "valid request",
			product:        inventory.Product{Sku: "somesku", Upc: "someupc", Name: "somename"},
			wantSku:        "somesku",
			wantStatusCode: 201,
		},
		{
			name:           "valid request with a long name",
			product:        inventory.Product{Sku: "someskuwithareallylongname", Upc: "longskuupc", Name: "somename"},
			wantSku:        "someskuwithareallylongname",
			wantStatusCode: 201,
		},
		{
			name:           "missing sku",
			product:        inventory.Product{Sku: "", Upc: "skurequiredupc", Name: "skurequiredname"},
			wantStatusCode: 400,
		},
		{
			name:           "missing upc",
			product:        inventory.Product{Sku: "upcreqsku", Upc: "", Name: "upcreqname"},
			wantStatusCode: 400,
		},
		{
			name:           "missing name",
			product:        inventory.Product{Sku: "namereqsku", Upc: "namerequpc", Name: ""},
			wantStatusCode: 400,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			request := api.CreateProductRequest{Product: test.product}
			res := testutil.Put(host()+"/api/v1/inventory", request, t)

			if res.StatusCode != test.wantStatusCode {
				t.Errorf("unexpected status got=%d want=%d", res.StatusCode, test.wantStatusCode)
			}

			body := &api.ProductResponse{}
			testutil.Unmarshal(res, body, t)
			if test.wantSku != "" && body.Sku != test.wantSku {
				t.Errorf("unexpected response sku got=%s want=%s", body.Sku, test.wantSku)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name string
		url  string

		wantMinRespLen int
		wantStatusCode int
	}{
		{
			name:           "valid request",
			url:            "/api/v1/inventory",
			wantMinRespLen: 2,
			wantStatusCode: 200,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			res, err := http.Get(host() + test.url)

			if err != nil {
				t.Errorf("unexpected error got=%s", err)
			}
			if res.StatusCode != test.wantStatusCode {
				t.Errorf("unexpected status got=%d want=%d", res.StatusCode, test.wantStatusCode)
			}

			body := []inventory.ProductInventory{}
			testutil.Unmarshal(res, &body, t)
			if len(body) < test.wantMinRespLen {
				t.Errorf("unexpected response len got=%d want=%d", len(body), test.wantMinRespLen)
			}
		})
	}
}

func host() string {
	return "http://localhost:" + cfg.Port.Value
}

```
