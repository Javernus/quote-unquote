package quote

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

var insertSQL = `
INSERT INTO quote (id, message, person, created_at, updated_at, ip)
VALUES ($1, $2, $3, $3, $4, $5)
`

func (r *Repo) Insert(ctx context.Context, quote Quote) error {
	_, err := r.db.Exec(
		ctx, insertSQL, quote.ID, quote.Message, quote.Person, quote.CreatedAt.UTC(),
		quote.IP,
	)
	if err != nil {
		return fmt.Errorf("execute sql: %w", err)
	}

	return nil
}

var selectSQL = `
SELECT id, message, person, created_at, ip
FROM quote
ORDER BY created_at DESC
LIMIT $1
`

func (r *Repo) FindAll(ctx context.Context, count int) ([]Quote, error) {
	rows, err := r.db.Query(ctx, selectSQL, count)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	res := []Quote{}

	for rows.Next() {
		var quote Quote

		err = rows.Scan(&quote.ID, &quote.Message, &quote.CreatedAt, &quote.IP)
		if err != nil {
			fmt.Println(err)
			continue
		}

		quote.CreatedAt = quote.CreatedAt.UTC()

		res = append(res, quote)
	}

	return res, nil
}

var countSQL = `
SELECT COUNT(*) FROM quote
`

func (r *Repo) Count(ctx context.Context) (int, error) {
	count := 0

	err := r.db.QueryRow(ctx, countSQL).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return count, nil
}
