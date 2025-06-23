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
	ErrEditConflict      = errors.New("edit conflict")
)

type Storage struct {
	Movies interface {
		GetAll(ctx context.Context, title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
		Create(ctx context.Context, movie *Movie) error
		Get(ctx context.Context, id int64) (*Movie, error)
		Update(ctx context.Context, movie *Movie) error
		Delete(ctx context.Context, id int64) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		Update(ctx context.Context, user *User) error
		GetForToken(ctx context.Context, tokenScope, tokenPlainText string) (*User, error)
	}
	Tokens interface {
		New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*Token, error)
		DeleteAllForUser(ctx context.Context, scope string, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {

	return Storage{
		Movies: &MovieStore{DB: db},
		Users:  &UserStore{DB: db},
		Tokens: &TokenStore{DB: db},
	}
}
