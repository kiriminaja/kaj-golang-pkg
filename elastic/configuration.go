package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// Config ...
func Config(cfg *Configuration) elasticsearch.Config {
	return elasticsearch.Config{
		Addresses: cfg.Address,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
}

type ElasticStructure struct {
	Shard        interface{}  `json:"_shards,omitempty"`
	Aggregations interface{}  `json:"aggregations,omitempty"`
	Data         *ElasticData `json:"hits"`
	Status       uint64       `json:"status,omitempty"`
	Error        *ErrResponse `json:"error,omitempty"`
	TimedOut     bool         `json:"timed_out"`
}

type ErrResponse struct {
	RootCause interface{} `json:"root_cause"`
	Type      string      `json:"type"`
	Reason    string      `json:"reason"`
	Line      uint64      `json:"line"`
	Column    string      `json:"column"`
}

type ElasticData struct {
	Hits     []ElasticSource `json:"hits"`
	Total    uint            `json:"total"`
	MaxScore float64         `json:"max_score"`
}
type ElasticSource struct {
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Version   uint64                 `json:"_version"`
	Highlight interface{}            `json:"highlight,omitempty"`
}

type SearchResult struct {
	Aggregations interface{} `json:"aggregations,omitempty"`
	Hits         interface{} `json:"hits"`
	Total        uint        `json:"total"`
	MaxScore     float64     `json:"max_score,omitempty"`
}
