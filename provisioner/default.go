package provisioner

import (
	"github.com/autonomy/alterant/cache"
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

func (p *DefaultProvisioner) executeTask(task *task.Task) error {
	p.Logger.Info(1, "Attempting task: %s", task.Name)

	// export environment variables specific to the specified machine
	p.Environment.Set(p.Cfg.Environment)

	// create the links specified in the task
	err := p.Linker.CreateLinks(task.Links)
	if err != nil {
		return err
	}

	// execute the commands specified in the task
	err = p.Commander.Execute(task)
	if err != nil {
		return err
	}

	p.Logger.Info(1, "Task fulfilled: %s", task.Name)

	return nil
}

// Provision provisions a machine
func (p *DefaultProvisioner) Provision() error {
	p.Logger.Info(0, "Provisioning: %s", p.Cfg.Machine)

	// decrypt files
	err := p.Encrypter.DecryptFiles(p.Cfg)
	if err != nil {
		return err
	}

	for _, task := range p.Cfg.Tasks {
		err = p.executeTask(task)
		if err != nil {
			return err
		}
	}

	p.Logger.Info(0, "Provisioned: %s", p.Cfg.Machine)

	p.Logger.Info(0, "Caching: %s", p.Cfg.Machine)

	err = cache.WriteToFile(p.Cfg)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a machine's tasks
func (p *DefaultProvisioner) Update() error {
	cache, err := cache.ReadCache()
	if err != nil {
		return err
	}

	if machine, ok := cache.Machines[p.Cfg.Machine]; ok {
		if machine.SHA1 == p.Cfg.SHA1 {
			p.Logger.Info(1, "Machine is clean: %s", p.Cfg.Machine)

			return nil
		}

		p.Logger.Info(1, "Changes detected for machine: %s", p.Cfg.Machine)

		for i := len(p.Cfg.Tasks) - 1; i >= 0; i-- {
			task := p.Cfg.Tasks[i]
			if t, ok := machine.Tasks[task.Name]; ok {
				if t.SHA1 == task.SHA1 {
					p.Logger.Info(2, "Task is clean: %s", task.Name)

					p.Cfg.Tasks = append(p.Cfg.Tasks[:i], p.Cfg.Tasks[i+1:]...)
				} else {
					p.Logger.Info(2, "Task is dirty: %s", task.Name)
					for _, l := range t.Links {
						for i := len(task.Links) - 1; i >= 0; i-- {
							if l == task.Links[i].SHA1 {
								p.Logger.Info(2, "Link is clean: %s", l)

								task.Links = append(task.Links[:i], task.Links[i+1:]...)
							} else {
								p.Logger.Info(2, "Link has changed: %s", l)
							}
						}
					}

					for _, c := range t.Commands {
						for i := len(task.Commands) - 1; i >= 0; i-- {
							if c == task.Commands[i].SHA1 {
								p.Logger.Info(2, "Command is clean: %s", c)

								task.Commands = append(task.Commands[:i], task.Commands[i+1:]...)
							} else {
								p.Logger.Info(2, "Command has changed: %s", c)
							}
						}
					}
				}
			} else {
				p.Logger.Info(2, "New task detected: %s", task.Name)
				// TODO: check for rename
			}
		}
	}

	p.Provision()

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
