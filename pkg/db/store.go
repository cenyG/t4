package db

import (
	"context"
	"fmt"
	"log/slog"

	"T4_test_case/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Store interface {
	Select(ctx context.Context, dst any, query string, args ...interface{}) error
	Get(ctx context.Context, dst any, query string, args ...interface{}) error
	Insert(ctx context.Context, query string, args ...interface{}) (int64, error)
}

type store struct {
	db *sqlx.DB
}

func NewStore() (Store, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to db")
	}

	return &store{
		db: db,
	}, nil
}

// Select - select multiple rows
func (s *store) Select(ctx context.Context, dst any, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dst, query, args...)
}

// Get - get at most one row
func (s *store) Get(ctx context.Context, dst any, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dst, query, args...)
}

// Insert - get at most one row
func (s *store) Insert(ctx context.Context, query string, args ...interface{}) (int64, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "cant get connection")
	}

	var id int64
	err = conn.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "cannot get id")
	}

	return id, nil
}

// connectToDB - connect to pg database
func connectToDB() (*sqlx.DB, error) {
	dbCfg := config.Cfg.Common.DB
	user, password, host, port, dbName := dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DbName
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbName)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("pg connect to db: %s", connStr))

	return db, nil
}
