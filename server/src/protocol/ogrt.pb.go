// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol/ogrt.proto

package OGRT

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MessageType int32

const (
	MessageType_JobStartMsg        MessageType = 0
	MessageType_JobEndMsg          MessageType = 11
	MessageType_ProcessInfoMsg     MessageType = 12
	MessageType_ProcessResourceMsg MessageType = 16
	MessageType_SharedObjectMsg    MessageType = 13
	MessageType_ForkMsg            MessageType = 14
	MessageType_ExecveMsg          MessageType = 15
)

var MessageType_name = map[int32]string{
	0:  "JobStartMsg",
	11: "JobEndMsg",
	12: "ProcessInfoMsg",
	16: "ProcessResourceMsg",
	13: "SharedObjectMsg",
	14: "ForkMsg",
	15: "ExecveMsg",
}

var MessageType_value = map[string]int32{
	"JobStartMsg":        0,
	"JobEndMsg":          11,
	"ProcessInfoMsg":     12,
	"ProcessResourceMsg": 16,
	"SharedObjectMsg":    13,
	"ForkMsg":            14,
	"ExecveMsg":          15,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{0}
}

type JobStart struct {
	JobId                string   `protobuf:"bytes,100,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	StartTime            int64    `protobuf:"varint,101,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobStart) Reset()         { *m = JobStart{} }
func (m *JobStart) String() string { return proto.CompactTextString(m) }
func (*JobStart) ProtoMessage()    {}
func (*JobStart) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{0}
}

func (m *JobStart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobStart.Unmarshal(m, b)
}
func (m *JobStart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobStart.Marshal(b, m, deterministic)
}
func (m *JobStart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobStart.Merge(m, src)
}
func (m *JobStart) XXX_Size() int {
	return xxx_messageInfo_JobStart.Size(m)
}
func (m *JobStart) XXX_DiscardUnknown() {
	xxx_messageInfo_JobStart.DiscardUnknown(m)
}

var xxx_messageInfo_JobStart proto.InternalMessageInfo

func (m *JobStart) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

func (m *JobStart) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

type JobEnd struct {
	JobId                string   `protobuf:"bytes,200,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	EndTime              int64    `protobuf:"varint,201,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobEnd) Reset()         { *m = JobEnd{} }
func (m *JobEnd) String() string { return proto.CompactTextString(m) }
func (*JobEnd) ProtoMessage()    {}
func (*JobEnd) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{1}
}

func (m *JobEnd) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobEnd.Unmarshal(m, b)
}
func (m *JobEnd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobEnd.Marshal(b, m, deterministic)
}
func (m *JobEnd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobEnd.Merge(m, src)
}
func (m *JobEnd) XXX_Size() int {
	return xxx_messageInfo_JobEnd.Size(m)
}
func (m *JobEnd) XXX_DiscardUnknown() {
	xxx_messageInfo_JobEnd.DiscardUnknown(m)
}

var xxx_messageInfo_JobEnd proto.InternalMessageInfo

func (m *JobEnd) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

func (m *JobEnd) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

type SharedObject struct {
	Path                 string   `protobuf:"bytes,400,opt,name=path,proto3" json:"path,omitempty"`
	Signature            string   `protobuf:"bytes,401,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SharedObject) Reset()         { *m = SharedObject{} }
func (m *SharedObject) String() string { return proto.CompactTextString(m) }
func (*SharedObject) ProtoMessage()    {}
func (*SharedObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{2}
}

func (m *SharedObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SharedObject.Unmarshal(m, b)
}
func (m *SharedObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SharedObject.Marshal(b, m, deterministic)
}
func (m *SharedObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SharedObject.Merge(m, src)
}
func (m *SharedObject) XXX_Size() int {
	return xxx_messageInfo_SharedObject.Size(m)
}
func (m *SharedObject) XXX_DiscardUnknown() {
	xxx_messageInfo_SharedObject.DiscardUnknown(m)
}

var xxx_messageInfo_SharedObject proto.InternalMessageInfo

func (m *SharedObject) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *SharedObject) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

type Module struct {
	Name                 string   `protobuf:"bytes,700,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Module) Reset()         { *m = Module{} }
func (m *Module) String() string { return proto.CompactTextString(m) }
func (*Module) ProtoMessage()    {}
func (*Module) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{3}
}

func (m *Module) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Module.Unmarshal(m, b)
}
func (m *Module) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Module.Marshal(b, m, deterministic)
}
func (m *Module) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Module.Merge(m, src)
}
func (m *Module) XXX_Size() int {
	return xxx_messageInfo_Module.Size(m)
}
func (m *Module) XXX_DiscardUnknown() {
	xxx_messageInfo_Module.DiscardUnknown(m)
}

