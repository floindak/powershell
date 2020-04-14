// Copyright (c) 2017 Gorillalabs. All rights reserved.

package powershell

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type session struct {
	upstream Middleware
	name     string
}

func NewSession(upstream Middleware, config *SessionConfig) (Middleware, error) {
	asserted, ok := config.Credential.(credential)
	if ok {
		credentialParamValue, err := asserted.prepare(upstream)
		if err != nil {
			return nil, errors.Wrap(err, "Could not setup credentials")
		}

		config.Credential = credentialParamValue
	}

	name := "goSess" + randomString(8)
	args := strings.Join(config.ToArgs(), " ")

	_, _, err := upstream.Execute(fmt.Sprintf("$%s = New-PSSession %s", name, args))
	if err != nil {
		return nil, errors.Wrap(err, "Could not create new PSSession")
	}

	return &session{upstream, name}, nil
}

func (s *session) Execute(cmd string) (string, string, error) {
	return s.upstream.Execute(fmt.Sprintf("Invoke-Command -Session $%s -Script {%s}", s.name, cmd))
}

func (s *session) Close() error {
	s.upstream.Execute(fmt.Sprintf("Disconnect-PSSession -Session $%s", s.name))
	return s.upstream.Close()
}
