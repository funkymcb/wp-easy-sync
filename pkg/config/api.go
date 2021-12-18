package config

import (
	"fmt"
	"time"

	"github.com/savsgio/atreugo/v11"
)

type APICfg struct {
	Port        int           `default:"8080" yaml:"port"`
	Token       string        `yaml:"token"`
	AuthOffsete time.Duration `yaml:"auth-offset"`
}

func (config *APICfg) AtreugoConfig() atreugo.Config {
	return atreugo.Config{
		Addr: fmt.Sprintf(":%d", config.Port),
	}
}
