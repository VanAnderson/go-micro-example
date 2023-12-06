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
