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
	if !task.Queued {
		return nil
	}

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

	if cachedMachine, ok := cache.Machines[p.Cfg.Machine]; ok {
		if cachedMachine.SHA1 == p.Cfg.SHA1 {
			p.Logger.Info(0, "Machine up to date: %s", p.Cfg.Machine)

			return nil
		}

		p.Logger.Info(0, "Preparing machine for update: %s", p.Cfg.Machine)

		for _, task := range p.Cfg.Tasks {
			if cachedTask, ok := cachedMachine.Tasks[task.Name]; ok {
				if cachedTask.SHA1 == task.SHA1 {
					task.Queued = false
				} else {
					p.Logger.Info(1, "Task queued for update: %s", task.Name)

					// update the links
					for _, SHA1 := range cachedTask.Links {
						if link, ok := task.Links[SHA1]; ok {
							link.Queued = false
						}
					}

					// update the commands
					for _, SHA1 := range cachedTask.Commands {
						if command, ok := task.Commands[SHA1]; ok {
							command.Queued = false
						}
					}
				}

				// update the cache
				delete(cachedMachine.Tasks, task.Name)
				cache.AddTask(cachedMachine, task)
			} else {
				for cachedTaskName, cachedTask := range cachedMachine.Tasks {
					if cachedTask.SHA1 == task.SHA1 {
						p.Logger.Info(1, "Task renamed: %s -> %s", cachedTaskName, task.Name)

						task.Queued = false

						// update the cache
						delete(cachedMachine.Tasks, task.Name)
						cache.AddTask(cachedMachine, task)

						break
					}
				}
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
