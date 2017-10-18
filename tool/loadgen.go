package main

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

func generateLoad() *ProcessInfo {
	p := &ProcessInfo{
		Binpath:   proto.String("/usr/bin/true"),
		Pid:       proto.Int32(100),
		ParentPid: proto.Int32(50),
		Signature: proto.String("9e68d822-3416-4857-b69a-5a7d17d881af"),
		JobId:     proto.String("J1"),
		Username:  proto.String("wunderwuzzi"),
		Hostname:  proto.String("localhost"),
		Cmdline:   proto.String("true"),
		Time:      proto.Int64(int64(time.Now().Nanosecond())),
	}

	buffer := make([]byte, 16384)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(MessageType_ProcessInfoMsg))
	data, _ := proto.Marshal(p)
	binary.BigEndian.PutUint32(buffer[4:8], uint32(len(data)))
	copy(buffer[8:], data)

	conn, err := net.Dial("udp", "localhost:7971")
	if err != nil {
		return nil
	}
	defer conn.Close()

	//simple write
	for {
		conn.Write(buffer[:len(data)+8])
	}
	return p
}
