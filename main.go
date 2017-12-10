package main

import (
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/auto53/lib"
	"github.com/rs/zerolog"
)

type cliConfig struct {
	Config   string        `arg:"path to the formatting rules configuration file`
	Dry      bool          `arg:"help:run without performing modifications"`
	Interval time.Duration `arg:"help:interval between periodic state retrieval`
	Listen   bool          `arg:"help:listen for API requests`
	Once     bool          `arg:"help:run one time and exit"`
	Port     int           `arg:"help:port to listen for API requests`
}

var (
	args = &cliConfig{
		Config:   "./auto53.yaml",
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
		FormattingRules: rules,
	})
	must(err)

	records, err := a.ListZoneRecords("/hostedzone/Z1UYP7K3ZF7TLR")
	must(err)

	for _, record := range records {
		logger.Info().Interface("record", record).Msg("aaa")
	}

	zones, err := a.ListZones()
	must(err)

	for _, zone := range zones {
		logger.Info().Interface("zone", zone).Msg("bbb")
	}

	asgs, err := a.ListAutoscalingGroups([]string{"wedeploy-swarm-worker-xyz1"})
	must(err)

	for _, asg := range asgs {
		logger.Info().Interface("asg", asg).Msg("ccc")
	}

	records, err = lib.CreateRecords(asgs, rules)
	must(err)

	for _, record := range records {
		logger.Info().Interface("created-record", record).Msg("ddd")
	}

}
