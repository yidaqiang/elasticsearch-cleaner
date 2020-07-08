package types

import (
	"github.com/olivere/elastic"
	"testing"
)

func TestElasticSearchManager_DeleteIndex(t *testing.T) {
	client, _ := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"), elastic.SetSniff(false))
	manager := ElasticSearchManager{Client: client}
	manager.DeleteIndex("test-bulk-example")
}

func TestElasticSearchManager_DeleteDocs(t *testing.T) {
	client, _ := elastic.NewClient(elastic.SetURL("http://192.168.72.41:30200"), elastic.SetSniff(false))
	manager := ElasticSearchManager{Client: client}
	manager.DeleteDocs("mysql", 1)
}
