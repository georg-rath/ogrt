// Code generated by protoc-gen-go.
// source: protocol/ogrt.proto
// DO NOT EDIT!

/*
Package OGRT is a generated protocol buffer package.

It is generated from these files:
	protocol/ogrt.proto

It has these top-level messages:
	JobStart
	JobEnd
	SharedObject
	Module
	ProcessInfo
	JobInfo
*/
package main

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type MessageType int32

const (
	MessageType_JobStartMsg     MessageType = 10
	MessageType_JobEndMsg       MessageType = 11
	MessageType_ProcessInfoMsg  MessageType = 12
	MessageType_SharedObjectMsg MessageType = 13
	MessageType_ForkMsg         MessageType = 14
	MessageType_ExecveMsg       MessageType = 15
)

var MessageType_name = map[int32]string{
	10: "JobStartMsg",
	11: "JobEndMsg",
	12: "ProcessInfoMsg",
	13: "SharedObjectMsg",
	14: "ForkMsg",
	15: "ExecveMsg",
}
var MessageType_value = map[string]int32{
	"JobStartMsg":     10,
	"JobEndMsg":       11,
	"ProcessInfoMsg":  12,
	"SharedObjectMsg": 13,
	"ForkMsg":         14,
	"ExecveMsg":       15,
}

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}
func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (x *MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MessageType_value, data, "MessageType")
	if err != nil {
		return err
	}
	*x = MessageType(value)
	return nil
}

