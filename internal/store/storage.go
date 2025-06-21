package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeDuration    = time.Second * 5
	ErrorNotFound        = errors.New("resource not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	EditConflict         = errors.New("edit conflict")
)

type Storage struct {
	Movies interface {
		GetAll(ctx context.Context, title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
		Create(ctx context.Context, movie *Movie) error
		Get(ctx context.Context, id int64) (*Movie, error)
		Update(ctx context.Context, movie *Movie) error
		Delete(ctx context.Context, id int64) error
	}
}

func NewStorage(db *sql.DB) Storage {

	return Storage{
		Movies: &MovieStore{db},
	}
}
