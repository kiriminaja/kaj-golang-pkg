package elastic

import (
	"context"

	olivere "github.com/olivere/elastic/v7"
)

// Configuration ...
type Configuration struct {
	Address  []string
	Username string
	Password string
}

// Client ..
type Client interface {
	Version() (map[string]interface{}, error)
	Mapping(ctx context.Context, index string, field interface{}) (map[string]interface{}, error)
	Insert(ctx context.Context, index, id string, field interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, index, id string) (map[string]interface{}, error)
	BulkRequest() *olivere.BulkService
	BulkInsert(ctx context.Context, bulkRequest *olivere.BulkService, index, id string, r interface{}) *olivere.BulkService
	BulkUpdate(ctx context.Context, bulkRequest *olivere.BulkService, index, id string, r interface{}) *olivere.BulkService
	BulkDelete(ctx context.Context, bulkRequest *olivere.BulkService, index, id string) *olivere.BulkService
	Get(ctx context.Context, index string, id string) (map[string]interface{}, error)
	Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResult, error)
}
