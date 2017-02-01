package output

import (
	"log"
	"protocol"
)

type NullOutput struct {
	OGWriter
}

func (fw *NullOutput) Open(params string) {
	log.Printf("null: open with params: '%s'", params)
}

func (fw *NullOutput) PersistJobStart(job_start *OGRT.JobStart) {
}

func (fw *NullOutput) PersistJobEnd(job_end *OGRT.JobEnd) {
}

func (fw *NullOutput) PersistProcessInfo(process_info *OGRT.ProcessInfo) {
	log.Println("null: persist process info")
}

func (fw *NullOutput) Close() {
	log.Println("null: close")
}
