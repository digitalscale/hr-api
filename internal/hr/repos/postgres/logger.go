package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

type logger struct{}

func (m *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if sql, ok := data["sql"]; ok {
		log.Printf("[debug] [postgres] %s", sql)
	}
}