var xxx_messageInfo_Module proto.InternalMessageInfo

func (m *Module) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// sent at start of process
type ProcessInfo struct {
	Uuid                 []byte          `protobuf:"bytes,299,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Binpath              string          `protobuf:"bytes,300,opt,name=binpath,proto3" json:"binpath,omitempty"`
	Pid                  int32           `protobuf:"varint,301,opt,name=pid,proto3" json:"pid,omitempty"`
	ParentPid            int32           `protobuf:"varint,302,opt,name=parent_pid,json=parentPid,proto3" json:"parent_pid,omitempty"`
	Time                 int64           `protobuf:"varint,303,opt,name=time,proto3" json:"time,omitempty"`
	Signature            string          `protobuf:"bytes,304,opt,name=signature,proto3" json:"signature,omitempty"`
	JobId                string          `protobuf:"bytes,305,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Username             string          `protobuf:"bytes,306,opt,name=username,proto3" json:"username,omitempty"`
	Hostname             string          `protobuf:"bytes,307,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Cmdline              string          `protobuf:"bytes,308,opt,name=cmdline,proto3" json:"cmdline,omitempty"`
	Cwd                  string          `protobuf:"bytes,313,opt,name=cwd,proto3" json:"cwd,omitempty"`
	EnvironmentVariables []string        `protobuf:"bytes,309,rep,name=environment_variables,json=environmentVariables,proto3" json:"environment_variables,omitempty"`
	Arguments            []string        `protobuf:"bytes,310,rep,name=arguments,proto3" json:"arguments,omitempty"`
	SharedObjects        []*SharedObject `protobuf:"bytes,311,rep,name=shared_objects,json=sharedObjects,proto3" json:"shared_objects,omitempty"`
	LoadedModules        []*Module       `protobuf:"bytes,312,rep,name=loaded_modules,json=loadedModules,proto3" json:"loaded_modules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ProcessInfo) Reset()         { *m = ProcessInfo{} }
func (m *ProcessInfo) String() string { return proto.CompactTextString(m) }
func (*ProcessInfo) ProtoMessage()    {}
func (*ProcessInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{4}
}

func (m *ProcessInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessInfo.Unmarshal(m, b)
}
func (m *ProcessInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessInfo.Marshal(b, m, deterministic)
}
func (m *ProcessInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessInfo.Merge(m, src)
}
func (m *ProcessInfo) XXX_Size() int {
	return xxx_messageInfo_ProcessInfo.Size(m)
}
func (m *ProcessInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessInfo proto.InternalMessageInfo

func (m *ProcessInfo) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *ProcessInfo) GetBinpath() string {
	if m != nil {
		return m.Binpath
	}
	return ""
}

func (m *ProcessInfo) GetPid() int32 {
	if m != nil {
		return m.Pid
	}
	return 0
}

func (m *ProcessInfo) GetParentPid() int32 {
	if m != nil {
		return m.ParentPid
	}
	return 0
}

func (m *ProcessInfo) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ProcessInfo) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *ProcessInfo) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

func (m *ProcessInfo) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ProcessInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *ProcessInfo) GetCmdline() string {
	if m != nil {
		return m.Cmdline
	}
	return ""
}

