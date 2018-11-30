package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/georg-rath/ogrt/src/output"
	"github.com/georg-rath/ogrt/src/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/vrischmann/go-metrics-influxdb"
)

var Version string

var config Configuration

type Output struct {
	Type    string
	Params  string
	Workers int
	Writer  output.OGWriter
}

type InfluxMetrics struct {
	Interval uint32
	URL      string
	Database string
	User     string
	Password string
}

type Configuration struct {
	Address          string
	Port             int
	MaxReceiveBuffer uint32
	DebugEndpoint    bool
	PrintMetrics     uint32
	WebAPI           bool
	WebAPIAddress    string
	Outputs          map[string]Output
	InfluxMetrics    InfluxMetrics
}

var exitChannel chan bool

var outputs map[string][]Output
var output_channels map[string]chan interface{}
var outstandingOutput sync.WaitGroup
var output_wait sync.WaitGroup
var shutdown = false

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("ogrt-server %s", Version)

	exitChannel = make(chan bool)

	if _, err := toml.DecodeFile("ogrt.conf", &config); err != nil {
		log.Fatal(err)
	}

	/* expose metrics as HTTP endpoint */
	if config.DebugEndpoint == true {
		exp.Exp(metrics.DefaultRegistry)
		go http.ListenAndServe(":8080", nil)
		log.Printf("Instantiated DebugEndpoint at Port 8080 (http://0.0.0.0:8080/debug/metrics)")
	}

	if config.WebAPI == true {
		if config.WebAPIAddress == "" {
			config.WebAPIAddress = ":8080"
		}
		go StartWebAPI(config.WebAPIAddress)
	}

	// Listen for incoming connections.
	listen_string := fmt.Sprintf("%s:%d", config.Address, config.Port)
	ServerAddr, err := net.ResolveUDPAddr("udp", listen_string)
	if err != nil {
		log.Fatal("Error resolving UDP address:", err.Error())
	}
	listener, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer listener.Close()

	outputs = make(map[string][]Output)
	output_channels = make(map[string]chan interface{})

	/* instantiate all outputs */
	for name, out := range config.Outputs {
		output_channels[name] = make(chan interface{})
		for worker_id := 0; worker_id < config.Outputs[name].Workers; worker_id++ {
			var output_ Output
			switch out.Type {
			case "JsonOverTcp":
				output_.Writer = new(output.JsonOverTcpOutput)
			case "JsonElasticSearch":
				fallthrough
			case "JsonElasticSearch3":
				output_.Writer = new(output.JsonElasticSearch3Output)
			case "JsonElasticSearch5":
				output_.Writer = new(output.JsonElasticSearch5Output)
			case "PgSqlAggregator":
				output_.Writer = new(output.PgSqlAggregatorOutput)
			case "JsonFile":
				output_.Writer = new(output.JsonFileOutput)
			case "Null":
				output_.Writer = new(output.NullOutput)
			default:
				log.Fatal("Unkown output type: ", out.Type)
			}
			output_.Writer.Open(out.Params)

			outputs[name] = append(outputs[name], output_)
			go writeToOutput(name, worker_id, &output_, output_channels[name])
		}

		metrics.Register("output_"+name, metrics.NewTimer())
		log.Printf("Instantiated output '%s' of type '%s' with parameters: '%s'", name, config.Outputs[name].Type, config.Outputs[name].Params)
	}

	/* Setup signal handler for SIGKILL and SIGTERM */
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.\n", sig)
		shutdown = true
		listener.Close()
		log.Println("Waiting for outstanding output...")
		outstandingOutput.Wait()
		for _, c := range output_channels {
			close(c)
		}
		output_wait.Wait()
		exitChannel <- true
	}(sigc)

	/* register timer for receive() */
	receive_timer := metrics.NewTimer()
	metrics.Register("input_receive", receive_timer)

	/* output metrics on stderr */
	if config.PrintMetrics > 0 {
		log.Printf("printing metrics every %d seconds", config.PrintMetrics)
		go metrics.LogScaled(metrics.DefaultRegistry, time.Duration(config.PrintMetrics)*time.Second, time.Millisecond, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}

	if config.InfluxMetrics.Interval > 0 {
		log.Printf("sending metrics every %d seconds to %s (db: %s) as %s", config.InfluxMetrics.Interval, config.InfluxMetrics.URL, config.InfluxMetrics.Database, config.InfluxMetrics.User)
		go influxdb.InfluxDB(
			metrics.DefaultRegistry, // metrics registry
			time.Second*time.Duration(config.InfluxMetrics.Interval), // interval
			config.InfluxMetrics.URL,                                 // the InfluxDB url
			config.InfluxMetrics.Database,                            // your InfluxDB database
			config.InfluxMetrics.User,                                // your InfluxDB user
			config.InfluxMetrics.Password,                            // your InfluxDB password
		)
	}

	bufferPool := NewStaticPool(int(config.MaxReceiveBuffer))

	// Read the data waiting on the connection and put it in the data buffer
	for {
		packetBuffer := bufferPool.Get()
		// Read header from the connection
		n, addr, err := listener.ReadFromUDP(packetBuffer)
		if err == io.EOF {
			bufferPool.Put(packetBuffer)
			continue
		} else if err != nil {
			bufferPool.Put(packetBuffer)
			if shutdown == true {
				break
			}
			log.Printf("error while receiving from %s, bytes: %d, error: %s", addr, n, err)
			continue
		}

		receive_timer.Time(func() {
			// Decode type and length of packet from header
			msg_type := binary.BigEndian.Uint32(packetBuffer[0:4])
			msg_length := binary.BigEndian.Uint32(packetBuffer[4:8])

			// allocate a buffer as big as the payload and read the rest of the packet
			data := packetBuffer[8 : msg_length+8]

			go func() {
				var msg interface{}

				switch msg_type {
				case uint32(OGRT.MessageType_ProcessInfoMsg):
					pi := &OGRT.ProcessInfo{}

					err = proto.Unmarshal(data, pi)
					if err != nil {
						log.Printf("Error decoding ExecveMsg: %s\n", err)
						bufferPool.Put(packetBuffer)
						return
					}
					msg = pi
				case uint32(OGRT.MessageType_ProcessResourceMsg):
					pri := &OGRT.ProcessResourceInfo{}

					err = proto.Unmarshal(data, pri)
					if err != nil {
						log.Printf("Error decoding ResourceInfoMsg: %s\n", err)
						bufferPool.Put(packetBuffer)
						return
					}

					msg = pri
				default:
					log.Println("unkown message type", msg_type)
					return
				}
				bufferPool.Put(packetBuffer)

				for _, c := range output_channels {
					outstandingOutput.Add(1)
					c <- msg
				}
			}()
		})
	}

	<-exitChannel
	log.Println(bufferPool.InFlight(), "packet buffers still in use")
	log.Println("Thank you for using OGRT.")
}

