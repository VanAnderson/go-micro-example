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



### mocks.go
```
package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sksmith/go-micro-example/testutil"
)

type MockConn struct {
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	ExecFunc     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	BeginFunc    func(ctx context.Context) (pgx.Tx, error)
	*testutil.CallWatcher
}

func NewMockConn() MockConn {
	return MockConn{
		QueryFunc:    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) { return nil, nil },
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row { return nil },
		ExecFunc:     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) { return nil, nil },
		BeginFunc:    func(ctx context.Context) (pgx.Tx, error) { return NewMockPgxTx(), nil },
		CallWatcher:  testutil.NewCallWatcher(),
	}
}

func (c *MockConn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	c.AddCall(ctx, sql, args)
	return c.QueryFunc(ctx, sql, args)
}

func (c *MockConn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	c.AddCall(ctx, sql, args)
	return c.QueryRowFunc(ctx, sql, args)
}

func (c *MockConn) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	c.AddCall(ctx, sql, args)
	return c.ExecFunc(ctx, sql, args)
}

func (c *MockConn) Begin(ctx context.Context) (pgx.Tx, error) {
	c.AddCall(ctx)
	return c.BeginFunc(ctx)
}

type MockTransaction struct {
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error

	MockConn
	*testutil.CallWatcher
}

func NewMockTransaction() *MockTransaction {
	return &MockTransaction{
		MockConn:     NewMockConn(),
		CommitFunc:   func(ctx context.Context) error { return nil },
		RollbackFunc: func(ctx context.Context) error { return nil },
		CallWatcher:  testutil.NewCallWatcher(),
	}
}

func (t *MockTransaction) Commit(ctx context.Context) error {
	t.AddCall(ctx)
	return t.CommitFunc(ctx)
}

func (t *MockTransaction) Rollback(ctx context.Context) error {
	t.AddCall(ctx)
	return t.RollbackFunc(ctx)
}

type MockPgxTx struct {
	*testutil.CallWatcher
}

func NewMockPgxTx() *MockPgxTx {
	return &MockPgxTx{
		CallWatcher: testutil.NewCallWatcher(),
	}
}

func (m *MockPgxTx) Begin(ctx context.Context) (pgx.Tx, error) {
	m.AddCall(ctx)
	return nil, nil
}

func (m *MockPgxTx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error) {
	m.AddCall(ctx, f)
	return nil
}

func (m *MockPgxTx) Commit(ctx context.Context) error {
	m.AddCall(ctx)
	return nil
}

func (m *MockPgxTx) Rollback(ctx context.Context) error {
	m.AddCall(ctx)
	return nil
}

func (m *MockPgxTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	m.AddCall(ctx, tableName, columnNames, rowSrc)
	return 0, nil
}

func (m *MockPgxTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	m.AddCall(ctx, b)
	return nil
}

func (m *MockPgxTx) LargeObjects() pgx.LargeObjects {
	m.AddCall()
	return pgx.LargeObjects{}
}

func (m *MockPgxTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	m.AddCall(ctx, name, sql)
	return nil, nil
}

func (m *MockPgxTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	m.AddCall(ctx, sql, arguments)
	return nil, nil
}

func (m *MockPgxTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	m.AddCall(ctx, sql, args)
	return nil, nil
}

func (m *MockPgxTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	m.AddCall(ctx, sql, args)
	return nil
}

func (m *MockPgxTx) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	m.AddCall(ctx, sql, args, scans, f)
	return nil, nil
}

func (m *MockPgxTx) Conn() *pgx.Conn {
	m.AddCall()
	return nil
}

```
