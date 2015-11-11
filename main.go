package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

type globalFlags struct {
	password string
}

type provisionFlags struct {
	global   *globalFlags
	links    bool
	commands bool
	parents  bool
	clobber  bool
}

type encryptFlags struct {
	global *globalFlags
	remove bool
}

type decryptFlags struct {
	global *globalFlags
	remove bool
}

func main() {
	var ignoreArgCheck bool

	global := &globalFlags{}

	app := cli.NewApp()

	app.Usage = "Alter your machine with ease."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "password to encrypt/decrypt a file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "provision",
			Aliases: []string{"p"},
			Usage:   "provision a machine",
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "links",
					Usage: "provision links, defaults to true",
				},
				cli.BoolTFlag{
					Name:  "commands",
					Usage: "provision commands, defaults to true",
				},
				cli.BoolFlag{
					Name:  "parents",
					Usage: "make parent directories as needed, defaults to false",
				},
				cli.BoolFlag{
					Name:  "clobber",
					Usage: "remove existing files/directories before linking, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					cli.ShowCommandHelp(c, "provision")
					os.Exit(1)
				}

				ignoreArgCheck = true

				argMachine := c.Args().First()
				argTasks := c.Args().Tail()

				// export the machine name to the environment
				log.Printf("Exporting MACHINE: %s\n", argMachine)
				os.Setenv("MACHINE", argMachine)

				global.password = c.GlobalString("password")

				flags := &provisionFlags{global: global}
				flags.links = c.BoolT("links")
				flags.commands = c.BoolT("commands")
				flags.parents = c.BoolT("parents")
				flags.clobber = c.BoolT("clobber")

				cmd := app.Command("decrypt")
				cmd.Run(c)

				cfg, err := requireConfig()
				if err != nil {
					log.Fatal(err)
				}

				err = cfg.filterTasks(argMachine, argTasks)
				if err != nil {
					log.Fatal(err)
				}

				err = provisionMachine(argMachine, cfg, flags)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "encrypt files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove original files, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					cli.ShowCommandHelp(c, "encrypt")
					os.Exit(1)
				}

				global.password = c.GlobalString("password")

				flags := &encryptFlags{global: global}
				flags.remove = c.BoolT("remove")

				cfg, err := requireConfig()
				if err != nil {
					log.Fatal(err)
				}

				err = encryptFiles(cfg.Encrypted, flags)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "decrypt files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove encrypted files, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				if !ignoreArgCheck {
					if len(c.Args()) > 0 {
						cli.ShowCommandHelp(c, "decrypt")
						os.Exit(1)
					}
				}

				global.password = c.GlobalString("password")

				flags := &decryptFlags{global: global}
				flags.remove = c.BoolT("remove")

				cfg, err := requireConfig()
				if err != nil {
					log.Fatal(err)
				}

				if len(cfg.Encrypted) > 0 {
					err = decryptFiles(cfg.Encrypted, flags)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Println("No encrypted files are specified in alter.yaml")
				}
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "clean provisioned links",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove encrypted files, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				// TODO: clean this up
				// if len(c.Args()) < 1 {
				// 	cli.ShowCommandHelp(c, "clean")
				// 	os.Exit(1)
				// }
				//
				// cfg, err := requireConfig()
				// if err != nil {
				// 	log.Fatal(err)
				// }
				//
				// machine := c.Args().First()
				//
				// removeTasks := map[string]struct{}{}
				// if len(c.Args()) > 1 {
				// 	for i := 1; i < len(c.Args()); i++ {
				// 		removeTasks[c.Args()[i]] = struct{}{}
				// 	}
				//
				// 	err = cleanMachine(machine, removeTasks, cfg)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// } else {
				// 	err = cleanMachine(machine, nil, cfg)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// }
			},
		},
	}

	app.Run(os.Args)
}
