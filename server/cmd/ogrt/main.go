package main

import (
	"fmt"
	"log"
	"os"

	"github.com/georg-rath/ogrt/pkg/inspect"
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
			{
				Name:    "inspect",
				Aliases: []string{"i"},
				Usage:   "inspect a binary signed by ogrt",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						fmt.Println("please specify a binary to inspect.")
						return nil
					}
					ExecInspect(c.Args().First())
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
	PrintMetrics     uint32
	WebAPIAddress    string
}

func ExecServer() {
	// load global config
	var gConf ServerConf
	if _, err := toml.DecodeFile("ogrt.conf", &gConf); err != nil {
		log.Fatal("error parsing configuration:", err)
	}

	server := server.New()
	server.Address = gConf.Address
	server.Port = gConf.Port
	server.MaxReceiveBuffer = gConf.MaxReceiveBuffer
	server.PrintMetrics = gConf.PrintMetrics
	if gConf.WebAPIAddress == "" {
		gConf.WebAPIAddress = ":8080"
	}
	server.WebAPIAddress = gConf.WebAPIAddress

	// load output config
	var raw interface{}
	if _, err := toml.DecodeFile("ogrt.conf", &raw); err != nil {
		log.Fatal("error parsing configuration:", err)
	}
	conf := raw.(map[string]interface{})
	if outputs, found := conf["Outputs"]; found {
		for name, config := range outputs.(map[string]interface{}) {
			server.AddOutput(name, config.(map[string]interface{}))
		}
	}

	// instantiate librarian
	if libConf, found := conf["Librarian"]; found {
		server.EnableLibrarian(libConf.(map[string]interface{}))
	}

	server.Start()

	select {}
}

func ExecInspect(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	sigs, err := inspect.FindSignatures(f)
	for _, sig := range sigs {
		fmt.Println(sig)
	}
}
