package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/repository"
	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/storage/sqlite/migrations"
	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/storage/sqlite/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"

	_ "modernc.org/sqlite"
)

type Storage struct {
	uri string
	dev bool
	db  *bun.DB
}

func New(uri string) *Storage {
	return &Storage{uri: uri}
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) migrate(ctx context.Context) error {
	migrator := migrate.NewMigrator(s.db, migrations.Migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return fmt.Errorf("migrator init error: %w", err)
	}
	if err := migrator.Lock(ctx); err != nil {
		return fmt.Errorf("migrator lock error: %w", err)
	}
	defer migrator.Unlock(ctx) //nolint:errcheck
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}
	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
		return nil
	}
	fmt.Printf("migrated to %s\n", group)
	return nil
}

func (s *Storage) Init(ctx context.Context) error {
    // I initially somehow set it to WAL and had troubles with
    // additional files (-wal and -shm), which required that
    // I map a directory in my docker volumes instead of a single
    // file.
    // DELETE doesn't create -wal and -shm, which is more appealing
    // to me at the moment.
    uriWithParams := fmt.Sprintf("file:%s?_journal_mode=DELETE", s.uri)
	sqlDB, err := sql.Open("sqlite", uriWithParams)
	if err != nil {
		return err
	}
	s.db = bun.NewDB(sqlDB, sqlitedialect.New())
	if s.dev {
		s.db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	err = s.migrate(ctx)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}


func (s *Storage) GetItems(ctx context.Context) ([]repository.Item, error) {
	dbItems := make([]models.Item, 0)
	query := s.db.NewSelect().Model(&dbItems)
	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}
	items := make([]repository.Item, len(dbItems))
	for i := range dbItems {
		items[i] = repositoryItemFromDBItem(dbItems[i])
	}
	return items, nil
}

func repositoryItemFromDBItem(dbItem models.Item) repository.Item {
	return repository.Item{
		ID:          dbItem.ID,
		Name:        dbItem.Name,
	}
}
