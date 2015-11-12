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

// State represents the state of alterant
type State struct {
	Logger      *logWrapper.LogWrapper
	Environment *environment.Environment
	Encrypter   encrypter.Encrypter
	Linker      linker.Linker
	Commander   commander.Commander
	Cfg         *config.Config
	Machine     string
}

// NewState returns an instance of `State`
func NewState(argMachine string, cfg *config.Config, c *cli.Context) *State {
	logger := logWrapper.NewLogWrapper(c.GlobalBool("verbose"))

	s := &State{
		Logger:      logger,
		Environment: environment.NewEnvironment(argMachine, logger),
		Encrypter: encrypter.NewDefaultEncryption(c.GlobalString("password"),
			c.BoolT("remove"), logger),
		Linker: linker.NewDefaultLinker(c.BoolT("links"), c.Bool("parents"),
			c.Bool("clobber"), logger),
		Commander: commander.NewDefaultCommander(c.BoolT("commands"), logger),
		Cfg:       cfg,
		Machine:   argMachine,
	}
	return s
}
