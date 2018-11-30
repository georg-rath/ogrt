package output

import (
	"github.com/georg-rath/ogrt/src/protocol"
)

type NullOutput struct {
	OGWriter
}

func (fw *NullOutput) Open(params string) {
}

func (fw *NullOutput) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *NullOutput) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *NullOutput) PersistProcessResourceInfo(process_info *OGRT.ProcessResourceInfo) {
}

func (fw *NullOutput) PersistProcessInfo(process_info *OGRT.ProcessInfo) {
}

func (fw *NullOutput) Close() {
}
