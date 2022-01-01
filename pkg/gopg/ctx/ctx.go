package ctxgopg

import (
	"context"
	"log"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// SetupContext setup db key via setupFunc.
func SetupContext(setupFunc func(key interface{}, value interface{}), db *pg.DB) {
	setupFunc(dbContextKey{}, db)
}

// Extract extracts orm.DB from context.
func Extract(ctx context.Context) orm.DB {
	if ctx.Value(txContextKey{}) == nil {
		return ctx.Value(dbContextKey{}).(orm.DB)
	}

	return ctx.Value(txContextKey{}).(orm.DB)
}

// RunInTransaction runs function with with transaction extended request.
// If transaction in base context not found, it creates it.
func RunInTransaction(ctx context.Context, fn func(ctx context.Context) (err error)) (err error) {
	if tx := ctx.Value(txContextKey{}); tx != nil {
		return fn(ctx)
	}
	db := extractDatabase(ctx)
	if db == nil {
		log.Println("WARNING: *pg.DB not exists in context, but ctxgopg.RunInTransaction() run")
		return fn(ctx)
	}
	return db.RunInTransaction(func(tx *pg.Tx) error {
		return fn(context.WithValue(ctx, txContextKey{}, tx))
	})
}

// extractDatabase
func extractDatabase(ctx context.Context) *pg.DB {
	v := ctx.Value(dbContextKey{})
	if v == nil {
		return nil
	}

	return ctx.Value(dbContextKey{}).(*pg.DB)
}

type dbContextKey struct{}
type txContextKey struct{}
