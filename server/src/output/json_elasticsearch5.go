package output

import (
	"context"
	"github.com/georg-rath/ogrt/src/protocol"
	"log"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

type JsonElasticSearch5Output struct {
	OGWriter
	client *elastic.Client
	index  string
	bulk   *elastic.BulkProcessor
}

func (fw *JsonElasticSearch5Output) Open(params string) {
	param_split := strings.Split(params, ":")
	if len(param_split) != 4 {
		panic("Wrong parameter specification for JsonOverElasticSearch - did you supply it in the format \"protocol:host:port:index\"?")
	}
	protocol := param_split[0]
	host := param_split[1]
	port := param_split[2]
	index := param_split[3]

	fw.index = index

	client, err := elastic.NewClient(elastic.SetURL(protocol+"://"+host+":"+port), elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err, ". Did you use the right version of ElasticSearch? I am expecting version 5.")
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

func (fw *JsonElasticSearch5Output) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *JsonElasticSearch5Output) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *JsonElasticSearch5Output) PersistProcessInfo(process_info *OGRT.ProcessInfo) {
	req := elastic.NewBulkIndexRequest().Index(fw.index).Type("process").Doc(process_info)
	fw.bulk.Add(req)
}

func (fw *JsonElasticSearch5Output) Close() {
	fw.bulk.Flush()
	fw.bulk.Close()
}
