package lib

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func FormattingRulesFromYamlFile(file string) (rules []*FormattingRule, err error) {
	finfo, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.Wrapf(err,
				"configuration file %s not found",
				file)
			return
		}

		err = errors.Wrapf(err,
			"unexpected error looking for config file %s",
			file)
		return
	}

	configContent, err := ioutil.ReadAll(finfo)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly read config file %s",
			file)
		return
	}

	err = yaml.Unmarshal(configContent, &rules)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't properly parse yaml config file %s",
			file)
		return
	}

	return
}
