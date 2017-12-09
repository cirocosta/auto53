package lib

import (
	"os"

	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Auto struct {
	logger      zerolog.Logger
	route53     *route53.Route53
	autoscaling *autoscaling.AutoScaling
}

type AutoConfig struct{}

func NewAuto(cfg AutoConfig) (a Auto, err error) {
	a.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "auto").
		Logger()

	sess, err := session.NewSession()
	if err != nil {
		err = errors.Wrapf(err, "failed to create aws session")
		return
	}

	a.route53 = route53.New(sess)
	a.autoscaling = autoscaling.New(sess)

	return
}

func (a *Auto) ListAsgInstances(asg string) (instances []*Instance, err error) {
	return
}
