package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/kiriminaja/kaj-golang-pkg/logger"
	"github.com/kiriminaja/kaj-golang-pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	client   *mongo.Client
	db       *mongo.Database
	logField []logger.Field
	cfg      *Config
}

func NewMongoClient(cfg *Config) (Adapter, error) {
	x := &mongoDB{
		logField: []logger.Field{
			logger.EventName("mongodb:log"),
		},
		cfg: cfg,
	}
	session, db, err := CreateSession(cfg)
	if err != nil {
		return nil, err
	}
	x.client = session
	x.db = db
	return x, nil
}

func (m *mongoDB) OptionCollection() *options.CollectionOptions {
	return options.Collection()
}

func (m *mongoDB) OptionChangeStream() *options.ChangeStreamOptions {
	return options.ChangeStream()
}

func (m *mongoDB) OptionFind() *options.FindOptions {
	return &options.FindOptions{}
}

func (m *mongoDB) OptionFindOne() *options.FindOneOptions {
	return &options.FindOneOptions{}
}

func (m *mongoDB) SetCollection(name string, opts *options.CollectionOptions) *mongo.Collection {
	return m.db.Collection(name, opts)
}

func (m *mongoDB) Pipeline() mongo.Pipeline {
	return mongo.Pipeline{}
}

func (m *mongoDB) Find(ctx context.Context, name string, opts *options.CollectionOptions, filter interface{}, optFind *options.FindOptions) (*mongo.Cursor, error) {
	defer m.debugInfo(filter, []interface{}{
		optFind,
	}, time.Now())
	result, err := m.SetCollection(name, opts).Find(ctx, filter, optFind)
	return result, err
}

func (m *mongoDB) Fetch(ctx context.Context, name string, opts *options.CollectionOptions, filter interface{}, optFind *options.FindOneOptions) *mongo.SingleResult {
	defer m.debugInfo(filter, []interface{}{
		optFind,
	}, time.Now())
	result := m.SetCollection(name, opts).FindOne(ctx, filter, optFind)
	return result
}

func (m *mongoDB) Ping(ctx context.Context) {
	m.client.Ping(ctx, nil)
}

func (m *mongoDB) Upsert(ctx context.Context, collection string, id uint64, data interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": data}
	opts := options.Update().SetUpsert(true)
	defer m.debugInfo(data, []interface{}{
		filter,
	}, time.Now())
	return m.db.Collection(collection, options.Collection()).UpdateOne(ctx, filter, update, opts)
}

func (m *mongoDB) Delete(ctx context.Context, collection string, id uint64) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(util.ToString(id))
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	defer m.debugInfo(filter, nil, time.Now())
	return m.db.Collection(collection).DeleteOne(ctx, filter)
}

func (m *mongoDB) Watch(ctx context.Context, collection string,
	opts *options.CollectionOptions, pipeline mongo.Pipeline, optStream *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if m.cfg.Replica == "" {
		return nil, errors.New("PLEASE CONNECT TO REPLICA SET")
	}
	result, err := m.SetCollection(collection, opts).Watch(ctx, pipeline, optStream)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *mongoDB) debugInfo(query interface{}, args []interface{}, start time.Time) {
	if !m.cfg.Debug {
		return
	}

	if util.Environtment() != "prod" {
		m.logField = append(m.logField, logger.Any("args", args))

	}
	m.logField = append(m.logField, logger.Any("query", query))
	elapsed := time.Since(start)
	m.logField = append(m.logField, logger.Any("duration", elapsed.Seconds()))
	logger.Info(logger.SetMessageFormat("log"), m.logField...)

}
