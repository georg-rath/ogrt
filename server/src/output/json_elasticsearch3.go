package output

import (
	"github.com/georg-rath/ogrt/src/protocol"
	"log"
	"strings"

	"gopkg.in/olivere/elastic.v3"
)

type JsonElasticSearch3Output struct {
	OGWriter
	client *elastic.Client
	index  string
}

func (fw *JsonElasticSearch3Output) Open(params string) {
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
		log.Fatal(err, ". Did you use the right version of ElasticSearch? I am expecting version 3.")
	}
	fw.client = client
}

func (fw *JsonElasticSearch3Output) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *JsonElasticSearch3Output) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *JsonElasticSearch3Output) PersistProcessInfo(process_info *OGRT.ProcessInfo) {
	_, err := fw.client.Index().Index(fw.index).Type("process").BodyJson(process_info).Do()
	if err != nil {
		log.Println("Could not index JSON in ElasticSearch: ", err)
	}
}

func (fw *JsonElasticSearch3Output) Close() {
}
