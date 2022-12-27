package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	ctxmongo "github.com/kabelsea-sandbox/slice/pkg/mongo/ctx"

	slcmongo "github.com/kabelsea-sandbox/slice/pkg/mongo"

	"github.com/kabelsea-sandbox/slice"
)

// Bundle provides mongo db integration.
type Bundle struct {
	Hosts      string `envconfig:"MONGODB_HOSTS" required:"true"`
	Database   string `envconfig:"MONGODB_DATABASE" required:"true"`
	Username   string `envconfig:"MONGODB_USERNAME" required:"true"`
	Password   string `envconfig:"MONGODB_PASSWORD" required:"true"`
	ReplicaSet string `envconfig:"MONGODB_REPLICA_SET"`
}

// Build implements bundle interface.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewDatabase)
}

// Boot implements BootShutdown interface.
func (b *Bundle) Boot(ctx context.Context, interactor slice.Container) error {
	return interactor.Invoke(b.SetupSliceContext)
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

// NewDatabase creates mongo database.
func (b *Bundle) NewDatabase() (*mongo.Database, error) {
	return slcmongo.NewConnection(b.Hosts, b.Database, b.Username, b.Password, b.ReplicaSet)
}

// SetupSliceContext setups slice context.
func (b *Bundle) SetupSliceContext(ctx *slice.Context, db *mongo.Database) {
	ctxmongo.SetupContext(ctx.Set, db)
}
