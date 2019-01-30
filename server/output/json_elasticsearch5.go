package output

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/georg-rath/ogrt/protocol"
	"github.com/olivere/elastic"
)

type JsonElasticSearch5Output struct {
	Emitter

	client *elastic.Client
	index  string
	bulk   *elastic.BulkProcessor

	indexCache    map[string]struct{}
	indexCreateMu sync.Mutex
	completionFn  func(n int)
}

func (fw *JsonElasticSearch5Output) indexForTime(t time.Time, prefix string) string {
	return fmt.Sprintf("%s-%s-%s", fw.index, prefix, t.Format("2006-01-02"))
}

// createIndexIfNotExists checks in a cache if the index exists, before it hits the elastic server
func (fw *JsonElasticSearch5Output) createIndexIfNotExists(index, mapping string) error {
	if _, found := fw.indexCache[index]; !found {
		// to prevent a race. one thread creating the index, while the other wants to create it too
		fw.indexCreateMu.Lock()
		defer fw.indexCreateMu.Unlock()

		exists, err := fw.client.IndexExists(index).Do(context.Background())
		if err != nil {
			return err
		}
		if !exists {
			createIndex, err := fw.client.CreateIndex(index).BodyString(mapping).Do(context.Background())
			if err != nil {
				return err
			}
			if !createIndex.Acknowledged {
				log.Println("creation of index", index, "not acked")
			}
			fw.indexCache[index] = struct{}{}
		} else {
			fw.indexCache[index] = struct{}{}
		}
	}
	return nil
}

func (fw *JsonElasticSearch5Output) Open(completionFn func(n int), config map[string]interface{}) {
	url := config["URL"].(string)
	index := config["Index"].(string)

	fw.index = index

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err, ". Did you use the right version of ElasticSearch? I am expecting version >= 5.")
	}
	fw.client = client

	bulk, err := client.BulkProcessor().Name("es5-bulk").
		Workers(1).
		BulkActions(1000).
		BulkSize(2 << 20).
		FlushInterval(2 * time.Second).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fw.bulk = bulk

	fw.completionFn = completionFn
	fw.indexCache = make(map[string]struct{})
}

func (fw *JsonElasticSearch5Output) EmitProcessStart(msg *msg.ProcessStart) {
	index := fw.indexForTime(convertFromMsToTime(msg.Time), "start")
	if err := fw.createIndexIfNotExists(index, startIndexMapping); err != nil {
		log.Fatal("creating start index:", err)
	}
	req := elastic.NewBulkIndexRequest().Index(index).Type("processStart").Doc(msg)
	fw.bulk.Add(req)
	fw.completionFn(-1)
}

func (fw *JsonElasticSearch5Output) EmitProcessEnd(msg *msg.ProcessEnd) {
	index := fw.indexForTime(convertFromMsToTime(msg.StartTime), "end")
	if err := fw.createIndexIfNotExists(index, endIndexMapping); err != nil {
		log.Fatal("creating end index:", err)
	}
	req := elastic.NewBulkIndexRequest().Index(index).Type("processEnd").Doc(msg)
	fw.bulk.Add(req)
	fw.completionFn(-1)
}

func (fw *JsonElasticSearch5Output) Close() {
	fw.bulk.Flush()
	fw.bulk.Close()
}

const startIndexMapping = `
{
	"mappings": {
		"processStart": {
			"properties": {
				"uuid": {
					"type": "keyword"
				},
				"binpath": {
					"type": "text"
				},
				"signature": {
					"type": "keyword"
				},
				"job_id": {
					"type": "keyword"
				},
				"username": {
					"type": "keyword"
				},
				"hostname": {
					"type": "keyword"
				},
				"parent_pid": {
					"type": "long"
				},
				"pid": {
					"type": "long"
				},
				"shared_objects": {
					"type": "nested",
					"properties": {
						"path": {
							"type": "text"
						},
						"signature": {
							"type": "keyword"
						}
					}
				},
				"time": {
					"type": "date",
					"format": "epoch_millis"
				}
			}
		}
  }
}`

const endIndexMapping = `
{
	"mappings": {
		"processEnd": {
			"properties": {
				"time": {
					"type": "date",
					"format": "epoch_millis"
				},
				"startTime": {
          "type": "date",
          "format": "epoch_millis"
        },
				"uuid": {
					"type": "keyword"
				}
			}
		}
	}
}
`
