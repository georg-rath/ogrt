package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/georg-rath/ogrt/protocol"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

type Librarian struct {
	*log.Logger

	elastic                  *elastic.Client
	elasticIndex             string
	elasticSearchWindow      int
	elasticResourceBatchSize int
}

func NewLibrarian(config map[string]interface{}) (lib *Librarian) {
	url := config["ElasticURL"].(string)
	index := config["ElasticIndex"].(string)
	searchWindow := config["ElasticMaxSearchWindow"].(int64)
	resourceBatchSize := config["ElasticResourceBatchSize"].(int64)

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return &Librarian{
		Logger: log.New(os.Stdout, "librarian: ", log.Flags()),

		elastic:                  client,
		elasticIndex:             index,
		elasticSearchWindow:      int(searchWindow),
		elasticResourceBatchSize: int(resourceBatchSize),
	}
}

func (l *Librarian) Start(router *gin.Engine) {
	router.POST("/job", l.GetJob)
	router.OPTIONS("/job", func(c *gin.Context) {})
}

func (l *Librarian) GetJob(c *gin.Context) {
	jobRequest := &struct {
		JobID string
	}{}

	err := c.BindJSON(jobRequest)
	if err != nil {
		abortWithError(c, 400, "Malformed request:"+err.Error())
		return
	}

	query := elastic.NewTermQuery("job_id", jobRequest.JobID)
	//query := elastic.NewBoolQuery()
	//query = query.Must(elastic.NewTermQuery("hostname", "denovo1"))
	//query = query.MustNot(elastic.NewTermQuery("job_id", "UNKNOWN"))

	// batch queries
	var procInfos []*ProcessInfo
	var searchAfter float64
	start, end := 0, l.elasticSearchWindow

	for {
		searchResult, err := l.elastic.Search().
			Index("test-start-*").
			Query(query).
			Sort("time", true).
			SearchAfter(searchAfter).
			From(start).Size(end).
			Do(context.Background())
		if err != nil {
			// Get *elastic.Error which contains additional information
			e := err.(*elastic.Error)
			l.Printf("Elastic failed with status %d and error %s.", e.Status, e.Details)
			panic(err)
		}
		l.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

		// TotalHits is another convenience function that works even when something goes wrong.
		l.Printf("Job: Found a total of %d records (%d in this request)\n", searchResult.TotalHits(), len(searchResult.Hits.Hits))
		// Here's how you iterate through the search results with full control over each step.
		if searchResult.Hits.TotalHits > 0 {
			// Iterate through results
			searchAfter = searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1].Sort[0].(float64)
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index
				var ps msg.ProcessStart
				err := json.Unmarshal(*hit.Source, &ps)
				if err != nil {
					// deserialization failed
					panic(err)
				}

				procInfos = append(procInfos, NewProcessInfoFromStartMsg(&ps))
			}

		}
		if len(searchResult.Hits.Hits) < end {
			break
		}
	}
	l.addResourcesToProcessInfos(procInfos)
	tree := l.treeFromProcessInfos(procInfos)
	c.JSON(200, map[string]interface{}{
		"data": map[string]interface{}{
			"type": "ProcessInfoTrees",
			"id":   jobRequest.JobID,
			"attributes": map[string]interface{}{
				"trees": tree,
			},
		},
	})
}

func (l *Librarian) addResourcesToProcessInfos(procInfos []*ProcessInfo) {
	// batch queries
	var start, step, end int
	step = l.elasticResourceBatchSize
	end = min(len(procInfos), step)
	for start < len(procInfos) {
		l.Println("doing batch from", start, "to", end)
		query := elastic.NewBoolQuery()
		matches := []elastic.Query{}
		for _, procInfo := range procInfos[start:end] {
			matches = append(matches, elastic.NewMatchQuery("uuid", procInfo.uuidBytes))
		}

		query.Should(matches...)

		//src, err := query.Source()
		//if err != nil {
		//	panic(err)
		//}
		//data, err := json.MarshalIndent(src, "", "  ")
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(string(data))

		searchResult, err := l.elastic.Search().
			Index("test-end-*"). // search in index "test-end"
			Query(query).        // specify the query
			From(0).Size(l.elasticResourceBatchSize).
			Do(context.Background())
		if err != nil {
			e, _ := err.(*elastic.Error)
			l.Printf("Elastic failed with status %d and error %s.", e.Status, e.Details)
		}
		l.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

		// TotalHits is another convenience function that works even when something goes wrong.
		l.Printf("Found a total of %d records\n", searchResult.TotalHits())

		// Here's how you iterate through the search results with full control over each step.
		if searchResult.Hits.TotalHits > 0 {
			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index

				var pe msg.ProcessEnd
				err := json.Unmarshal(*hit.Source, &pe)
				if err != nil {
					// Deserialization failed
					panic(err)
				}

				for _, procInfo := range procInfos {
					if bytes.Equal(procInfo.uuidBytes, pe.Uuid) {
						procInfo.AddEndMessage(&pe)
						break
					}
				}
			}
		}

		start += step
		end = min(start+step, len(procInfos)-1)
	}
}

type ProcessInfoTree []*ProcessInfoTreeNode
type ProcessInfoTreeNode struct {
	*ProcessInfo

	Children ProcessInfoTree `json:"children,omitempty"`
}

func (l *Librarian) treeFromProcessInfos(procInfos []*ProcessInfo) ProcessInfoTree {
	l.Println(len(procInfos), "to treeify")
	lookup := make(map[int32]*ProcessInfoTreeNode)
	for _, processInfo := range procInfos {
		lookup[processInfo.Pid] = &ProcessInfoTreeNode{ProcessInfo: processInfo}
	}

	tree := make([]*ProcessInfoTreeNode, 0, 3)
	for _, processInfo := range lookup {
		if parent, found := lookup[processInfo.ParentPid]; !found {
			// node has no known parent, put it onto the root
			tree = append(tree, processInfo)
			continue
		} else {
			// node has a known parent, add it to the parent
			parent.Children = append(parent.Children, processInfo)
		}
	}

	if len(tree) > 1 {
		parent := tree[0].ParentPid
		for _, node := range tree[1:] {
			if parent != node.ParentPid {
				return tree
			}
		}

		l.Println("all top-level nodes have the same parent, reparenting...")
		return []*ProcessInfoTreeNode{&ProcessInfoTreeNode{ProcessInfo: &ProcessInfo{Binpath: "autogenerated root node"}, Children: tree}}
	}
	return tree
}
