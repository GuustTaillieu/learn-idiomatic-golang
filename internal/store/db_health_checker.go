package store

import (
	"context"
	"database/sql"
)

type DBHealthChecker struct {
	db *sql.DB
}

func NewDBHealthChecker(db *sql.DB) *DBHealthChecker {
	return &DBHealthChecker{db: db}
}

func (h *DBHealthChecker) Ping(ctx context.Context) error {
	return h.db.PingContext(ctx)
}
