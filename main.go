package main

// TODO: create a `prepare` command that deletes any decrypted files
// TODO: DRY the common app command code
import (
	"log"
	"os"

	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/provisioner"
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
				cmd := app.Command("decrypt")
				cmd.Run(c)

				argMachine := c.Args().First()
				argTasks := c.Args().Tail()

				cfg, err := config.AcquireConfig(argMachine)
				if err != nil {
					log.Fatal(err)
				}

				err = cfg.FilterTasks(argMachine, argTasks)
				if err != nil {
					log.Fatal(err)
				}

				provisioner := provisioner.NewDefaultProvisioner(argMachine, cfg, c)

				err = provisioner.Provision()
				if err != nil {
					log.Fatal(err)
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

				argMachine := c.Args().First()
				argTasks := c.Args().Tail()

				cfg, err := config.AcquireConfig(argMachine)
				if err != nil {
					log.Fatal(err)
				}

				err = cfg.FilterTasks(argMachine, argTasks)
				if err != nil {
					log.Fatal(err)
				}

				provisioner := provisioner.NewDefaultProvisioner(argMachine, cfg, c)

				err = provisioner.Clean()
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
	}

	app.Run(os.Args)
}
