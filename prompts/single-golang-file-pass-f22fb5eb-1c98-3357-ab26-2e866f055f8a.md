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
package invrepo

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/go-micro-example/core"
	"github.com/sksmith/go-micro-example/core/inventory"
	"github.com/sksmith/go-micro-example/db"
)

type dbRepo struct {
	conn core.Conn
}

func NewPostgresRepo(conn core.Conn) *dbRepo {
	log.Info().Msg("creating inventory repository...")
	return &dbRepo{
		conn: conn,
	}
}

func (d *dbRepo) SaveProduct(ctx context.Context, product inventory.Product, options ...core.UpdateOptions) error {
	m := db.StartMetric("SaveProduct")
	tx := db.GetUpdateOptions(d.conn, options...)

	ct, err := tx.Exec(ctx, `
		UPDATE products
           SET upc = $2, name = $3
         WHERE sku = $1;`,
		product.Sku, product.Upc, product.Name)
	if err != nil {
		m.Complete(nil)
		return errors.WithStack(err)
	}
	if ct.RowsAffected() == 0 {
		_, err := tx.Exec(ctx, `
		INSERT INTO products (sku, upc, name)
                      VALUES ($1, $2, $3);`,
			product.Sku, product.Upc, product.Name)
		if err != nil {
			m.Complete(err)
			return err
		}
	}
	m.Complete(nil)
	return nil
}

func (d *dbRepo) SaveProductInventory(ctx context.Context, productInventory inventory.ProductInventory, options ...core.UpdateOptions) error {
	m := db.StartMetric("SaveProductInventory")
	tx := db.GetUpdateOptions(d.conn, options...)

	ct, err := tx.Exec(ctx, `
		UPDATE product_inventory
           SET available = $2
         WHERE sku = $1;`,
		productInventory.Sku, productInventory.Available)
	if err != nil {
		m.Complete(nil)
		return errors.WithStack(err)
	}
	if ct.RowsAffected() == 0 {
		insert := `INSERT INTO product_inventory (sku, available)
                      VALUES ($1, $2);`
		_, err := tx.Exec(ctx, insert, productInventory.Sku, productInventory.Available)
		m.Complete(err)
		if err != nil {
			return err
		}
	}
	m.Complete(nil)
	return nil
}

func (d *dbRepo) GetProduct(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.Product, error) {
	m := db.StartMetric("GetProduct")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	product := inventory.Product{}
	err := tx.QueryRow(ctx, `SELECT sku, upc, name FROM products WHERE sku = $1 `+forUpdate, sku).
		Scan(&product.Sku, &product.Upc, &product.Name)

	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return product, errors.WithStack(core.ErrNotFound)
		}
		return product, errors.WithStack(err)
	}

	m.Complete(nil)
	return product, nil
}

func (d *dbRepo) GetProductInventory(ctx context.Context, sku string, options ...core.QueryOptions) (inventory.ProductInventory, error) {
	m := db.StartMetric("GetProductInventory")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	productInventory := inventory.ProductInventory{}
	err := tx.QueryRow(ctx, `SELECT p.sku, p.upc, p.name, pi.available FROM products p, product_inventory pi WHERE p.sku = $1 AND p.sku = pi.sku `+forUpdate, sku).
		Scan(&productInventory.Sku, &productInventory.Upc, &productInventory.Name, &productInventory.Available)

	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return productInventory, errors.WithStack(core.ErrNotFound)
		}
		return productInventory, errors.WithStack(err)
	}

	m.Complete(nil)
	return productInventory, nil
}

func (d *dbRepo) GetAllProductInventory(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]inventory.ProductInventory, error) {
	m := db.StartMetric("GetAllProducts")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	products := make([]inventory.ProductInventory, 0)
	rows, err := tx.Query(ctx,
		`SELECT p.sku, p.upc, p.name, pi.available FROM products p, product_inventory pi WHERE p.sku = pi.sku ORDER BY p.sku LIMIT $1 OFFSET $2 `+forUpdate,
		limit, offset)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return products, errors.WithStack(core.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		product := inventory.ProductInventory{}
		err = rows.Scan(&product.Sku, &product.Upc, &product.Name, &product.Available)
		if err != nil {
			m.Complete(err)
			if err == pgx.ErrNoRows {
				return nil, errors.WithStack(core.ErrNotFound)
			}
			return nil, errors.WithStack(err)
		}
		products = append(products, product)
	}

	m.Complete(nil)
	return products, nil
}

func (d *dbRepo) GetProductionEventByRequestID(ctx context.Context, requestID string, options ...core.QueryOptions) (pe inventory.ProductionEvent, err error) {
	m := db.StartMetric("GetProductionEventByRequestID")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	pe = inventory.ProductionEvent{}
	err = tx.QueryRow(ctx, `SELECT id, request_id, sku, quantity, created FROM production_events `+forUpdate+` WHERE request_id = $1 `+forUpdate, requestID).
		Scan(&pe.ID, &pe.RequestID, &pe.Sku, &pe.Quantity, &pe.Created)

	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return pe, errors.WithStack(core.ErrNotFound)
		}
		return pe, errors.WithStack(err)
	}

	m.Complete(nil)
	return pe, nil
}

