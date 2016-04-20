package provisioner

import (
	"encoding/base64"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"crypto/sha1"

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

type DBEntryType int

const (
	DEPENDENCY DBEntryType = iota
	LINK
	COMMAND
)

type DBEntry struct {
	sha       string
	raw       string
	entryType DBEntryType
}

func NewDBEntry(name string, e []byte, entryType DBEntryType) DBEntry {
	db, err := bolt.Open("/home/vagrant/.alterant/cache.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := sha1.New()
	h.Write([]byte(string(e)))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		fmt.Printf("%#v", b)
		return nil
	})

	return DBEntry{sha: sha, raw: string(e), entryType: entryType}
}

func hashTask(t *task.Task) error {
	var entries []DBEntry

	for _, dep := range t.Dependencies {
		d, err := yaml.Marshal(&dep)
		if err != nil {
			return err
		}

		entries = append(entries, NewDBEntry(t.Name, d, DEPENDENCY))
	}

	for _, link := range t.Links {
		l, err := yaml.Marshal(&link)
		if err != nil {
			return err
		}

		entries = append(entries, NewDBEntry(t.Name, l, LINK))
	}

	for _, command := range t.Commands {
		c, err := yaml.Marshal(&command)
		if err != nil {
			return err
		}

		entries = append(entries, NewDBEntry(t.Name, c, COMMAND))
	}

	fmt.Printf("%#v", entries)

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

		hashTask(task)

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
