package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Adapter interface {
	SetCollection(name string, opts *options.CollectionOptions) *mongo.Collection
	Ping(ctx context.Context)
	Watch(ctx context.Context, collection string,
		opts *options.CollectionOptions, pipeline mongo.Pipeline, optStream *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	Fetch(ctx context.Context, name string, opts *options.CollectionOptions,
		filter interface{}, optFind *options.FindOneOptions) *mongo.SingleResult
	Upsert(ctx context.Context, collection string, id uint64, data interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, collection string, id uint64) (*mongo.DeleteResult, error)
	OptionCollection() *options.CollectionOptions
	OptionChangeStream() *options.ChangeStreamOptions
	OptionFind() *options.FindOptions
	OptionFindOne() *options.FindOneOptions
	Pipeline() mongo.Pipeline
}

type Config struct {
	Username string
	Password string
	Port     string
	Host     string
	Timeout  int
	Name     string
	Debug    bool
	Replica  string
}
