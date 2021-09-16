package config

import (
	"fmt"
	"net/url"
)

type WordpressCfg struct {
	Host     string            `yaml:"host"`
	Path     string            `yaml:"path"`
	Username string            `yaml:"user"`
	Password string            `yaml:"pass"`
	Options  map[string]string `yaml:"options,omitempty"`
}

// OptionsURI will concatenate all options specified in config.yaml
// and add them to the url
func (c *WordpressCfg) optionsURI() string {
	optionsURI := url.Values{}

	for key, value := range c.Options {
		optionsURI.Add(key, value)
	}

	return optionsURI.Encode()
}

// APIRequestURI parses the URI for the API Request to wordpress
func (c *WordpressCfg) APIRequestURI(endpoint string, page int) string {
	path := fmt.Sprintf("%s%s",
		config.Wordpress.Path,
		endpoint,
	)
	requestURI := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     path,
		RawQuery: fmt.Sprintf("%s&page=%d", c.optionsURI(), page),
	}

	return requestURI.String()
}