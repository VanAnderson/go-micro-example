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



### main.go
```
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
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

	"github.com/common-nighthawk/go-figure"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	cfg := config.Load("config")

	configLogging(cfg)
	printLogHeader(cfg)
	cfg.Print()

	dbPool := configDatabase(ctx, cfg)

	iq := queue.NewInventoryQueue(ctx, cfg)

	ir := invrepo.NewPostgresRepo(dbPool)

	invService := inventory.NewService(ir, iq)

	ur := usrrepo.NewPostgresRepo(dbPool)

	userService := user.NewService(ur)

	r := api.ConfigureRouter(cfg, invService, invService, userService)

	_ = queue.NewProductQueue(ctx, cfg, invService)

	log.Info().Str("port", cfg.Port.Value).Int64("startTimeMs", time.Since(start).Milliseconds()).Msg("listening")
	log.Fatal().Err(http.ListenAndServe(":"+cfg.Port.Value, r))
}

func printLogHeader(cfg *config.Config) {
	if cfg.Log.Structured.Value {
		log.Info().Str("application", cfg.AppName.Value).
			Str("revision", cfg.Revision.Value).
			Str("version", cfg.AppVersion.Value).
			Str("sha1ver", cfg.Sha1Version.Value).
			Str("build-time", cfg.BuildTime.Value).
			Str("profile", cfg.Profile.Value).
			Str("config-source", cfg.Config.Source.Value).
			Str("config-branch", cfg.Config.Spring.Branch.Value).
			Send()
	} else {
		f := figure.NewFigure(cfg.AppName.Value, "", true)
		f.Print()

		log.Info().Msg("=============================================")
		log.Info().Msg(fmt.Sprintf("      Revision: %s", cfg.Revision.Value))
		log.Info().Msg(fmt.Sprintf("       Profile: %s", cfg.Profile.Value))
		log.Info().Msg(fmt.Sprintf(" Config Server: %s - %s", cfg.Config.Source.Value, cfg.Config.Spring.Branch.Value))
		log.Info().Msg(fmt.Sprintf("   Tag Version: %s", cfg.AppVersion.Value))
		log.Info().Msg(fmt.Sprintf("  Sha1 Version: %s", cfg.Sha1Version.Value))
		log.Info().Msg(fmt.Sprintf("    Build Time: %s", cfg.BuildTime.Value))
		log.Info().Msg("=============================================")
	}
}

func configDatabase(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	dbPool, err := db.ConnectDb(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	return dbPool
}

func configLogging(cfg *config.Config) {
	log.Info().Msg("configuring logging...")

	if !cfg.Log.Structured.Value {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	level, err := zerolog.ParseLevel(cfg.Log.Level.Value)
	if err != nil {
		log.Warn().Str("loglevel", cfg.Log.Level.Value).Err(err).Msg("defaulting to info")
		level = zerolog.InfoLevel
	}
	log.Info().Str("loglevel", level.String()).Msg("setting log level")
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

```