func (m *ProcessInfo) GetCwd() string {
	if m != nil {
		return m.Cwd
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

// sent at end of process
type ProcessResourceInfo struct {
	Uuid []byte `protobuf:"bytes,99,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Time int64  `protobuf:"varint,104,opt,name=time,proto3" json:"time,omitempty"`
	// resource info from getrusage()
	RuUtime              int64    `protobuf:"varint,105,opt,name=ru_utime,json=ruUtime,proto3" json:"ru_utime,omitempty"`
	RuStime              int64    `protobuf:"varint,106,opt,name=ru_stime,json=ruStime,proto3" json:"ru_stime,omitempty"`
	RuMaxrss             int64    `protobuf:"varint,107,opt,name=ru_maxrss,json=ruMaxrss,proto3" json:"ru_maxrss,omitempty"`
	RuMinflt             int64    `protobuf:"varint,108,opt,name=ru_minflt,json=ruMinflt,proto3" json:"ru_minflt,omitempty"`
	RuMajflt             int64    `protobuf:"varint,109,opt,name=ru_majflt,json=ruMajflt,proto3" json:"ru_majflt,omitempty"`
	RuInblock            int64    `protobuf:"varint,110,opt,name=ru_inblock,json=ruInblock,proto3" json:"ru_inblock,omitempty"`
	RuOublock            int64    `protobuf:"varint,111,opt,name=ru_oublock,json=ruOublock,proto3" json:"ru_oublock,omitempty"`
	RuNvcsw              int64    `protobuf:"varint,112,opt,name=ru_nvcsw,json=ruNvcsw,proto3" json:"ru_nvcsw,omitempty"`
	RuNivcsw             int64    `protobuf:"varint,113,opt,name=ru_nivcsw,json=ruNivcsw,proto3" json:"ru_nivcsw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessResourceInfo) Reset()         { *m = ProcessResourceInfo{} }
func (m *ProcessResourceInfo) String() string { return proto.CompactTextString(m) }
func (*ProcessResourceInfo) ProtoMessage()    {}
func (*ProcessResourceInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{5}
}

func (m *ProcessResourceInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessResourceInfo.Unmarshal(m, b)
}
func (m *ProcessResourceInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessResourceInfo.Marshal(b, m, deterministic)
}
func (m *ProcessResourceInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessResourceInfo.Merge(m, src)
}
func (m *ProcessResourceInfo) XXX_Size() int {
	return xxx_messageInfo_ProcessResourceInfo.Size(m)
}
func (m *ProcessResourceInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessResourceInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessResourceInfo proto.InternalMessageInfo

func (m *ProcessResourceInfo) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *ProcessResourceInfo) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuUtime() int64 {
	if m != nil {
		return m.RuUtime
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuStime() int64 {
	if m != nil {
		return m.RuStime
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuMaxrss() int64 {
	if m != nil {
		return m.RuMaxrss
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuMinflt() int64 {
	if m != nil {
		return m.RuMinflt
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuMajflt() int64 {
	if m != nil {
		return m.RuMajflt
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuInblock() int64 {
	if m != nil {
		return m.RuInblock
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuOublock() int64 {
	if m != nil {
		return m.RuOublock
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuNvcsw() int64 {
	if m != nil {
		return m.RuNvcsw
	}
	return 0
}

func (m *ProcessResourceInfo) GetRuNivcsw() int64 {
	if m != nil {
		return m.RuNivcsw
	}
	return 0
}

type JobInfo struct {
	JobId                string         `protobuf:"bytes,400,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Processes            []*ProcessInfo `protobuf:"bytes,401,rep,name=processes,proto3" json:"processes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *JobInfo) Reset()         { *m = JobInfo{} }
func (m *JobInfo) String() string { return proto.CompactTextString(m) }
func (*JobInfo) ProtoMessage()    {}
func (*JobInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6cb1e1d94c0282fd, []int{6}
}

func (m *JobInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobInfo.Unmarshal(m, b)
}
func (m *JobInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobInfo.Marshal(b, m, deterministic)
}
func (m *JobInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobInfo.Merge(m, src)
}
func (m *JobInfo) XXX_Size() int {
	return xxx_messageInfo_JobInfo.Size(m)
}
func (m *JobInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_JobInfo.DiscardUnknown(m)
}

var xxx_messageInfo_JobInfo proto.InternalMessageInfo

func (m *JobInfo) GetJobId() string {
	if m != nil {
		return m.JobId
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
	proto.RegisterType((*JobStart)(nil), "OGRT.JobStart")
	proto.RegisterType((*JobEnd)(nil), "OGRT.JobEnd")
	proto.RegisterType((*SharedObject)(nil), "OGRT.SharedObject")
	proto.RegisterType((*Module)(nil), "OGRT.Module")
	proto.RegisterType((*ProcessInfo)(nil), "OGRT.ProcessInfo")
	proto.RegisterType((*ProcessResourceInfo)(nil), "OGRT.ProcessResourceInfo")
	proto.RegisterType((*JobInfo)(nil), "OGRT.JobInfo")
}

func init() { proto.RegisterFile("protocol/ogrt.proto", fileDescriptor_6cb1e1d94c0282fd) }

var fileDescriptor_6cb1e1d94c0282fd = []byte{
	// 700 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x94, 0xcb, 0x6e, 0x13, 0x31,
	0x14, 0x86, 0x49, 0xd2, 0x26, 0x99, 0x93, 0x4b, 0x5b, 0x87, 0x56, 0x2e, 0x55, 0x51, 0x94, 0x55,
	0xc4, 0xa2, 0x20, 0xe8, 0x0a, 0xb1, 0x40, 0x48, 0x05, 0xb5, 0x52, 0xda, 0x6a, 0x52, 0xd8, 0x8e,
	0xe6, 0xe2, 0x26, 0x4e, 0x67, 0xec, 0x60, 0x8f, 0xd3, 0xf2, 0x0a, 0xac, 0xe0, 0x39, 0xb8, 0xdf,
	0x41, 0x62, 0xc1, 0x12, 0x9e, 0x85, 0x97, 0x40, 0xb6, 0x67, 0x92, 0xa1, 0xbb, 0xf9, 0xff, 0xcf,
	0xc7, 0x9e, 0xe3, 0xdf, 0x36, 0x74, 0xa6, 0x82, 0xa7, 0x3c, 0xe4, 0xf1, 0x4d, 0x3e, 0x12, 0xe9,
	0x8e, 0x51, 0x68, 0xe9, 0xe8, 0x91, 0x7b, 0xd2, 0xbb, 0x0f, 0xf5, 0x03, 0x1e, 0x0c, 0x53, 0x5f,
	0xa4, 0x68, 0x1d, 0xaa, 0x13, 0x1e, 0x78, 0x34, 0xc2, 0x51, 0xb7, 0xd4, 0x77, 0xdc, 0xe5, 0x09,
	0x0f, 0xf6, 0x23, 0xb4, 0x0d, 0x20, 0x35, 0xf7, 0x52, 0x9a, 0x10, 0x4c, 0xba, 0xa5, 0x7e, 0xc5,
	0x75, 0x8c, 0x73, 0x42, 0x13, 0xd2, 0xbb, 0x07, 0xd5, 0x03, 0x1e, 0xec, 0xb1, 0x08, 0x6d, 0xcc,
	0xeb, 0x7f, 0x97, 0x8a, 0x13, 0x5c, 0x83, 0x3a, 0x61, 0x91, 0x2d, 0xff, 0x53, 0x32, 0xf5, 0x35,
	0xc2, 0x22, 0x53, 0xfd, 0x00, 0x9a, 0xc3, 0xb1, 0x2f, 0x48, 0x74, 0x14, 0x4c, 0x48, 0x98, 0xa2,
	0x0e, 0x2c, 0x4d, 0xfd, 0x74, 0x8c, 0x5f, 0x54, 0xcc, 0x0c, 0x46, 0xa0, 0x6d, 0x70, 0x24, 0x1d,
	0x31, 0x3f, 0x55, 0x82, 0xe0, 0x97, 0x96, 0x2c, 0x9c, 0xde, 0x36, 0x54, 0x07, 0x3c, 0x52, 0x31,
	0xd1, 0xd5, 0xcc, 0x4f, 0x08, 0xfe, 0xb9, 0x6c, 0xab, 0xb5, 0xe8, 0xfd, 0xad, 0x40, 0xe3, 0x58,
	0xf0, 0x90, 0x48, 0xb9, 0xcf, 0x4e, 0xb9, 0x1e, 0xa4, 0x14, 0x8d, 0xf0, 0xab, 0x72, 0xb7, 0xd4,
	0x6f, 0xba, 0x46, 0xa0, 0x4d, 0xa8, 0x05, 0x94, 0x99, 0xa5, 0x5f, 0x97, 0x4d, 0x71, 0xae, 0xd1,
	0x1a, 0x54, 0xa6, 0x34, 0xc2, 0x6f, 0xb4, 0xbd, 0xec, 0xea, 0x6f, 0x74, 0x1d, 0x60, 0xea, 0x0b,
	0xc2, 0x52, 0x4f, 0x93, 0xb7, 0x96, 0x38, 0xd6, 0x3a, 0xa6, 0x91, 0x5e, 0xc2, 0x74, 0xfb, 0xae,
	0x6c, 0xba, 0x35, 0xe2, 0xff, 0x2e, 0xde, 0x97, 0x2f, 0x75, 0x51, 0xd8, 0xbd, 0x0f, 0xe5, 0xe2,
	0xee, 0x6d, 0x41, 0x5d, 0x49, 0x22, 0x4c, 0x5f, 0x1f, 0x2d, 0x99, 0x1b, 0x1a, 0x8e, 0xb9, 0x4c,
	0x0d, 0xfc, 0x94, 0xc1, 0xdc, 0xd0, 0x3d, 0x85, 0x49, 0x14, 0x53, 0x46, 0xf0, 0xe7, 0xac, 0xa7,
	0x4c, 0xeb, 0x9e, 0xc2, 0xf3, 0x08, 0xff, 0xb0, 0xb6, 0xfe, 0x46, 0xbb, 0xb0, 0x4e, 0xd8, 0x8c,
	0x0a, 0xce, 0x12, 0xdd, 0xd8, 0xcc, 0x17, 0xd4, 0x0f, 0x62, 0x22, 0xf1, 0x97, 0x72, 0xb7, 0xd2,
	0x77, 0xdc, 0xab, 0x05, 0xfa, 0x24, 0x87, 0xba, 0x29, 0x5f, 0x8c, 0x94, 0x36, 0x25, 0xfe, 0x6a,
	0x47, 0x2e, 0x1c, 0x74, 0x17, 0xda, 0xd2, 0xc4, 0xeb, 0x71, 0x93, 0xaf, 0xc4, 0xdf, 0xf4, 0x98,
	0xc6, 0x6d, 0xb4, 0xa3, 0x8f, 0xdf, 0x4e, 0x31, 0x7b, 0xb7, 0x25, 0x0b, 0x4a, 0xa2, 0x5d, 0x68,
	0xc7, 0xdc, 0x8f, 0x48, 0xe4, 0x25, 0x26, 0x5d, 0x89, 0xbf, 0xdb, 0xda, 0xa6, 0xad, 0xb5, 0x99,
	0xbb, 0x2d, 0x3b, 0xc8, 0x2a, 0xd9, 0xfb, 0x55, 0x86, 0x4e, 0x96, 0xb6, 0x4b, 0x24, 0x57, 0x22,
	0x24, 0x26, 0x75, 0x94, 0xa5, 0x1e, 0x16, 0x42, 0x47, 0x59, 0x4c, 0xe3, 0x42, 0x4a, 0x9b, 0x50,
	0x17, 0xca, 0x53, 0xc6, 0xa7, 0xf6, 0xac, 0x0a, 0xf5, 0xb8, 0x80, 0xa4, 0x41, 0x93, 0x1c, 0x0d,
	0x0d, 0xda, 0x02, 0x47, 0x28, 0x2f, 0xf1, 0x2f, 0x84, 0x94, 0xf8, 0xcc, 0xb0, 0xba, 0x50, 0x03,
	0xa3, 0x73, 0x48, 0xd9, 0x69, 0x9c, 0xe2, 0x78, 0x0e, 0x8d, 0x9e, 0x57, 0x4e, 0x34, 0x4c, 0x16,
	0x95, 0x5a, 0xeb, 0xab, 0x27, 0x94, 0x47, 0x59, 0x10, 0xf3, 0xf0, 0x0c, 0x33, 0x7b, 0xf5, 0x84,
	0xda, 0xb7, 0x46, 0x86, 0xb9, 0xb2, 0x98, 0xe7, 0xf8, 0xc8, 0x1a, 0xd9, 0xff, 0xb2, 0x59, 0x28,
	0xcf, 0xf1, 0x34, 0xff, 0xdf, 0x43, 0x2d, 0xb3, 0x55, 0x19, 0x35, 0xec, 0x69, 0xbe, 0xea, 0xa1,
	0xd1, 0xbd, 0x21, 0xd4, 0x0e, 0x78, 0x60, 0x76, 0x6d, 0x71, 0x28, 0xb3, 0x0b, 0x99, 0x1d, 0xca,
	0x5b, 0xe0, 0x4c, 0xed, 0x26, 0x13, 0xa9, 0x6f, 0xa4, 0x8e, 0x65, 0xcd, 0xc6, 0x52, 0xb8, 0x6a,
	0xee, 0x62, 0xd0, 0x8d, 0xe7, 0x25, 0x68, 0x0c, 0x88, 0x94, 0xfe, 0x88, 0x9c, 0x3c, 0x9b, 0x12,
	0xb4, 0x02, 0x8d, 0xfc, 0xe1, 0x19, 0xc8, 0xd1, 0xea, 0x15, 0xd4, 0x02, 0xc7, 0xbe, 0x23, 0x5a,
	0x36, 0x10, 0x82, 0x76, 0x61, 0x26, 0xed, 0x35, 0xd1, 0x06, 0xa0, 0x4b, 0xd1, 0x6a, 0x7f, 0x15,
	0x75, 0x60, 0xa5, 0x78, 0x90, 0xb4, 0xd9, 0x42, 0x0d, 0xa8, 0x3d, 0xe4, 0xe2, 0x4c, 0x8b, 0xb6,
	0x9e, 0x7c, 0xef, 0x82, 0x84, 0x33, 0x53, 0xb0, 0x12, 0x54, 0xcd, 0x13, 0x78, 0xe7, 0x5f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x98, 0x74, 0x5b, 0xf8, 0x19, 0x05, 0x00, 0x00,
}
