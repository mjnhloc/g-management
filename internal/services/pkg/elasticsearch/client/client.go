package client

import (
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

type ClientInterface interface {
	// GetClient(cfg elasticsearch.Config) *elasticsearch.Client
	CreateIndex(indexName string, mapping string) error
}

type Client struct {
	esClient *elasticsearch.Client
}

func NewClient(cfg elasticsearch.Config) (ClientInterface, error) {
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{esClient: client}, nil
}

func (c *Client) CreateIndex(indexName string, mapping string) error {
	_, err := c.esClient.Indices.Create(
		indexName,
		c.esClient.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	return nil
}
