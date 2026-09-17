package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// db is conenction pool -> i.e. it doesnt create a new process everytime database is queried it take a connection from the pool execs and returns it
	// default is 0 which means unlimited
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// fail fast practice
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// if connection is closed then pingContext returns error elif the ctx timeouts then ctx return error smthng like context_Deadline_Exceeded
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping:  %w", err)
	}
	return db, nil
}
