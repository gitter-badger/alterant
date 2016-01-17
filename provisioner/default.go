package provisioner

import (
	"github.com/autonomy/alterant/commander"
	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/environment"
	"github.com/autonomy/alterant/linker"
	"github.com/autonomy/alterant/logger"
	"github.com/codegangsta/cli"
)

// DefaultProvisioner is a basic provisioner and is the default
type DefaultProvisioner struct {
	Logger      *logWrapper.LogWrapper
	Environment *environment.Environment
	Encrypter   encrypter.Encrypter
	Linker      linker.Linker
	Commander   commander.Commander
	Cfg         *config.Config
	Machine     string
}

// Provision provisions a machine
func (p *DefaultProvisioner) Provision(requests []string) error {
	p.Logger.Info("Provisioning: %s", p.Machine)

	// iterate over the request array to preserve the order of tasks
	for _, request := range requests {
		task := p.Cfg.Tasks[request]
		if task == nil {
			continue
		}

		p.Logger.Info("Attempting to fulfill request: %s", request)

		// export environment variables specific to the specified machine
		p.Environment.Set(p.Cfg.Environment)

		// create the links specified in the task
		err := p.Linker.CreateLinks(task.Links)
		if err != nil {
			return err
		}

		// execute the commands specified in the task
		// err := p.Commander.Execute(task)
		// if err != nil {
		// 	return err
		// }

		p.Logger.Info("Request fulfilled: %s", request)
	}

	p.Logger.Info("Provisioned: %s", p.Machine)

	return nil
}

// Clean removes provisioned links
func (p *DefaultProvisioner) Clean() error {
	p.Logger.Info("Cleaning: %s", p.Machine)
	// for _, task := range p.Cfg.Tasks {
	// 	p.Linker.RemoveLinks(task.Links)
	// }

	return nil
}

// NewDefaultProvisioner returns and instance of a `DefaultProvisioner`
func NewDefaultProvisioner(machine string, cfg *config.Config, c *cli.Context) *DefaultProvisioner {
	logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

	p := &DefaultProvisioner{
		Logger:      logger,
		Environment: environment.NewEnvironment(machine, logger),
		Encrypter: encrypter.NewDefaultEncryption(c.GlobalString("password"),
			c.BoolT("remove"), logger),
		Linker: linker.NewDefaultLinker(c.BoolT("links"), c.Bool("parents"),
			c.Bool("clobber"), logger),
		Commander: commander.NewDefaultCommander(c.BoolT("commands"), logger),
		Cfg:       cfg,
		Machine:   machine,
	}

	return p
}
