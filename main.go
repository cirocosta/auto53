package main

import (
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/auto53/lib"
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
	rules = []*lib.FormattingRule{
		{
			AutoScalingGroup: "wedeploy-swarm-worker-xyz1",
			Zone:             "private-wedeploy-xyz1",
			Record:           "{{ .Id }}-asg",
		},
		{
			AutoScalingGroup: "wedeploy-swarm-worker-xyz1",
			Zone:             "private-wedeploy-xyz1",
			Record:           "instances",
		},
	}
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
