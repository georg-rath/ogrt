package server

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/georg-rath/ogrt/output"
	"github.com/georg-rath/ogrt/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/rcrowley/go-metrics"
	"github.com/vrischmann/go-metrics-influxdb"
)

type Server struct {
	*log.Logger

	Address          string
	Port             int
	MaxReceiveBuffer int
	PrintMetrics     uint32
	WebAPIAddress    string

	Librarian *Librarian

	InfluxMetrics InfluxMetrics

	outputs           map[string]output.Emitter
	outstandingOutput sync.WaitGroup
	inShutdown        bool

	listener *net.UDPConn
}

type InfluxMetrics struct {
	Interval uint32
	URL      string
	Database string
	User     string
	Password string
}

func New() (s *Server) {
	s = &Server{}
	s.outputs = make(map[string]output.Emitter)
	s.Logger = log.New(os.Stdout, "server: ", log.Flags())
	output.DefaultCompletionFn = func(n int) {
		s.outstandingOutput.Add(n)
	}

	return
}

func (s *Server) Stop() {
	s.inShutdown = true
	s.listener.Close()
	log.Println("Waiting for outstanding output...")
	s.outstandingOutput.Wait()
}

func (s *Server) AddOutput(name string, config map[string]interface{}) {
	typ, ok := config["Type"]
	if !ok {
		s.Fatalf("no type specified for output '%s'\n", name)
	}

	var emitter output.Emitter
	switch typ {
	case "JsonOverTcp":
		emitter = &output.JsonOverTcpOutput{}
	case "JsonElasticSearch":
		fallthrough
	case "JsonElasticSearch5":
		emitter = &output.JsonElasticSearch5Output{}
	case "PgSqlAggregator":
		emitter = &output.PgSqlAggregatorOutput{}
	case "JsonFile":
		emitter = &output.JsonFileOutput{}
	case "Null":
		emitter = &output.NullOutput{}
	default:
		s.Fatalf("unkown output type '%s' for output '%s'", typ, name)
	}

	s.outputs[name] = emitter
	metrics.Register("output_"+name, metrics.NewTimer())

	s.outputs[name].Open(output.DefaultCompletionFn, config)
}

func (s *Server) EnableLibrarian(config map[string]interface{}) {
	s.Librarian = NewLibrarian(config)
}

func (s *Server) Start() {
	if s.WebAPIAddress != "" {
		go StartWebAPI(s.WebAPIAddress, s.Librarian)
	}

	// Listen for incoming connections.
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.Address, s.Port))
	if err != nil {
		s.Fatal("Error resolving UDP address:", err.Error())
	}

	if s.listener, err = net.ListenUDP("udp", addr); err != nil {
		s.Fatal("Error listening:", err.Error())
	}

	/* register timer for receive() */
	receive_timer := metrics.NewTimer()
	metrics.Register("input_receive", receive_timer)

	/* output metrics on stderr */
	if s.PrintMetrics > 0 {
		s.Printf("printing metrics every %d seconds", s.PrintMetrics)
		go metrics.LogScaled(metrics.DefaultRegistry, time.Duration(s.PrintMetrics)*time.Second, time.Millisecond, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}

	if s.InfluxMetrics.Interval > 0 {
		s.Printf("sending metrics every %d seconds to %s (db: %s) as %s", s.InfluxMetrics.Interval, s.InfluxMetrics.URL, s.InfluxMetrics.Database, s.InfluxMetrics.User)
		go influxdb.InfluxDB(
			metrics.DefaultRegistry,                             // metrics registry
			time.Second*time.Duration(s.InfluxMetrics.Interval), // interval
			s.InfluxMetrics.URL,                                 // the InfluxDB url
			s.InfluxMetrics.Database,                            // your InfluxDB database
			s.InfluxMetrics.User,                                // your InfluxDB user
			s.InfluxMetrics.Password,                            // your InfluxDB password
		)
	}

	go func() {
		bufferPool := NewStaticPool("pool_receive", s.MaxReceiveBuffer)

		// Read the data waiting on the connection and put it in the data buffer
		for {
			packetBuffer := bufferPool.Get()
			// Read header from the connection
			n, addr, err := s.listener.ReadFromUDP(packetBuffer)
			if err == io.EOF {
				bufferPool.Put(packetBuffer)
				continue
			} else if err != nil {
				bufferPool.Put(packetBuffer)
				if s.inShutdown == true {
					break
				}
				s.Printf("error while receiving from %s, bytes: %d, error: %s", addr, n, err)
				continue
			}

			receive_timer.Time(func() {
				// Decode type and length of packet from header
				msgType := binary.BigEndian.Uint32(packetBuffer[0:4])
				msgLength := binary.BigEndian.Uint32(packetBuffer[4:8])

				data := packetBuffer[8 : msgLength+8]

				go func() {
					switch msgType {
					case uint32(msg.MessageType_ProcessStartMsg):
						ps := &msg.ProcessStart{}

						err = proto.Unmarshal(data, ps)
						if err != nil {
							s.Printf("Error decoding ProcessStart message: %s\n", err)
							bufferPool.Put(packetBuffer)
							return
						}

						s.Printf("ProcessStart: JobId=%s, Pid=%d, UUID=%s\n", ps.JobId, ps.Pid, hex.EncodeToString(ps.Uuid))
						for name, output := range s.outputs {
							metric := metrics.Get("output_" + name).(metrics.Timer)
							metric.Time(func() {
								s.outstandingOutput.Add(1)
								output.EmitProcessStart(ps)
							})
						}
					case uint32(msg.MessageType_ProcessEndMsg):
						pe := &msg.ProcessEnd{}

						err = proto.Unmarshal(data, pe)
						if err != nil {
							s.Printf("Error decoding ProcessEnd message: %s\n", err)
							bufferPool.Put(packetBuffer)
							return
						}

						s.Printf("ProcessEnd: UUID=%s\n", hex.EncodeToString(pe.Uuid))

						for name, output := range s.outputs {
							metric := metrics.Get("output_" + name).(metrics.Timer)
							metric.Time(func() {
								s.outstandingOutput.Add(1)
								output.EmitProcessEnd(pe)
							})
						}
					default:
						s.Println("unkown message type", msgType)
						bufferPool.Put(packetBuffer)
						return
					}
				}()
			})
		}

		s.Println(bufferPool.InFlight(), "packet buffers still in use")
	}()
}
