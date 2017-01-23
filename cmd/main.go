package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/supereagle/goline/config"
	"github.com/supereagle/goline/server"
)

func main() {
	// Create the command line app
	app := cli.NewApp()
	app.Name = "pipeline"
	app.Usage = "Tools to manage CI/CD pipelines for Jenkins 2.0"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.HelpFlag,
		cli.VersionFlag,
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file path",
			Value: "config.json",
		},
	}
	app.Action = func(c *cli.Context) {
		// Read the config
		path := c.String("config")
		cfg, err := config.Read(path)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		// Create and start the server
		server, err := server.NewServer(cfg)
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		err = server.Start()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}

	app.Run(os.Args)
}
