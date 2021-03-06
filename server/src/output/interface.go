package output

import (
	"github.com/georg-rath/ogrt/src/protocol"
)

type OGWriter interface {
	PersistJobStart(msg *OGRT.JobStart)
	PersistJobEnd(msg *OGRT.JobEnd)
	PersistProcessInfo(msg *OGRT.ProcessInfo)
	Open(params string)
	Close()
}
