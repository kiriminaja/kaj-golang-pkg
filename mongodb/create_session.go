package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kiriminaja/kaj-golang-pkg/logger"
	"github.com/kiriminaja/kaj-golang-pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func CreateSession(cfg *Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	var connStringTemplate string
	var clientOptions *options.ClientOptions
	if cfg.Replica != "" {
		connStringTemplate = "mongodb://%s:%s/?replicaSet=%s"
		clientOptions = options.Client().ApplyURI(fmt.Sprintf(connStringTemplate, cfg.Host, cfg.Port, cfg.Replica))
	} else {
		connStringTemplate = "mongodb://%s:%s"
		clientOptions = options.Client().ApplyURI(fmt.Sprintf(connStringTemplate, cfg.Host, cfg.Port))
	}

	if cfg.Username != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}
	clientOptions.SetMaxConnIdleTime(time.Duration(util.StringToInt(os.Getenv("MONGO_DB_CONN_TIMEOUT"))) * time.Second)
	clientOptions.SetMaxPoolSize(util.StrToUint64(os.Getenv("MONGO_DB_MAX_POOL")))
	clientOptions.SetMinPoolSize(util.StrToUint64(os.Getenv("MONGO_DB_MIN_POOL")))
	clientOptions.SetTimeout(time.Duration(util.StringToInt(os.Getenv("MONGO_DB_TIMEOUT"))) * time.Second)

	clientOptions.SetWriteConcern(&writeconcern.WriteConcern{
		W: "majority",
	})

	clientOptions.SetReadConcern(&readconcern.ReadConcern{
		Level: "majority",
	})

	// Membuat koneksi klien MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(logger.SetMessageFormat("Error Connect to client mongo db: %s", err.Error()))
		return nil, nil, err
	}
	database := client.Database(cfg.Name, options.Database())

	// Verifikasi koneksi
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(logger.SetMessageFormat("Error Connect to mongo database: %s", err.Error()))
		return nil, nil, err
	}
	return client, database, nil
}
