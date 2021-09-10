package config

import (
	"fmt"
	"net/url"
)

type EasyvereinCfg struct {
	Host    string            `yaml:"host"`
	Path    string            `yaml:"path"`
	Token   string            `yaml:"token"`
	Options map[string]string `yaml:"options,omitempty"`
}

// optionsURI will concatenate all options specified in config.yaml
// and add them to the url
func (c *EasyvereinCfg) optionsURI() string {
	optionsURI := url.Values{}

	for key, value := range c.Options {
		optionsURI.Add(key, value)
	}

	return optionsURI.Encode()
}

// APIRequestURI parses the URI for the API Request to easyverein
func (c *EasyvereinCfg) APIRequestURI(endpoint string) string {
	path := fmt.Sprintf("%s%s",
		config.Easyverein.Path,
		endpoint,
	)
	requestURI := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     path,
		RawQuery: c.optionsURI(),
	}

	return requestURI.String()
}
