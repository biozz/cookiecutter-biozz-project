package mock

import (
	"context"
	"strings"
	"sync"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/repository"
	"github.com/oklog/ulid/v2"
)

type Storage struct {
	data sync.Map
}

func New(uri string) *Storage {
	return &Storage{data: sync.Map{}}
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

type StoreKey string

const (
	itemsKey       StoreKey = "items"
	debtEntriesKey StoreKey = "debt_entries"
)

func (s *Storage) Init(ctx context.Context) error {
	s.CreateItem(ctx, repository.Item{Name: "Хлеб"})
	s.CreateItem(ctx, repository.Item{Name: "Соль"})
	return nil
}

func (s *Storage) CreateItem(ctx context.Context, item repository.Item) (repository.Item, error) {
	items := []repository.Item{}
	itemsRaw, ok := s.data.Load(itemsKey)
	if ok {
		items = itemsRaw.([]repository.Item)
	}
	id := strings.ToLower(ulid.Make().String())
	item.ID = id
	items = append(items, item)
	s.data.Store(itemsKey, items)
	return item, nil
}

func (s *Storage) GetItems(ctx context.Context) ([]repository.Item, error) {
	var items []repository.Item
	itemsRaw, ok := s.data.Load(itemsKey)
	if ok {
		items = itemsRaw.([]repository.Item)
	}
	result := make([]repository.Item, 0)
	for _, item := range items {
		result = append(result, item)
	}
	return result, nil
}

