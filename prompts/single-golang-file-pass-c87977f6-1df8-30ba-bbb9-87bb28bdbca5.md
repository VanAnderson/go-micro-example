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



### repo.go
```
package usrrepo

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/user"
	"github.com/sksmith/go-micro-example/db"

	lru "github.com/hashicorp/golang-lru"
)

type dbRepo struct {
	conn core.Conn
	c    *lru.Cache
}

func NewPostgresRepo(conn core.Conn) user.Repository {
	log.Info().Msg("creating user repository...")

	l, err := lru.New(256)
	if err != nil {
		log.Warn().Err(err).Msg("unable to configure cache")
	}
	return &dbRepo{
		conn: conn,
		c:    l,
	}
}

func (r *dbRepo) Create(ctx context.Context, user *user.User, txs ...core.UpdateOptions) error {
	m := db.StartMetric("Create")
	tx := db.GetUpdateOptions(r.conn, txs...)

	_, err := tx.Exec(ctx, `
		INSERT INTO users (username, password, is_admin, created_at)
                      VALUES ($1, $2, $3, $4);`,
		user.Username, user.HashedPassword, user.IsAdmin, user.Created)
	if err != nil {
		m.Complete(err)
		return err
	}
	r.cache(*user)
	m.Complete(nil)
	return nil
}

func (r *dbRepo) Get(ctx context.Context, username string, txs ...core.QueryOptions) (user.User, error) {
	m := db.StartMetric("GetUser")
	tx, forUpdate := db.GetQueryOptions(r.conn, txs...)

	u, ok := r.getcache(username)
	if ok {
		return u, nil
	}

	query := `SELECT username, password, is_admin, created_at FROM users WHERE username = $1 ` + forUpdate

	log.Debug().Str("query", query).Str("username", username).Msg("getting user")

	err := tx.QueryRow(ctx, query, username).
		Scan(&u.Username, &u.HashedPassword, &u.IsAdmin, &u.Created)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return user.User{}, errors.WithStack(core.ErrNotFound)
		}
		return user.User{}, errors.WithStack(err)
	}

	r.cache(u)
	m.Complete(nil)
	return u, nil
}

func (r *dbRepo) Delete(ctx context.Context, username string, txs ...core.UpdateOptions) error {
	m := db.StartMetric("DeleteUser")
	tx := db.GetUpdateOptions(r.conn, txs...)

	_, err := tx.Exec(ctx, `DELETE FROM users WHERE username = $1`, username)

	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return errors.WithStack(core.ErrNotFound)
		}
		return errors.WithStack(err)
	}

	r.uncache(username)
	m.Complete(nil)
	return nil
}

func (r *dbRepo) cache(u user.User) {
	if r.c == nil {
		return
	}
	r.c.Add(u.Username, u)
}

func (r *dbRepo) uncache(username string) {
	if r.c == nil {
		return
	}
	r.c.Remove(username)
}

func (r *dbRepo) getcache(username string) (user.User, bool) {
	if r.c == nil {
		return user.User{}, false
	}

	v, ok := r.c.Get(username)
	if !ok {
		return user.User{}, false
	}
	u, ok := v.(user.User)
	return u, ok
}

```