func (d *dbRepo) SaveProductionEvent(ctx context.Context, event *inventory.ProductionEvent, options ...core.UpdateOptions) error {
	m := db.StartMetric("SaveProductionEvent")
	tx := db.GetUpdateOptions(d.conn, options...)

	insert := `INSERT INTO production_events (request_id, sku, quantity, created)
			       VALUES ($1, $2, $3, $4) RETURNING id;`

	err := tx.QueryRow(ctx, insert, event.RequestID, event.Sku, event.Quantity, event.Created).Scan(&event.ID)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return errors.WithStack(core.ErrNotFound)
		}
		return errors.WithStack(err)
	}
	m.Complete(nil)
	return nil
}

func (d *dbRepo) SaveReservation(ctx context.Context, r *inventory.Reservation, options ...core.UpdateOptions) error {
	m := db.StartMetric("SaveReservation")
	tx := db.GetUpdateOptions(d.conn, options...)

	insert := `INSERT INTO reservations (request_id, requester, sku, state, reserved_quantity, requested_quantity, created)
                      VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	err := tx.QueryRow(ctx, insert, r.RequestID, r.Requester, r.Sku, r.State, r.ReservedQuantity, r.RequestedQuantity, r.Created).Scan(&r.ID)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return errors.WithStack(core.ErrNotFound)
		}
		return errors.WithStack(err)
	}
	m.Complete(nil)
	return nil
}

func (d *dbRepo) UpdateReservation(ctx context.Context, ID uint64, state inventory.ReserveState, qty int64, options ...core.UpdateOptions) error {
	m := db.StartMetric("UpdateReservation")
	tx := db.GetUpdateOptions(d.conn, options...)

	update := `UPDATE reservations SET state = $2, reserved_quantity = $3 WHERE id=$1;`
	_, err := tx.Exec(ctx, update, ID, state, qty)
	m.Complete(err)
	if err != nil {
		return errors.WithStack(err)
	}
	m.Complete(nil)
	return nil
}

const reservationFields = "id, request_id, requester, sku, state, reserved_quantity, requested_quantity, created"

func (d *dbRepo) GetReservations(ctx context.Context, resOptions inventory.GetReservationsOptions, limit, offset int, options ...core.QueryOptions) ([]inventory.Reservation, error) {
	m := db.StartMetric("GetSkuOpenReserves")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	params := make([]interface{}, 0)
	params = append(params, limit)
	params = append(params, offset)

	whereClause := ""
	paramIdx := 2

	if resOptions.Sku != "" || resOptions.State != inventory.None {
		whereClause = " WHERE "
	}

	if resOptions.Sku != "" {
		if paramIdx > 2 {
			whereClause += " AND"
		}
		paramIdx++
		whereClause += " sku = $" + strconv.Itoa(paramIdx)
		params = append(params, resOptions.Sku)
	}

	if resOptions.State != inventory.None {
		if paramIdx > 2 {
			whereClause += " AND"
		}
		paramIdx++
		whereClause += " state = $" + strconv.Itoa(paramIdx)
		params = append(params, resOptions.State)
	}

	reservations := make([]inventory.Reservation, 0)
	rows, err := tx.Query(ctx,
		`SELECT `+reservationFields+` FROM reservations `+whereClause+` ORDER BY created ASC LIMIT $1 OFFSET $2 `+forUpdate,
		params...)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return reservations, errors.WithStack(core.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := inventory.Reservation{}
		err = rows.Scan(&r.ID, &r.RequestID, &r.Requester, &r.Sku, &r.State, &r.ReservedQuantity, &r.RequestedQuantity, &r.Created)
		if err != nil {
			m.Complete(err)
			return nil, err
		}
		reservations = append(reservations, r)
	}

	m.Complete(nil)
	return reservations, nil
}

func (d *dbRepo) GetReservationByRequestID(ctx context.Context, requestId string, options ...core.QueryOptions) (inventory.Reservation, error) {
	m := db.StartMetric("GetReservationByRequestID")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	r := inventory.Reservation{}
	err := tx.QueryRow(ctx,
		`SELECT `+reservationFields+` FROM reservations WHERE request_id = $1 `+forUpdate,
		requestId).Scan(&r.ID, &r.RequestID, &r.Requester, &r.Sku, &r.State, &r.ReservedQuantity, &r.RequestedQuantity, &r.Created)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return r, errors.WithStack(core.ErrNotFound)
		}
		return r, errors.WithStack(err)
	}

	m.Complete(nil)
	return r, nil
}

func (d *dbRepo) GetReservation(ctx context.Context, ID uint64, options ...core.QueryOptions) (inventory.Reservation, error) {
	m := db.StartMetric("GetReservation")
	tx, forUpdate := db.GetQueryOptions(d.conn, options...)

	r := inventory.Reservation{}
	err := tx.QueryRow(ctx,
		`SELECT `+reservationFields+` FROM reservations WHERE id = $1 `+forUpdate, ID).
		Scan(&r.ID, &r.RequestID, &r.Requester, &r.Sku, &r.State, &r.ReservedQuantity, &r.RequestedQuantity, &r.Created)
	if err != nil {
		m.Complete(err)
		if err == pgx.ErrNoRows {
			return r, errors.WithStack(core.ErrNotFound)
		}
		return r, errors.WithStack(err)
	}

	m.Complete(nil)
	return r, nil
}

func (d *dbRepo) BeginTransaction(ctx context.Context) (core.Transaction, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

```
