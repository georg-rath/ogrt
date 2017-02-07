package output

import (
	"context"
	"log"
	"protocol"
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

type JsonElasticSearch5Output struct {
	OGWriter
	client *elastic.Client
	index  string
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
}

func (fw *JsonElasticSearch5Output) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *JsonElasticSearch5Output) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *JsonElasticSearch5Output) PersistProcessInfo(process_info *OGRT.ProcessInfo) {
	// set time to milliseconds
	*process_info.Time = *process_info.Time * int64(1000)
	_, err := fw.client.Index().Index(fw.index).Type("process").BodyJson(process_info).Do(context.TODO())
	if err != nil {
		log.Println("Could not index JSON in ElasticSearch: ", err)
	}
}

func (fw *JsonElasticSearch5Output) Close() {
}
