package services

import (
	"context"

	"g-management/internal/models/classes/pkg/entity"
	"g-management/pkg/services/elasticsearch/client"
	"g-management/pkg/services/elasticsearch/document"

	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type ClassesServiceInterface interface {
	CheckExistAndIndexNewClassDoc(ctx context.Context, document document.Document) error
}

type classesService struct {
	esClientRepo client.ClientInterface
}

func NewClassesService(esClientRepo client.ClientInterface) ClassesServiceInterface {
	return &classesService{
		esClientRepo: esClientRepo,
	}
}

func (cs *classesService) CheckExistAndIndexNewClassDoc(ctx context.Context, document document.Document) error {
	exist, err := cs.esClientRepo.CheckExistIndex(ctx, entity.ClassesIndexNameElasticSearch)
	if err != nil {
		return err
	}

	if !exist {
		_, err := cs.esClientRepo.CreateIndex(
			ctx,
			entity.ClassesIndexNameElasticSearch,
			&create.Request{
				Mappings: &types.TypeMapping{
					Properties: map[string]types.Property{
						"name": &types.TextProperty{
							Type: "text",
						},
						"description": &types.TextProperty{
							Type: "text",
						},
					},
				},
			},
		)
		if err != nil {
			return err
		}
	}

	_, err = cs.esClientRepo.IndexDocument(ctx, document.IndexName(), document)
	if err != nil {
		return err
	}

	return nil
}
