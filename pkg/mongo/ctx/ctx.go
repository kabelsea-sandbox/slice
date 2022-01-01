package ctxmongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// SetupContext setup db key in context via setupFunc.
func SetupContext(setupFunc func(key interface{}, value interface{}), db *mongo.Database) {
	setupFunc(dbContextKey{}, db)
}

// Extract extracts *mongo.Database from context.
func Extract(ctx context.Context) *mongo.Database {
	return ctx.Value(dbContextKey{}).(*mongo.Database)
}

// RunInTransaction runs function in transaction.
func RunInTransaction(ctx context.Context, fn func(sctx mongo.SessionContext) error) error {
	db := extractDatabase(ctx)
	if db == nil {
		log.Println("WARNING: *mongo.Database not exists in context, but ctxmongo.RunInTransaction() runs")
		return mongo.WithSession(ctx, nil, fn)
	}
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}
	if err = mongo.WithSession(ctx, session, fn); err != nil {
		if err := session.AbortTransaction(ctx); err != nil {
			return err
		}
		return err
	}
	return session.CommitTransaction(ctx)
}

//  extractDatabase
func extractDatabase(ctx context.Context) *mongo.Database {
	v := ctx.Value(dbContextKey{})
	if v == nil {
		return nil
	}
	return v.(*mongo.Database)
}

type dbContextKey struct{}
