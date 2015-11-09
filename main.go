package main

import (
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

func main() {
	// require that the config is name 'alter.yaml'
	file := "alter.yaml"

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// require that alter.yaml exists in the cwd
	if _, err := os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cfg, err := loadConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	cfg.path = path.Join(cwd, file)
	cfg.cwd = path.Dir(cfg.path)

	app := cli.NewApp()

	app.Usage = "Alter your machine with ease."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "a password to encrypt/decrypt a file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "provision",
			Aliases: []string{"p"},
			Usage:   "provision a machine",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					cli.ShowCommandHelp(c, "provision")
					os.Exit(1)
				}

				cmd := app.Command("decrypt")
				cmd.Run(c)

				machine := c.Args()[0]

				err = performTasks(machine, cfg)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "encrypt files",
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					cli.ShowCommandHelp(c, "encrypt")
					os.Exit(1)
				}

				err = encryptFiles(cfg.Encrypted, c.GlobalString("password"))
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "decrypt files",
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					cli.ShowCommandHelp(c, "decrypt")
					os.Exit(1)
				}

				if len(cfg.Encrypted) > 0 {
					err = decryptFiles(cfg.Encrypted, c.GlobalString("password"))
					if err != nil {
						log.Fatal(err)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
