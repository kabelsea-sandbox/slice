package caching

import (
	"context"
	"log"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	ast "github.com/aerospike/aerospike-client-go/types"
	"github.com/eko/gocache/v3/store"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

//go:generate mockgen --package=cachingmock -destination=mocks/mock_aerospike.go . AerospikeStore

type AerospikeStore interface {
	store.StoreInterface
}

type aerospikeStore struct {
	client      *as.Client
	writePolicy *as.WritePolicy
	config      *AerospikeStoreConfig
}

type AerospikeStoreConfig struct {
	Namespace string
	SetName   string
	TTL       uint32
}

var (
	aerospikeStoreDefaultConfig = AerospikeStoreConfig{
		Namespace: "default",
		SetName:   "cache",
		TTL:       60,
	}
)

func NewAerospikeStore(client *as.Client, config *AerospikeStoreConfig) AerospikeStore {
	if err := mergo.Merge(config, aerospikeStoreDefaultConfig); err != nil {
		log.Panic(err)
	}
	return &aerospikeStore{
		client:      client,
		writePolicy: as.NewWritePolicy(0, config.TTL),
		config:      config,
	}
}

func (a *aerospikeStore) options(opts ...store.Option) *store.Options {
	options := &store.Options{}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (a *aerospikeStore) key(value any) (*as.Key, error) {
	asKey, err := as.NewKey(a.config.Namespace, a.config.SetName, value)
	if err != nil {
		return nil, errors.Wrap(err, "aerospike key failed")
	}
	return asKey, err
}

func (a *aerospikeStore) Get(ctx context.Context, key any) (any, error) {
	asKey, err := a.key(key)
	if err != nil {
		return nil, err
	}

	res, err := a.client.Get(nil, asKey, "data")
	if err != nil {
		if err != ast.ErrKeyNotFound {
			return nil, errors.Wrap(err, "aerospike get failed")
		}
		return nil, nil
	}

	val, ok := res.Bins["data"]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (a *aerospikeStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	return nil, time.Second, errors.New("get with ttl method not implemented")
}

func (a *aerospikeStore) Set(ctx context.Context, key any, value any, opts ...store.Option) error {
	asKey, err := a.key(key)
	if err != nil {
		return err
	}

	bins := as.BinMap{
		"data": value,
	}

	if err := a.client.Put(a.writePolicy, asKey, bins); err != nil {
		return errors.Wrap(err, "aerospike set failed")
	}
	return nil
}

func (a *aerospikeStore) Delete(ctx context.Context, key any) error {
	return errors.New("delete method not implemented")
}

func (a *aerospikeStore) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return errors.New("invalidate method not implemented")
}

func (a *aerospikeStore) Clear(ctx context.Context) error {
	return errors.New("clear method not implemented")
}

func (a *aerospikeStore) GetType() string {
	return "aerospike"
}
