package main

import (
	"log"
	"os"

	"github.com/georg-rath/ogrt/server"

	"github.com/BurntSushi/toml"
	"gopkg.in/urfave/cli.v2"
)

var Version string

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	app := &cli.App{
		Name:    "ogrt",
		Usage:   "OGRT runtime tracker",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "relay and process ogrt packets",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "none",
						Usage:   "location of configuration file",
					},
				},
				Action: func(c *cli.Context) error {
					ExecServer()
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}

type ServerConf struct {
	Address          string
	Port             int
	MaxReceiveBuffer int
	DebugEndpoint    bool
	PrintMetrics     uint32
	WebAPIAddress    string
}

func ExecServer() {
	// load global config
	var gConf ServerConf
	if _, err := toml.DecodeFile("ogrt.conf", &gConf); err != nil {
		log.Fatal("error parsing configuration:", err)
	}

	s := server.New()
	s.Address = gConf.Address
	s.Port = gConf.Port
	s.MaxReceiveBuffer = gConf.MaxReceiveBuffer
	s.DebugEndpoint = gConf.DebugEndpoint
	s.PrintMetrics = gConf.PrintMetrics
	if gConf.WebAPIAddress == "" {
		gConf.WebAPIAddress = ":8080"
	}
	s.WebAPIAddress = gConf.WebAPIAddress

	// load output config
	var raw interface{}
	if _, err := toml.DecodeFile("ogrt.conf", &raw); err != nil {
		log.Fatal("error parsing configuration:", err)
	}
	conf := raw.(map[string]interface{})
	if outputs, found := conf["Outputs"]; found {
		for name, config := range outputs.(map[string]interface{}) {
			s.AddOutput(name, config.(map[string]interface{}))
		}
	}

	s.Start()

	select {}
}
