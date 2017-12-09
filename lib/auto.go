package lib

import (
	"os"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/route53"
	"github.com/rs/zerolog"
)

type Auto struct {
	logger zerolog.Logger
}

type AutoConfig struct{}

func NewAuto(cfg AutoConfig) (a Auto, err error) {
	a.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "auto").
		Logger()

	return
}
