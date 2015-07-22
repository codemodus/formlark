package main

import (
	"github.com/codemodus/config"
)

type conf struct {
	*config.Config
	ServerDomain   string
	ServerPort     string
	FormPathPrefix string
	SMTPDomain     string
	SMTPUser       string
	SMTPPassword   string
	ValidDomains   []string
}

func (c *conf) InitPost() (err error) {
	if c.ServerPort == "" {
		c.ServerPort = ":54541"
	}
	return nil
}
