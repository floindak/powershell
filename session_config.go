// Copyright (c) 2017 Gorillalabs. All rights reserved.

package powershell

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

const (
	HTTPPort  = 5985
	HTTPSPort = 5986
)

type SessionConfig struct {
	ComputerName          string
	AllowRedirection      bool
	Authentication        string
	CertificateThumbprint string
	Credential            interface{}
	Port                  int
	UseSSL                bool
}

func (c *SessionConfig) ToArgs() []string {
	args := make([]string, 0)

	if c.ComputerName != "" {
		args = append(args, "-ComputerName")
		args = append(args, quoteArg(c.ComputerName))
	}

	if c.AllowRedirection {
		args = append(args, "-AllowRedirection")
	}

	if c.Authentication != "" {
		args = append(args, "-Authentication")
		args = append(args, quoteArg(c.Authentication))
	}

	if c.CertificateThumbprint != "" {
		args = append(args, "-CertificateThumbprint")
		args = append(args, quoteArg(c.CertificateThumbprint))
	}

	if c.Port > 0 {
		args = append(args, "-Port")
		args = append(args, strconv.Itoa(c.Port))
	}

	if asserted, ok := c.Credential.(string); ok {
		args = append(args, "-Credential")
		args = append(args, asserted) // do not quote, as it contains a variable name when using password auth
	}

	if c.UseSSL {
		args = append(args, "-UseSSL")
	}

	return args
}

type credential interface {
	prepare(Middleware) (interface{}, error)
}

type UserPasswordCredential struct {
	Username string
	Password string
}

func (c *UserPasswordCredential) prepare(s Middleware) (interface{}, error) {
	name := "goCred" + randomString(8)
	pwname := "goPass" + randomString(8)

	_, _, err := s.Execute(fmt.Sprintf("$%s = ConvertTo-SecureString -String %s -AsPlainText -Force", pwname, quoteArg(c.Password)))
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert password to secure string")
	}

	_, _, err = s.Execute(fmt.Sprintf("$%s = New-Object -TypeName 'System.Management.Automation.PSCredential' -ArgumentList %s, $%s", name, quoteArg(c.Username), pwname))
	if err != nil {
		return nil, errors.Wrap(err, "Could not create PSCredential object")
	}

	return fmt.Sprintf("$%s", name), nil
}
