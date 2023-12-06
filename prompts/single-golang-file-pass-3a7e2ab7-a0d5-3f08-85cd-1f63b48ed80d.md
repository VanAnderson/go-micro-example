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
