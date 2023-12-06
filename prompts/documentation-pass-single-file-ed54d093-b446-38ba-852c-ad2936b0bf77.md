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



### db.go
```
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/config"
	"github.com/sksmith/go-micro-example/core"
)

type dbconfig struct {
	timeZone              string
	sslMode               string
	poolMaxConns          int32
	poolMinConns          int32
	poolMaxConnLifetime   time.Duration
	poolMaxConnIdleTime   time.Duration
	poolHealthCheckPeriod time.Duration
}

type configOption func(cn *dbconfig)

func MinPoolConns(minConns int32) func(cn *dbconfig) {
	return func(c *dbconfig) {
		c.poolMinConns = minConns
	}
}

func MaxPoolConns(maxConns int32) func(cn *dbconfig) {
	return func(c *dbconfig) {
		c.poolMaxConns = maxConns
	}
}

func newDbConfig() dbconfig {
	return dbconfig{
		sslMode:               "disable",
		timeZone:              "UTC",
		poolMaxConns:          4,
		poolMinConns:          0,
		poolMaxConnLifetime:   time.Hour,
		poolMaxConnIdleTime:   time.Minute * 30,
		poolHealthCheckPeriod: time.Minute,
	}
}

func formatOption(url, option string, value interface{}) string {
	return url + " " + option + "=" + fmt.Sprintf("%v", value)
}

func addOptionsToConnStr(connStr string, options ...configOption) string {
	config := newDbConfig()
	for _, option := range options {
		option(&config)
	}

	connStr = formatOption(connStr, "sslmode", config.sslMode)
	connStr = formatOption(connStr, "TimeZone", config.timeZone)
	connStr = formatOption(connStr, "pool_max_conns", config.poolMaxConns)
	connStr = formatOption(connStr, "pool_min_conns", config.poolMinConns)
	connStr = formatOption(connStr, "pool_max_conn_lifetime", config.poolMaxConnLifetime)
	connStr = formatOption(connStr, "pool_max_conn_idle_time", config.poolMaxConnIdleTime)
	connStr = formatOption(connStr, "pool_health_check_period", config.poolHealthCheckPeriod)

	return connStr
}

func ConnectDb(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {

	log.Info().Str("host", cfg.Db.Host.Value).Str("name", cfg.Db.Name.Value).Msg("connecting to the database...")
	var err error

	if cfg.Db.Migrate.Value {
		log.Info().Msg("executing migrations")

		if err = RunMigrations(cfg); err != nil {
			log.Warn().Err(err).Msg("error executing migrations")
		}
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Db.Host.Value, cfg.Db.Port.Value, cfg.Db.User.Value, cfg.Db.Pass.Value, cfg.Db.Name.Value)

	var pool *pgxpool.Pool

	url := addOptionsToConnStr(connStr, MinPoolConns(int32(cfg.Db.Pool.MinSize.Value)), MaxPoolConns(int32(cfg.Db.Pool.MaxSize.Value)))
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	level, err := pgx.LogLevelFromString(cfg.Db.LogLevel.Value)
	if err != nil {
		return nil, err
	}
	poolConfig.ConnConfig.Logger = logger{level: level}

	for {
		pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err != nil {
			log.Error().Err(err).Msg("failed to create connection pool... retrying")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return pool, nil
}

type logger struct {
	level pgx.LogLevel
}

func (l logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if l.level < level {
		return
	}
	var evt *zerolog.Event
	switch level {
	case pgx.LogLevelTrace:
		evt = log.Trace()
	case pgx.LogLevelDebug:
		evt = log.Debug()
	case pgx.LogLevelInfo:
		evt = log.Info()
	case pgx.LogLevelWarn:
		evt = log.Warn()
	case pgx.LogLevelError:
		evt = log.Error()
	case pgx.LogLevelNone:
		evt = log.Info()
	default:
		evt = log.Info()
	}

	for k, v := range data {
		evt.Interface(k, v)
	}

	evt.Msg(msg)
}

func RunMigrations(cfg *config.Config) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User.Value,
		cfg.Db.Pass.Value,
		cfg.Db.Host.Value,
		cfg.Db.Port.Value,
		cfg.Db.Name.Value)

	m, err := migrate.New("file:"+cfg.Db.MigrationFolder.Value, connStr)
	if err != nil {
		return err
	}
	if cfg.Db.Clean.Value {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		log.Info().Msg("schema is up to date")
	}

	return nil
}

func GetQueryOptions(cn core.Conn, options ...core.QueryOptions) (conn core.Conn, forUpdate string) {
	conn = cn
	forUpdate = ""
	if len(options) > 0 {
		conn = options[0].Tx

		if options[0].ForUpdate {
			forUpdate = "FOR UPDATE"
		}
	}

	return conn, forUpdate
}

func GetUpdateOptions(cn core.Conn, options ...core.UpdateOptions) (conn core.Conn) {
	conn = cn
	if len(options) > 0 {
		conn = options[0].Tx
	}

	return conn
}

```
