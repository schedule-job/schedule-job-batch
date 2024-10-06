package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Request struct {
	Name    string
	Payload map[string]interface{}
}

func (p *PostgresSQL) GetRequests(id string) (*Request, error) {
	data, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		request := Request{}
		queryErr := client.QueryRow(ctx, "SELECT name, payload FROM request WHERE id = $1", id).Scan(
			&request.Name,
			&request.Payload,
		)
		if queryErr != nil {
			return nil, queryErr
		}
		return request, nil
	})

	if err != nil {
		return nil, err
	}

	request, check := data.(Request)

	if !check {
		return nil, errors.New("failed")
	}

	return &request, nil
}
