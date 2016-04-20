package provisioner

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/cache"
	"github.com/autonomy/alterant/commander"
	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/encrypter"
	"github.com/autonomy/alterant/environment"
	"github.com/autonomy/alterant/linker"
	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/task"
	"github.com/boltdb/bolt"
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

func (p *DefaultProvisioner) cacheTask(tb *bolt.Bucket, task *task.Task) error {
	// save the task to cache.db
	p.Logger.Info(1, "Caching task: %s", task.Name)

	t, err := tb.CreateBucketIfNotExists([]byte(task.Name))
	if err != nil {
		return err
	}

	t.Put([]byte("SHA1"), []byte(task.SHA1))

	for _, dep := range task.Dependencies {
		err = t.Put([]byte(cache.SHAFromString(dep)), []byte(dep))
		if err != nil {
			return err
		}
	}

	for _, link := range task.Links {
		l, err := yaml.Marshal(&link)
		if err != nil {
			return err
		}

		err = t.Put([]byte(link.SHA1), l)
		if err != nil {
			return err
		}
	}

	for _, command := range task.Commands {
		c, err := yaml.Marshal(&command)
		if err != nil {
			return err
		}

		err = t.Put([]byte(command.SHA1), c)
		if err != nil {
			return err
		}
	}

	return nil
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
func (p *DefaultProvisioner) Provision(requests []*task.Task) error {
	db, err := bolt.Open("/home/vagrant/.alterant/cache.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	p.Logger.Info(0, "Provisioning: %s", p.Cfg.Machine)

	db.Update(func(tx *bolt.Tx) error {
		mb, err := tx.CreateBucketIfNotExists([]byte(p.Cfg.Machine))
		if err != nil {
			return err
		}

		mb.Put([]byte("SHA1"), []byte(p.Cfg.Sha1))

		tb, err := mb.CreateBucketIfNotExists([]byte("tasks"))
		if err != nil {
			return err
		}

		// decrypt files
		err = p.Encrypter.DecryptFiles(p.Cfg)
		if err != nil {
			return err
		}

		for _, task := range requests {
			err = p.executeTask(task)
			if err != nil {
				return err
			}

			err = p.cacheTask(tb, task)
			if err != nil {
				return err
			}
		}

		return nil
	})

	p.Logger.Info(0, "Provisioned: %s", p.Cfg.Machine)

	return nil
}

// Update updates a machine's tasks
func (p *DefaultProvisioner) Update(cfg *config.Config) error {
	db, err := bolt.Open("/home/vagrant/.alterant/cache.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(cfg.Machine))
		if mb == nil {
			fmt.Printf("Machine not found: %s\n", cfg.Machine)
		}

		for _, task := range cfg.Order {
			tb := mb.Bucket([]byte("tasks"))
			if tb == nil {
				fmt.Println("Tasks not found in cache")

				return nil
			}

			t := tb.Bucket([]byte(task.Name))
			if t == nil {
				fmt.Printf("Task not found in cache: %s\n", task.Name)

				c := tb.Cursor()

				// check if the task has been renamed
				for k, v := c.First(); k != nil; k, v = c.Next() {
					t = tb.Bucket([]byte(k))

					sha1 := t.Get([]byte("SHA1"))

					if task.SHA1 == string(sha1) {
						fmt.Printf("Task renamed: %s -> %s", k, task.Name)

						tb.Delete([]byte(k))

						err = p.cacheTask(tb, task)
						if err != nil {
							return err
						}

						break
					}

					fmt.Printf("value=%s", v)
				}

				fmt.Printf("Adding new task: %s\n", task.Name)

				err = p.executeTask(task)
				if err != nil {
					return err
				}

				err = p.cacheTask(tb, task)
				if err != nil {
					return err
				}

				return nil
			}

			sha1 := t.Get([]byte("SHA1"))

			if task.SHA1 == string(sha1) {
				fmt.Printf("Task is clean: %s\n", task.Name)
				continue
			} else {
				fmt.Printf("Updating task: %s\n", task.Name)

				err = tb.DeleteBucket([]byte(task.Name))
				if err != nil {
					return err
				}

				p.executeTask(task)
				p.cacheTask(tb, task)
			}
		}
		return nil
	})

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
