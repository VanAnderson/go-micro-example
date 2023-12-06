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
