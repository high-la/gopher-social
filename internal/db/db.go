package db

import (
	"context"
	"database/sql"
	"time"
)

// constructor
func New(addr string, maxOpenConnections int, maxIdleConnections int, maxIdleTime time.Duration) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxIdleConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, err
}
