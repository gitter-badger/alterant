package main

// TODO: create a `prepare` command that deletes any decrypted files
// TODO: DRY the common app command code
// TODO: add an option to specify which branch the repo should be on
// TODO: files that are encrypted should not be decrypted unless rquired by a task
// TODO: add a `require` field to ensure task dependencies are fulfilled
// TODO: add encryption groups that can be encrypted with different passwords
// TODO: add `update` command that cleans the current environment, pulls the
// updated repo and reprovisions the machine
// TODO: add a `update_strategy` section to the config that can indicate options
// for updating
// TODO: add the option to use an alter.yaml from a git repo.

import (
	"fmt"
	"log"
	"os"

	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/provisioner"
	"github.com/autonomy/alterant/repo"
	"github.com/codegangsta/cli"
)

var version string

func main() {
	var ignoreArgCheck bool

	app := cli.NewApp()

	app.Usage = "Alter your machine with ease."
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "password to encrypt/decrypt a file",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "use verbose logging",
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

				// decrypt files before provisioning
				ignoreArgCheck = true
				// cmd := app.Command("decrypt")
				// cmd.Run(c)

				machine, err := repo.CurrentMachine()
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				requests := c.Args()

				err = repo.OpenMachineByName(machine)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				cfg, err := config.AcquireConfig(machine)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				provisioner := provisioner.NewDefaultProvisioner(machine, cfg, c)

				err = provisioner.Provision(requests)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "clean provisioned links",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					cli.ShowCommandHelp(c, "clean")
					os.Exit(1)
				}

				// machine := c.Args().First()
				// requests := c.Args().Tail()
				//
				// cfg, err := config.AcquireConfig(machine)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				//
				// provisioner := provisioner.NewDefaultProvisioner(machine, cfg, c)
				//
				// err = provisioner.Clean()
				// if err != nil {
				// 	log.Fatal(err)
				// }
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

				cfg, err := config.AcquireConfig("")
				if err != nil {
					log.Fatal(err)
				}

				logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

				if len(cfg.Encrypted) > 0 {
					encrypter := encrypter.NewDefaultEncryption(c.GlobalString("password"),
						c.BoolT("remove"), logger)

					err = encrypter.EncryptFiles(cfg.Encrypted)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Println("No encrypted files are specified in alter.yaml")
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

				cfg, err := config.AcquireConfig("")
				if err != nil {
					log.Fatal(err)
				}

				logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

				if len(cfg.Encrypted) > 0 {
					encrypter := encrypter.NewDefaultEncryption(c.GlobalString("password"),
						c.BoolT("remove"), logger)

					err = encrypter.DecryptFiles(cfg.Encrypted)
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Println("No encrypted files are specified in alter.yaml")
				}
			},
		},
		{
			Name:    "machine",
			Aliases: []string{"m"},
			Usage:   "manage a machine",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: "create a new machine",
					Action: func(c *cli.Context) {
						if len(c.Args()) != 1 {
							cli.ShowSubcommandHelp(c)
							os.Exit(1)
						}
						machine := c.Args().First()

						err := repo.CreateMachine(machine)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					},
				},
				{
					Name:  "open",
					Usage: "open an existing machine",
					Action: func(c *cli.Context) {
						if len(c.Args()) != 1 {
							cli.ShowSubcommandHelp(c)
							os.Exit(1)
						}
						machine := c.Args().First()

						err := repo.OpenMachineByName(machine)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					},
				},
				{
					Name:  "list",
					Usage: "list available machines",
					Action: func(c *cli.Context) {
						if len(c.Args()) != 0 {
							cli.ShowSubcommandHelp(c)
							os.Exit(1)
						}

						err := repo.ListMachines()
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
