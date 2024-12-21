package elastic

import (
	"encoding/json"
	"errors"

	"github.com/kiriminaja/kaj-golang-pkg/util"
)

func castElastic(v interface{}) (*ElasticStructure, error) {
	result := &ElasticStructure{}
	err := json.Unmarshal([]byte(util.DumpToString(v)), result)
	if err != nil {
		return nil, err
	}
	if result.Status == 400 {
		return nil, errors.New(result.Error.Reason)
	}
	return result, nil
}

func elasticPresentations(elasticData *ElasticStructure) (*SearchResult, error) {
	data := make([]interface{}, 0)
	for _, i := range elasticData.Data.Hits {
		i.Source["_version"] = i.Version
		data = append(data, i.Source)
	}
	result := &SearchResult{
		Hits:         data,
		Aggregations: elasticData.Aggregations,
		Total:        elasticData.Data.Total,
		MaxScore:     elasticData.Data.MaxScore,
	}
	return result, nil
}