func writeToOutput(name string, id int, output *Output, messages chan interface{}) {
	output_wait.Add(1)
	for message := range messages {
		switch message := message.(type) {
		default:
			log.Printf("unexpected type %T", message)
			outstandingOutput.Done()
		case *OGRT.ProcessInfo:
			metric := metrics.Get("output_" + name).(metrics.Timer)
			metric.Time(func() {
				output.Writer.PersistProcessInfo(message)
			})
			outstandingOutput.Done()
		case *OGRT.ProcessResourceInfo:
			metric := metrics.Get("output_" + name).(metrics.Timer)
			metric.Time(func() {
				output.Writer.PersistProcessResourceInfo(message)
			})
			outstandingOutput.Done()
		}
	}
	output.Writer.Close()
	log.Printf("output %s [%d]: closed output.", name, id)
	output_wait.Done()
}

type StaticPool struct {
	pool     sync.Pool
	size     int
	inFlight int64
}

func (s StaticPool) Get() []byte {
	atomic.AddInt64(&s.inFlight, 1)
	return s.pool.Get().([]byte)
}

func (s StaticPool) Put(buf []byte) {
	if cap(buf) != s.size {
		panic("buffer with wrong size returned to pool")
	}
	atomic.AddInt64(&s.inFlight, -1)
	s.pool.Put(buf)
}

func (s StaticPool) InFlight() int64 {
	inFlight := atomic.LoadInt64(&s.inFlight)
	return inFlight
}

func NewStaticPool(size int) *StaticPool {
	return &StaticPool{
		pool: sync.Pool{
			New: func() interface{} { return make([]byte, size) },
		},
		size: size,
	}
}
