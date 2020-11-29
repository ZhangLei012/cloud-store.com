package esclient

import (
	"cloud-store.com/service/search/config"
	"github.com/olivere/elastic"
)

var client *elastic.Client

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"user_file":{
			"properties":{
				"userName":{
					"type":"keyword"
				},
				"fileSha1":{
					"type":"keyword"
				},
				"fileName":{
					"type":"keyword"
				},
				"uploadAt":{
					"type":"string"
				},
				"fileSize":{
					"type":"long"
				},
				"location":{
					"type":"string"
				}
			}
		}
	}
}`

func init() {
	if client != nil {
		return
	}
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(config.ElasticURL),
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	/*
		exists, err := client.IndexExists("cloudstore").Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !exists {
			// Create a new index.
			createIndex, err := client.CreateIndex("cloudstore").BodyString(mapping).Do(context.Background())
			if err != nil {
				// Handle error
				panic(err)
			}
			if !createIndex.Acknowledged {
				// Not acknowledged
				log.Printf("Warning: create index not acknowledged")
			}
		}
	*/
}
func Client() *elastic.Client {
	return client
}
