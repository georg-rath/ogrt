package output

import (
	"github.com/georg-rath/ogrt/protocol"
)

type Emitter interface {
	Open(completionFn func(n int), config map[string]interface{})
	Close()

	EmitProcessStart(msg *msg.ProcessStart)
	EmitProcessEnd(msg *msg.ProcessEnd)
}

var DefaultCompletionFn func(n int)
