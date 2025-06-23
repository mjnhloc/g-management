package query

import (
	"encoding/json"

	"g-management/internal/models/classes/pkg/entity"
	clientRepository "g-management/pkg/services/elasticsearch/client"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	eTypes "github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/graphql-go/graphql"
)

func NewGetSearchClassesQuery(
	types map[string]*graphql.Object,
	elasticSearchClientRepo clientRepository.ClientInterface,
) *graphql.Field {
	return &graphql.Field{
		Type: types["class_elasticsearch"],
		Args: graphql.FieldConfigArgument{
			"keyword": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			keyword := params.Args["keyword"].(string)

			res, err := elasticSearchClientRepo.Search(
				params.Context,
				entity.ClassesIndexNameElasticSearch,
				&search.Request{
					Query: &eTypes.Query{
						MultiMatch: &eTypes.MultiMatchQuery{
							Query:  keyword,
							Fields: []string{"name", "description"},
						},
					},
				},
			)
			if err != nil {
				return nil, err
			}

			results := []entity.ClassDocument{}
			for _, hit := range res.Hits.Hits {
				bytes, err := hit.Source_.MarshalJSON()
				if err != nil {
					return nil, err
				}

				doc := entity.ClassDocument{}
				err = json.Unmarshal(bytes, &doc)
				if err != nil {
					return nil, err
				}

				results = append(results, doc)
			}

			return results, nil
		},
	}
}
