package output

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/georg-rath/ogrt/protocol"
)

type JsonOverTcpOutput struct {
	Emitter

	encoder      *json.Encoder
	completionFn func(n int)
}

func (fw *JsonOverTcpOutput) Open(completionFn func(n int), config map[string]interface{}) {
	host := config["host"].(string)
	port := config["port"].(int)

	connex, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}
	fw.encoder = json.NewEncoder(connex)
	fw.completionFn = completionFn
}

func (fw *JsonOverTcpOutput) EmitProcessStart(ps *msg.ProcessStart) {
	if err := fw.encoder.Encode(ps); err != nil {
		log.Println("Could not send JSON over TCP: ", err)
	}
	fw.completionFn(-1)
}

func (fw *JsonOverTcpOutput) EmitProcessEnd(pe *msg.ProcessEnd) {
	if err := fw.encoder.Encode(pe); err != nil {
		log.Println("Could not send JSON over TCP: ", err)
	}
	fw.completionFn(-1)
}

func (fw *JsonOverTcpOutput) Close() {
}
