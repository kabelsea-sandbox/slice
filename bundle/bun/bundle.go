package bun

import (
	"context"

	"slice"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// Bundle is a postgres bundle. It provides configured database instance and provide db into slice.Context.
type Bundle struct {
	ConnString string `envconfig:"postgres_conn_string" required:"True"`
}

// Build provide database to di.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewDB)
}

// Boot boots bun bundle. Setup context.
func (b *Bundle) Boot(_ context.Context, interactor slice.Container) (err error) {
	return
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewDB creates db.
func (b *Bundle) NewDB() (*bun.DB, error) {
	config, err := pgx.ParseConfig(b.ConnString)
	if err != nil {
		return nil, errors.Wrap(err, "postgres connection failed")
	}

	config.PreferSimpleProtocol = true

	db := bun.NewDB(
		stdlib.OpenDB(*config),
		pgdialect.New(),
		[]bun.DBOption{
			bun.WithDiscardUnknownColumns(),
		}...,
	)

	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		),
	)

	return db, nil
}
