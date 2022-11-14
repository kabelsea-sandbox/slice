package gopg

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/pkg/di"
	"github.com/kabelsea-sandbox/slice/pkg/gopg"
	ctxgopg "github.com/kabelsea-sandbox/slice/pkg/gopg/ctx"
)

// Bundle is a postgres bundle. It provides configured database instance and provide db into slice.Context.
type Bundle struct {
	Host     string `envconfig:"postgres_host" default:"postgres"`
	Port     string `envconfig:"postgres_port" default:"5432"`
	User     string `envconfig:"postgres_user"`
	Password string `envconfig:"postgres_password"`
	Database string `envconfig:"postgres_db"`
	SSLMode  bool   `envconfig:"postgres_ssl" default:"false"`
}

// Build provide database to di.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewDB, di.As(new(orm.DB)))
}

// Boot boots gopg bundle. Setup context.
func (b *Bundle) Boot(_ context.Context, interactor slice.Container) (err error) {
	return interactor.Invoke(b.SetupSliceContext)
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewDB creates db.
func (b *Bundle) NewDB() (*pg.DB, error) {
	return gopg.NewConnection(b.Host, b.Port, b.User, b.Database, b.Password)
}

// SetupSliceContext setups slice context.
func (b *Bundle) SetupSliceContext(ctx *slice.Context, db *pg.DB) {
	ctxgopg.SetupContext(ctx.Set, db)
}
