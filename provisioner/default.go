package provisioner

import (
	"fmt"

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
		fmt.Println(err)
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
func (p *DefaultProvisioner) Provision(requests []*task.Task) error {
	p.Logger.Info(0, "Provisioning: %s", p.Cfg.Machine)

	// decrypt files
	err := p.Encrypter.DecryptFiles(p.Cfg)
	if err != nil {
		return err
	}

	for _, task := range requests {
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

// TODO: execute only the components that have been changed e.g. if a command
// has been updated but everything else is the same, then only execute that
// command
// Update updates a machine's tasks
func (p *DefaultProvisioner) Update(cfg *config.Config) error {
	// db, err := bolt.Open("/home/vagrant/.alterant/cache.db", 0600, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// db.Update(func(tx *bolt.Tx) error {
	// 	mb := tx.Bucket([]byte(cfg.Machine))
	// 	if mb == nil {
	// 		fmt.Printf("Machine not found: %s\n", cfg.Machine)
	// 	}

	// 	for _, task := range cfg.Order {
	// 		taskAdded := false
	// 		taskRenamed := false

	// 		tb := mb.Bucket([]byte("tasks"))
	// 		if tb == nil {
	// 			fmt.Println("Tasks not found in cache")

	// 			return nil
	// 		}

	// 		t := tb.Bucket([]byte(task.Name))

	// 		if t == nil {
	// 			fmt.Printf("Task not found in cache: %s\n", task.Name)

	// 			// check if the task has been renamed
	// 			c := tb.Cursor()

	// 			for k, v := c.First(); k != nil; k, v = c.Next() {
	// 				t = tb.Bucket([]byte(k))

	// 				sha1 := t.Get([]byte("SHA1"))

	// 				if task.SHA1 == string(sha1) {
	// 					taskRenamed = true
	// 					fmt.Printf("Task renamed: %s -> %s", k, task.Name)

	// 					tb.Delete([]byte(k))

	// 					err = p.cacheTask(tb, task)
	// 					if err != nil {
	// 						return err
	// 					}

	// 					break
	// 				}

	// 				// TODO: I don't want to have to use `v`
	// 				fmt.Printf("value=%s", v)
	// 			}

	// 			// if the task is not a rename then it must be a new task
	// 			if !taskRenamed {
	// 				fmt.Printf("Adding new task: %s\n", task.Name)

	// 				err = p.executeTask(task)
	// 				if err != nil {
	// 					return err
	// 				}

	// 				err = p.cacheTask(tb, task)
	// 				if err != nil {
	// 					return err
	// 				}

	// 				taskAdded = true
	// 			}
	// 		}

	// 		// if the task is not new, then it possibly needs updating
	// 		if !taskAdded {
	// 			sha1 := t.Get([]byte("SHA1"))

	// 			if task.SHA1 == string(sha1) {
	// 				fmt.Printf("Task is clean: %s\n", task.Name)
	// 				continue
	// 			} else {
	// 				fmt.Printf("Updating task: %s\n", task.Name)

	// 				// err = tb.DeleteBucket([]byte(task.Name))
	// 				// if err != nil {
	// 				// 	return err
	// 				// }

	// 				lb := t.Bucket([]byte("links"))
	// 				if lb == nil {
	// 					fmt.Println("Links not found")
	// 					return fmt.Errorf("Links not found")
	// 				}

	// 				var links []*link.Link
	// 				for _, link := range task.Links {
	// 					l := lb.Get([]byte(link.SHA1))
	// 					if l == nil {
	// 						fmt.Println("Link exists")
	// 						links = append(links, link)
	// 					}
	// 				}
	// 				task.Links = links

	// 				cb := t.Bucket([]byte("commands"))
	// 				if cb == nil {
	// 					fmt.Println("Commands not found")
	// 					return fmt.Errorf("Commands not found")
	// 				}

	// 				var commands []*command.Command
	// 				for _, command := range task.Commands {
	// 					c := cb.Get([]byte(command.SHA1))
	// 					if c == nil {
	// 						commands = append(commands, command)
	// 					}
	// 				}

	// 				task.Commands = commands
	// 				err = p.executeTask(task)
	// 				if err != nil {
	// 					return err
	// 				}

	// 				err = p.cacheTask(tb, task)
	// 				if err != nil {
	// 					return err
	// 				}

	// 			}
	// 		}
	// 	}

	// 	return nil
	// })

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
