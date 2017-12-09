package lib

import (
	"os"

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
