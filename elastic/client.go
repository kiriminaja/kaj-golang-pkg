package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	olivere "github.com/olivere/elastic/v7"
)

type clientElastic struct {
	client        *elasticsearch.Client
	clientOlivere *olivere.Client
}

// NewClient ..
func NewClient(cfg elasticsearch.Config) (Client, error) {
	x := &clientElastic{}
	var err error
	x.clientOlivere, err = x.clientBulk(cfg)
	if err != nil {
		return x, err
	}
	x.client, err = x.connect(cfg)
	if err != nil {
		return x, err
	}
	return x, nil
}

// Connect ...
func (e *clientElastic) connect(cfg elasticsearch.Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(cfg)
}

func (e *clientElastic) clientBulk(cfg elasticsearch.Config) (*olivere.Client, error) {
	return olivere.NewClient(
		olivere.SetSniff(false),
		olivere.SetURL(cfg.Addresses...),
		olivere.SetBasicAuth(cfg.Username, cfg.Password),
	)

}

func (e *clientElastic) Version() (map[string]interface{}, error) {
	var result map[string]interface{}
	res, err := e.client.Info()
	if err != nil {
		return nil, err
	}

	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *clientElastic) Get(ctx context.Context, index string, id string) (map[string]interface{}, error) {
	res, err := e.client.Get(index, id)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *clientElastic) Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResult, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	resultFix, err := castElastic(result)
	if err != nil {
		return nil, err
	}
	return elasticPresentations(resultFix)
}

func (e *clientElastic) Mapping(ctx context.Context, index string, field interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	res, err := e.client.Indices.Create(
		index,
		e.client.Indices.Create.WithBody(strings.NewReader(string(data))),
	)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil

}

func (e *clientElastic) Insert(ctx context.Context, index, id string, field interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}
	res, err := e.client.Index(
		index,
		strings.NewReader(string(data)),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithDocumentID(id),
	)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil

}

func (e *clientElastic) Delete(ctx context.Context, index, id string) (map[string]interface{}, error) {
	res, err := e.client.Delete(index, id)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *clientElastic) BulkRequest() *olivere.BulkService {
	return e.clientOlivere.Bulk()
}

func (e *clientElastic) BulkInsert(ctx context.Context, bulkRequest *olivere.BulkService, index, id string, r interface{}) *olivere.BulkService {
	req := olivere.NewBulkIndexRequest().Index(index).Type("_doc").Id(id).Doc(r)
	return bulkRequest.Add(req)
}

func (e *clientElastic) BulkUpdate(ctx context.Context, bulkRequest *olivere.BulkService, index, id string, r interface{}) *olivere.BulkService {
	req := olivere.NewBulkUpdateRequest().Index(index).Type("_doc").Id(id).Doc(r).DocAsUpsert(true)
	return bulkRequest.Add(req)
}

func (e *clientElastic) BulkDelete(ctx context.Context, bulkRequest *olivere.BulkService, index, id string) *olivere.BulkService {
	req := olivere.NewBulkDeleteRequest().Index(index).Type("_doc").Id(id)
	return bulkRequest.Add(req)
}
