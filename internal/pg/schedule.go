package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Schedule struct {
	Name    string
	Payload map[string]string
}

func (p *PostgresSQL) GetSchedule(id string) (*Schedule, error) {
	data, err := p.usePostgresSQL(func(client *pgx.Conn, ctx context.Context) (result interface{}, err error) {
		schedule := Schedule{}
		queryErr := client.QueryRow(ctx, "SELECT name, payload FROM schedule WHERE id = $1", id).Scan(
			&schedule.Name,
			&schedule.Payload,
		)
		if queryErr != nil {
			return nil, queryErr
		}
		return schedule, nil
	})

	if err != nil {
		return nil, err
	}

	schedule, check := data.(Schedule)

	if !check {
		return nil, errors.New("failed")
	}

	return &schedule, nil
}
