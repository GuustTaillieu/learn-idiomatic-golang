package health

import (
	"context"
	"database/sql"
)

type DatabaseChecker struct {
	db *sql.DB
}

func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

func (h *DatabaseChecker) Ping(ctx context.Context) error {
	return h.db.PingContext(ctx)
}
