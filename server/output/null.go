package output

import (
	"github.com/georg-rath/ogrt/protocol"
)

type NullOutput struct {
	Emitter

	completionFn func(n int)
}

func (fw *NullOutput) Open(completionFn func(n int), config map[string]interface{}) {
	fw.completionFn = completionFn
}

func (fw *NullOutput) EmitProcessStart(msg *msg.ProcessStart) {
	fw.completionFn(-1)
}

func (fw *NullOutput) EmitProcessEnd(msg *msg.ProcessEnd) {
	fw.completionFn(-1)
}

func (fw *NullOutput) Close() {
}
