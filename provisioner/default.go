package provisioner

import (
	"github.com/autonomy/alterant/config"
	"github.com/codegangsta/cli"
)

// DefaultProvisioner is a basic provisioner and is the default
type DefaultProvisioner struct {
	state *State
}

// Provision provisions a machine
func (p *DefaultProvisioner) Provision() error {
	p.state.Logger.Info("Provisioning: %s", p.state.Machine)
	//
	// // iterate over the request array to preserve the order of tasks
	for _, request := range p.state.Cfg.Machines[p.state.Machine].Requests {
		task := p.state.Cfg.Tasks[request]
		if task == nil {
			continue
		}

		p.state.Logger.Info("Performing task: %s", request)

		// export environment variables specific to the specified machine
		p.state.Environment.Set(p.state.Machine, p.state.Cfg.Machines[p.state.Machine])

		// create the links specified in the task
		// err := p.state.Linker.CreateLinks(task.Links)
		// if err != nil {
		// 	return err
		// }

		// execute the commands specified in the task
		err := p.state.Commander.Execute(task)
		if err != nil {
			return err
		}

		p.state.Logger.Info("Request fulfilled: %s", request)
	}

	p.state.Logger.Info("Provisioned: %s", p.state.Machine)

	return nil
}

// Clean removes provisioned links
func (p *DefaultProvisioner) Clean() error {
	p.state.Logger.Info("Cleaning: %s", p.state.Machine)
	// for _, task := range p.state.Cfg.Tasks {
	// 	p.state.Linker.RemoveLinks(task.Links)
	// }

	return nil
}

// NewDefaultProvisioner returns and instance of a `DefaultProvisioner`
func NewDefaultProvisioner(argMachine string, cfg *config.Config, c *cli.Context) *DefaultProvisioner {
	state := NewState(argMachine, cfg, c)

	return &DefaultProvisioner{state: state}
}
