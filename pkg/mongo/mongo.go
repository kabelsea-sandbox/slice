package mongo

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewConnection creates new mongo connection with provided params.
func NewConnection(
	hosts string, database string, username string, password string, replicaSet string,
) (*mongo.Database, error) {
	clientOptions := options.Client().
		SetHosts(strings.Split(hosts, ",")).
		SetHeartbeatInterval(10 * time.Second).
		SetLocalThreshold(2 * time.Second).
		SetServerSelectionTimeout(2 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetRetryWrites(true).
		SetAuth(options.Credential{
			Username:   username,
			Password:   password,
			AuthSource: database,
		})

	if replicaSet != "" {
		clientOptions.
			SetReplicaSet(replicaSet).
			SetReadPreference(readpref.Primary())
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.TODO()); err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
