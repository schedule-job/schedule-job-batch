package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Request struct {
	ID      string
	Name    string
	Payload map[string]interface{}
}

func (p *PostgresSQL) GetRequest(id string) (*Request, error) {
	data, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		request := Request{}
		queryErr := client.QueryRow(ctx, "SELECT job_id, name, payload FROM request WHERE job_id = $1", id).Scan(
			&request.ID,
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

func (p *PostgresSQL) GetIdsByRequests() ([]string, error) {
	data, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		rows, queryErr := client.Query(ctx, "SELECT job_id FROM request")
		if err != nil {
			return nil, queryErr
		}
		ids := []string{}
		for rows.Next() {
			id := ""
			scanErr := rows.Scan(&id)
			ids = append(ids, id)
			if scanErr != nil {
				continue
			}
		}
		return ids, nil
	})

	if err != nil {
		return nil, err
	}

	ids, check := data.([]string)

	if !check {
		return nil, errors.New("failed")
	}

	return ids, nil
}

func (p *PostgresSQL) GetRequests() (*[]Request, error) {
	data, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		rows, queryErr := client.Query(ctx, "SELECT job_id, name, payload FROM request")
		if err != nil {
			return nil, queryErr
		}
		requests := []Request{}
		for rows.Next() {
			request := Request{}
			scanErr := rows.Scan(&request.ID,
				&request.Name,
				&request.Payload)
			requests = append(requests, request)
			if scanErr != nil {
				continue
			}
		}
		return requests, nil
	})

	if err != nil {
		return nil, err
	}

	requests, check := data.([]Request)

	if !check {
		return nil, errors.New("failed")
	}

	return &requests, nil
}
