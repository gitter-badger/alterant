package provisioner

import (
	"github.com/autonomy/alterant/commander"
	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/environment"
	"github.com/autonomy/alterant/linker"
	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/task"
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
}

// Provision provisions a machine
func (p *DefaultProvisioner) Provision(requests []*task.Task) error {
	p.Logger.Info(0, "Provisioning: %s", p.Cfg.Machine)

	// decrypt files
	err := p.Encrypter.DecryptFiles(p.Cfg)
	if err != nil {
		return err
	}

	for _, task := range requests {
		p.Logger.Info(1, "Attempting task: %s", task.Name)

		// export environment variables specific to the specified machine
		p.Environment.Set(p.Cfg.Environment)

		// create the links specified in the task
		err = p.Linker.CreateLinks(task.Links)
		if err != nil {
			return err
		}

		// execute the commands specified in the task
		err = p.Commander.Execute(task)
		if err != nil {
			return err
		}

		p.Logger.Info(1, "Task fulfilled: %s", task.Name)
	}

	p.Logger.Info(0, "Provisioned: %s", p.Cfg.Machine)

	return nil
}

// Remove removes provisioned tasks
func (p *DefaultProvisioner) Remove(requests []string) error {
	p.Logger.Info(2, "Removing: %s", p.Cfg.Machine)
	for _, task := range p.Cfg.Tasks {
		p.Linker.RemoveLinks(task.Links)
	}

	return nil
}

// NewDefaultProvisioner returns an instance of a `DefaultProvisioner`
func NewDefaultProvisioner(cfg *config.Config, c *cli.Context) *DefaultProvisioner {
	logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

	p := &DefaultProvisioner{
		Logger:      logger,
		Environment: environment.NewEnvironment(cfg.Machine, logger),
		Encrypter: encrypter.NewDefaultEncryption(c.GlobalString("password"),
			c.GlobalString("private"), c.GlobalString("public"), c.BoolT("remove"), logger),
		Linker: linker.NewDefaultLinker(c.BoolT("links"), c.Bool("parents"),
			c.Bool("clobber"), logger),
		Commander: commander.NewDefaultCommander(c.BoolT("commands"), logger),
		Cfg:       cfg,
	}

	return p
}
