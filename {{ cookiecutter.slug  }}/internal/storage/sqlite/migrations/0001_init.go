package migrations

import (
	"context"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/internal/storage/sqlite/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*models.Item)(nil)).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*models.Item)(nil)).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
