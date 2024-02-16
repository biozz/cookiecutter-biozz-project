package repository

import (
	"context"
)

type Storage interface {
	Init(ctx context.Context) error
	Close(ctx context.Context) error

	GetItems(ctx context.Context) ([]Item, error)
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
