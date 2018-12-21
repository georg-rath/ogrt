package output

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/georg-rath/ogrt/protocol"
)

type JsonFileOutput struct {
	Emitter

	path         string
	mu           sync.Mutex
	completionFn func(n int)
}

func (fw *JsonFileOutput) Open(completionFn func(n int), config map[string]interface{}) {
	fw.path = config["Path"].(string)
	if _, err := os.Stat(fw.path); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(fw.path, 0700)
		} else {
			log.Panic(err)
		}
	}
	fw.completionFn = completionFn
}

func (fw *JsonFileOutput) EmitProcessStart(message *msg.ProcessStart) {
	filepath := path.Join(fw.path, message.GetJobId()+".start")
	fw.mu.Lock()
	defer fw.mu.Unlock()

	var jobData []*msg.ProcessStart

	if rawJson, err := ioutil.ReadFile(filepath); err == nil {
		if err = json.Unmarshal(rawJson, &jobData); err != nil {
			log.Panic(err)
		}
	}
	jobData = append(jobData, message)

	b, err := json.Marshal(jobData)
	if err != nil {
		log.Panic(err)
	}

	file, err := os.Create(filepath)
	if _, err := file.Write(b); err != nil {
		log.Panic(err)
	}
	file.Close()

	fw.completionFn(-1)
}

func (fw *JsonFileOutput) EmitProcessEnd(msg *msg.ProcessEnd) {
	fw.completionFn(-1)
}

func (fw *JsonFileOutput) Close() {
}
