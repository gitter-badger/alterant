package main

// TODO: create a `prepare` command that deletes any decrypted files
// TODO: DRY the common app command code
// TODO: files that are encrypted should not be decrypted unless rquired by a task
// TODO: add a `require` field to ensure task dependencies are fulfilled
// TODO: add encryption groups that can be encrypted with different passwords/keys
// TODO: add `update` command that cleans the current environment, pulls the
// updated repo and reprovisions the machine
// TODO: add an `update_strategy` section to the config that can indicate options
// for updating

import (
	"log"
	"os"
	"path"

	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/provisioner"
	"github.com/autonomy/alterant/repo"
	"github.com/codegangsta/cli"
)

var version string

func main() {
	alterantDir := os.Getenv("HOME") + "/.alterant"

	if _, err := os.Stat(alterantDir); os.IsNotExist(err) {
		if err := os.MkdirAll(alterantDir, 0700); err != nil {
			log.Fatal(err)
		}
	}

	app := cli.NewApp()

	app.Usage = "Alter your machine with ease."
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "password to encrypt/decrypt a file",
		},
		cli.StringFlag{
			Name:  "private",
			Value: os.Getenv("HOME") + "/.alterant/secring.gpg",
			Usage: "openpgp ASCII armored private key",
		},
		cli.StringFlag{
			Name:  "public",
			Value: os.Getenv("HOME") + "/.alterant/pubring.gpg",
			Usage: "openpgp ASCII armored public key",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "use verbose logging",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "provision",
			Usage: "provision a machine",
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
				if len(c.Args()) < 2 {
					cli.ShowSubcommandHelp(c)
					os.Exit(1)
				}

				url := c.Args()[0]
				requestedMachine := c.Args()[1]

				err := repo.CloneToAlterantDir(url, requestedMachine, alterantDir)
				if err != nil {
					log.Fatal(err)
				}

				err = os.Chdir(path.Join(alterantDir, requestedMachine))
				if err != nil {
					log.Fatal(err)
				}

				// machine, err := repo.CurrentMachine()
				// if err != nil {
				// 	log.Fatal(err)
				// }

				cfg, err := config.AcquireConfig(requestedMachine)
				if err != nil {
					log.Fatal(err)
				}

				var requests []string
				if len(c.Args()) > 2 {
					requests = c.Args()[2:]
				} else {
					requests = cfg.Order
				}

				provisioner := provisioner.NewDefaultProvisioner(cfg, c)

				err = provisioner.Provision(requests)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "remove",
			Usage: "remove provisioned tasks",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					cli.ShowCommandHelp(c, "remove")
					os.Exit(1)
				}

				machine, err := repo.CurrentMachine()
				if err != nil {
					log.Fatal(err)
				}

				cfg, err := config.AcquireConfig(machine)
				if err != nil {
					log.Fatal(err)
				}

				var requests []string
				if len(c.Args()) == 0 {
					requests = cfg.Order
				} else {
					requests = c.Args()
				}

				provisioner := provisioner.NewDefaultProvisioner(cfg, c)

				err = provisioner.Remove(requests)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
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
					log.Fatal(err)
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

				err := repo.OpenMachine(machine)
				if err != nil {
					log.Fatal(err)
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
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "encrypt",
			Usage: "encrypt files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove unencrypted file, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				machine, err := repo.CurrentMachine()
				if err != nil {
					log.Fatal(err)
				}

				cfg, err := config.AcquireConfig(machine)
				if err != nil {
					log.Fatal(err)
				}

				logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

				encrypter := encrypter.NewDefaultEncryption(c.GlobalString("password"),
					c.GlobalString("private"), c.GlobalString("public"), c.BoolT("remove"), logger)

				err = encrypter.EncryptFiles(cfg)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "decrypt",
			Usage: "decrypt files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "remove",
					Usage: "remove encrypted file, defaults to false",
				},
			},
			Action: func(c *cli.Context) {
				machine, err := repo.CurrentMachine()
				if err != nil {
					log.Fatal(err)
				}

				cfg, err := config.AcquireConfig(machine)
				if err != nil {
					log.Fatal(err)
				}

				logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

				encrypter := encrypter.NewDefaultEncryption(c.GlobalString("password"),
					c.GlobalString("private"), c.GlobalString("public"), c.BoolT("remove"), logger)

				err = encrypter.DecryptFiles(cfg)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:        "gen-key",
			Usage:       "generate a private/public key pair",
			Description: "Takes a name, comment, and email address as arguments.",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 3 {
					cli.ShowCommandHelp(c, "gen-key")
					os.Exit(1)
				}

				name := c.Args()[0]
				comment := c.Args()[1]
				email := c.Args()[2]

				err := encrypter.NewKeyPair(name, comment, email)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
