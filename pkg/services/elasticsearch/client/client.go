package client

import (
	"context"

	"g-management/pkg/services/elasticsearch/document"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
)

type ClientInterface interface {
	CreateIndex(ctx context.Context, indexName string, request *create.Request) (*create.Response, error)
	IndexDocument(ctx context.Context, indexName string, document document.Document) (*index.Response, error)
	GetDocument(ctx context.Context, indexName, documentID string) (*get.Response, error)
	Search(ctx context.Context, indexName string, request *search.Request) (*search.Response, error)
	CheckExistIndex(ctx context.Context, indexName string) (bool, error)
}

type Client struct {
	esClient *elasticsearch.TypedClient
}

func NewClient(cfg elasticsearch.Config) (ClientInterface, error) {
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{esClient: client}, nil
}

func (c *Client) CreateIndex(ctx context.Context, indexName string, request *create.Request) (*create.Response, error) {
	res, err := c.esClient.Indices.Create(indexName).
		Request(request).
		Do(nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) IndexDocument(ctx context.Context, indexName string, document document.Document) (*index.Response, error) {
	res, err := c.esClient.Index(indexName).
		Request(document).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetDocument(ctx context.Context, indexName, documentID string) (*get.Response, error) {
	res, err := c.esClient.Get(indexName, documentID).Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Search(ctx context.Context, indexName string, request *search.Request) (*search.Response, error) {
	res, err := c.esClient.Search().
		Index(indexName).
		Request(request).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) CheckExistIndex(ctx context.Context, indexName string) (bool, error) {
	return c.esClient.Indices.Exists(indexName).IsSuccess(ctx)
}
