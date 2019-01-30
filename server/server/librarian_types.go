package server

import (
	"github.com/georg-rath/ogrt/protocol"

	"github.com/gofrs/uuid"
)

type ProcessInfo struct {
	Uuid                 string              `json:"uuid"`
	uuidBytes            []byte              `json:"-"`
	Binpath              string              `json:"binpath"`
	Pid                  int32               `json:"pid"`
	ParentPid            int32               `json:"parent_pid"`
	StartTime            int64               `json:"start_time"`
	EndTime              int64               `json:"end_time"`
	Signature            string              `json:"signature"`
	JobId                string              `json:"job_id"`
	Username             string              `json:"username"`
	Hostname             string              `json:"hostname"`
	Cmdline              string              `json:"cmdline"`
	Cwd                  string              `json:"cwd"`
	EnvironmentVariables []string            `json:"environment_variables,omitempty"`
	SharedObjects        []*msg.SharedObject `json:"shared_objects,omitempty"`
	LoadedModules        []*msg.Module       `json:"loaded_modules,omitempty"`

	/* resource info from getrusage() */
	ResourceInfoAvailable bool  `json:"resource_info_available"`
	Ru_utime              int64 `json:"ru_utime,omitempty"`
	Ru_stime              int64 `json:"ru_stime,omitempty"`
	Ru_maxrss             int64 `json:"ru_maxrss,omitempty"`
	Ru_minflt             int64 `json:"ru_minflt,omitempty"`
	Ru_majflt             int64 `json:"ru_majflt,omitempty"`
	Ru_inblock            int64 `json:"ru_inblock,omitempty"`
	Ru_oublock            int64 `json:"ru_oublock,omitempty"`
	Ru_nvcsw              int64 `json:"ru_nvcsw,omitempty"`
	Ru_nivcsw             int64 `json:"ru_nivcsw,omitempty"`
}

func (pi *ProcessInfo) AddEndMessage(pe *msg.ProcessEnd) {
	pi.EndTime = pe.Time
	pi.Ru_utime = pe.RuUtime
	pi.Ru_stime = pe.RuStime
	pi.Ru_maxrss = pe.RuMaxrss
	pi.Ru_minflt = pe.RuMinflt
	pi.Ru_majflt = pe.RuMajflt
	pi.Ru_inblock = pe.RuInblock
	pi.Ru_oublock = pe.RuOublock
	pi.Ru_nvcsw = pe.RuNvcsw
	pi.Ru_nivcsw = pe.RuNivcsw

	pi.ResourceInfoAvailable = true
}

func NewProcessInfoFromStartMsg(ps *msg.ProcessStart) (pi *ProcessInfo) {
	return &ProcessInfo{
		Uuid:                 uuid.FromBytesOrNil(ps.Uuid).String(),
		uuidBytes:            ps.Uuid,
		Binpath:              ps.Binpath,
		Pid:                  ps.Pid,
		ParentPid:            ps.ParentPid,
		StartTime:            ps.Time,
		Signature:            ps.Signature,
		JobId:                ps.JobId,
		Username:             ps.Username,
		Hostname:             ps.Hostname,
		Cmdline:              ps.Cmdline,
		Cwd:                  ps.Cwd,
		EnvironmentVariables: ps.EnvironmentVariables,
		SharedObjects:        ps.SharedObjects,
		LoadedModules:        ps.LoadedModules,

		ResourceInfoAvailable: false,
	}
}
