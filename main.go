package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/auto53/lib"
	"github.com/rs/zerolog"
)

type cliConfig struct {
	Config   string        `arg:"help:path to the formatting rules configuration file`
	Debug    bool          `arg:"help:activates debug-level logging"`
	Dry      bool          `arg:"help:run without performing modifications"`
	Interval time.Duration `arg:"help:interval between periodic state retrieval`
	Listen   bool          `arg:"help:listen for API requests`
	Once     bool          `arg:"help:run one time and exit"`
	Port     int           `arg:"help:port to listen for API requests`
}

var (
	args = &cliConfig{
		Config:   "./auto53.yaml",
		Debug:    false,
		Dry:      false,
		Interval: 2 * time.Minute,
		Listen:   false,
		Once:     false,
		Port:     8080,
	}
	logger = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
)

func must(err error) {
	if err == nil {
		return
	}

	logger.Fatal().
		Err(err).
		Msg("main execution failed")
}

func main() {
	arg.MustParse(args)

	rules, err := lib.FormattingRulesFromYamlFile(args.Config)
	must(err)

	a, err := lib.NewAuto(lib.AutoConfig{
		Debug:           args.Debug,
		FormattingRules: rules,
	})
	must(err)

	asgs, err := a.GetAutoScalingGroups()
	must(err)

	currentRecords := []*lib.Record{}

	zonesRecords, err := a.GetZonesRecords()
	must(err)

	for _, records := range zonesRecords {
		currentRecords = append(currentRecords, records...)
	}

	desiredRecords, err := lib.CreateRecords(asgs, rules)
	must(err)

	evals, err := lib.GetEvaluations(currentRecords, desiredRecords)
	must(err)

	if args.Dry {
		fmt.Println("")
		lib.ShowAutoScalingGroupsTable(asgs)

		fmt.Println("")
		lib.ShowEvalsTable(evals)
		return
	}

	err = a.ExecuteEvaluations(evals)
	must(err)

}
