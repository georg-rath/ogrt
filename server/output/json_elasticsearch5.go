package output

import (
	"context"
	"log"
	"time"

	"github.com/georg-rath/ogrt/protocol"
	"gopkg.in/olivere/elastic.v5"
)

type JsonElasticSearch5Output struct {
	Emitter

	client       *elastic.Client
	index        string
	bulk         *elastic.BulkProcessor
	completionFn func(n int)
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
}

func (fw *JsonElasticSearch5Output) EmitProcessStart(msg *msg.ProcessStart) {
	msg.Time = msg.Time * int64(1000)
	req := elastic.NewBulkIndexRequest().Index(fw.index).Type("process").Doc(msg)
	fw.bulk.Add(req)
	fw.completionFn(-1)
}

func (fw *JsonElasticSearch5Output) EmitProcessEnd(msg *msg.ProcessEnd) {
	fw.completionFn(-1)
}

func (fw *JsonElasticSearch5Output) Close() {
	fw.bulk.Flush()
	fw.bulk.Close()
}