type JobStart struct {
	JobId            *string `protobuf:"bytes,100,req,name=job_id" json:"job_id,omitempty"`
	StartTime        *int64  `protobuf:"varint,101,req,name=start_time" json:"start_time,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *JobStart) Reset()         { *m = JobStart{} }
func (m *JobStart) String() string { return proto.CompactTextString(m) }
func (*JobStart) ProtoMessage()    {}

func (m *JobStart) GetJobId() string {
	if m != nil && m.JobId != nil {
		return *m.JobId
	}
	return ""
}

func (m *JobStart) GetStartTime() int64 {
	if m != nil && m.StartTime != nil {
		return *m.StartTime
	}
	return 0
}

type JobEnd struct {
	JobId            *string `protobuf:"bytes,200,req,name=job_id" json:"job_id,omitempty"`
	EndTime          *int64  `protobuf:"varint,201,req,name=end_time" json:"end_time,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *JobEnd) Reset()         { *m = JobEnd{} }
func (m *JobEnd) String() string { return proto.CompactTextString(m) }
func (*JobEnd) ProtoMessage()    {}

func (m *JobEnd) GetJobId() string {
	if m != nil && m.JobId != nil {
		return *m.JobId
	}
	return ""
}

func (m *JobEnd) GetEndTime() int64 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

type SharedObject struct {
	Path             *string `protobuf:"bytes,400,req,name=path" json:"path,omitempty"`
	Signature        *string `protobuf:"bytes,401,opt,name=signature" json:"signature,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SharedObject) Reset()         { *m = SharedObject{} }
func (m *SharedObject) String() string { return proto.CompactTextString(m) }
func (*SharedObject) ProtoMessage()    {}

func (m *SharedObject) GetPath() string {
	if m != nil && m.Path != nil {
		return *m.Path
	}
	return ""
}

func (m *SharedObject) GetSignature() string {
	if m != nil && m.Signature != nil {
		return *m.Signature
	}
	return ""
}

type Module struct {
	Name             *string `protobuf:"bytes,700,req,name=name" json:"name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Module) Reset()         { *m = Module{} }
func (m *Module) String() string { return proto.CompactTextString(m) }
func (*Module) ProtoMessage()    {}

func (m *Module) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type ProcessInfo struct {
	Binpath              *string         `protobuf:"bytes,300,req,name=binpath" json:"binpath,omitempty"`
	Pid                  *int32          `protobuf:"varint,301,req,name=pid" json:"pid,omitempty"`
	ParentPid            *int32          `protobuf:"varint,302,req,name=parent_pid" json:"parent_pid,omitempty"`
	Time                 *int64          `protobuf:"varint,303,req,name=time" json:"time,omitempty"`
	Signature            *string         `protobuf:"bytes,304,opt,name=signature" json:"signature,omitempty"`
	JobId                *string         `protobuf:"bytes,305,opt,name=job_id" json:"job_id,omitempty"`
	Username             *string         `protobuf:"bytes,306,opt,name=username" json:"username,omitempty"`
	Hostname             *string         `protobuf:"bytes,307,opt,name=hostname" json:"hostname,omitempty"`
	Cmdline              *string         `protobuf:"bytes,308,opt,name=cmdline" json:"cmdline,omitempty"`
	EnvironmentVariables []string        `protobuf:"bytes,309,rep,name=environment_variables" json:"environment_variables,omitempty"`
	Arguments            []string        `protobuf:"bytes,310,rep,name=arguments" json:"arguments,omitempty"`
	SharedObjects        []*SharedObject `protobuf:"bytes,311,rep,name=shared_objects" json:"shared_objects,omitempty"`
	LoadedModules        []*Module       `protobuf:"bytes,312,rep,name=loaded_modules" json:"loaded_modules,omitempty"`
	XXX_unrecognized     []byte          `json:"-"`
}

func (m *ProcessInfo) Reset()         { *m = ProcessInfo{} }
func (m *ProcessInfo) String() string { return proto.CompactTextString(m) }
func (*ProcessInfo) ProtoMessage()    {}

func (m *ProcessInfo) GetBinpath() string {
	if m != nil && m.Binpath != nil {
		return *m.Binpath
	}
	return ""
}

func (m *ProcessInfo) GetPid() int32 {
	if m != nil && m.Pid != nil {
		return *m.Pid
	}
	return 0
}

func (m *ProcessInfo) GetParentPid() int32 {
	if m != nil && m.ParentPid != nil {
		return *m.ParentPid
	}
	return 0
}

func (m *ProcessInfo) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *ProcessInfo) GetSignature() string {
	if m != nil && m.Signature != nil {
		return *m.Signature
	}
	return ""
}

func (m *ProcessInfo) GetJobId() string {
	if m != nil && m.JobId != nil {
		return *m.JobId
	}
	return ""
}

func (m *ProcessInfo) GetUsername() string {
	if m != nil && m.Username != nil {
		return *m.Username
	}
	return ""
}

func (m *ProcessInfo) GetHostname() string {
	if m != nil && m.Hostname != nil {
		return *m.Hostname
	}
	return ""
}

func (m *ProcessInfo) GetCmdline() string {
	if m != nil && m.Cmdline != nil {
		return *m.Cmdline
	}
	return ""
}

func (m *ProcessInfo) GetEnvironmentVariables() []string {
	if m != nil {
		return m.EnvironmentVariables
	}
	return nil
}

func (m *ProcessInfo) GetArguments() []string {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *ProcessInfo) GetSharedObjects() []*SharedObject {
	if m != nil {
		return m.SharedObjects
	}
	return nil
}

func (m *ProcessInfo) GetLoadedModules() []*Module {
	if m != nil {
		return m.LoadedModules
	}
	return nil
}

type JobInfo struct {
	JobId            *string        `protobuf:"bytes,400,req,name=job_id" json:"job_id,omitempty"`
	Processes        []*ProcessInfo `protobuf:"bytes,401,rep,name=processes" json:"processes,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *JobInfo) Reset()         { *m = JobInfo{} }
func (m *JobInfo) String() string { return proto.CompactTextString(m) }
func (*JobInfo) ProtoMessage()    {}

func (m *JobInfo) GetJobId() string {
	if m != nil && m.JobId != nil {
		return *m.JobId
	}
	return ""
}

func (m *JobInfo) GetProcesses() []*ProcessInfo {
	if m != nil {
		return m.Processes
	}
	return nil
}

func init() {
	proto.RegisterEnum("OGRT.MessageType", MessageType_name, MessageType_value)
}