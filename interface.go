package powershell

import "io"

type Shell interface {
	Execute(cmd string) (string, string, error)
	Close() error
}

type Middleware Shell

type SSHSession interface {
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
	StderrPipe() (io.Reader, error)
	Start(string) error
	Wait() error
}
