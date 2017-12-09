package main

import (
	"time"
	"os"

	"github.com/cirocosta/auto53/lib"
	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
)

type cliConfig struct {
	Interval time.Duration `arg:"help:interval between state retrieval`
	Dry      bool          `arg:"help:run once without performing modifications"`
}

var (
	args = &cliConfig{
		Interval: 1 * time.Minute,
	}
	logger = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
)

func must (err error) {
	if err == nil {
		return
	}

	logger.Fatal().
		Err(err).
		Msg("main execution failed")
}

func main() {
	arg.MustParse(args)

	_, err := lib.NewAuto(lib.AutoConfig{})
	must(err)
}
