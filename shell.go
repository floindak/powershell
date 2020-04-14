package powershell

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	winodwsISEBinary     = "powershell"
	powershellCoreBinary = "pwsh"
)

var (
	binary = winodwsISEBinary
)

func SwitchPowerShellVersion(isPowershellCore bool) {
	if isPowershellCore {
		binary = powershellCoreBinary
	} else {
		binary = winodwsISEBinary
	}
}

type shell struct {
	waiter func() error
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
	sync.Mutex
}

func Local() (Shell, error) {
	startCMD := []string{binary, "-NoExit", "-Command", "-"}
	command := exec.Command(startCMD[0], startCMD[1:]...)
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stdin stream")
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stdout stream")
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stderr stream")
	}

	err = command.Start()
	return &shell{
		waiter: command.Wait,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, err
}

func SSH(sess SSHSession) (*shell, error) {
	startCMD := []string{binary, "-NoExit", "-Command", "-"}
	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stdin stream")
	}

	stdout, err := sess.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stdout stream")
	}

	stderr, err := sess.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Could not get hold of the PowerShell's stderr stream")
	}

	err = sess.Start(strings.Join(startCMD, " "))
	return &shell{
		waiter: sess.Wait,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, err
}

func (s *shell) Execute(cmd string) (string, string, error) {
	if s.waiter == nil {
		return "", "", errors.Wrap(errors.New(cmd), "Cannot execute commands on closed shells.")
	}
	s.Lock()
	defer s.Unlock()

	outBoundary := createBoundary()
	errBoundary := createBoundary()

	// wrap the command in special markers so we know when to stop reading from the pipes
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, outBoundary, errBoundary, "\n")

	_, err := s.stdin.Write([]byte(full))
	if err != nil {
		return "", "", errors.Wrap(errors.Wrap(err, cmd), "Could not send PowerShell command")
	}

	// read stdout and stderr
	sout := ""
	serr := ""

	waiter := &sync.WaitGroup{}
	waiter.Add(2)

	go streamReader(s.stdout, outBoundary, &sout, waiter)
	go streamReader(s.stderr, errBoundary, &serr, waiter)

	waiter.Wait()

	if len(serr) > 0 {
		return sout, serr, errors.Wrap(errors.New(cmd), serr)
	}

	return sout, serr, nil
}

func (s *shell) Close() error {
	defer func() {
		s.waiter = nil
		s.stdin = nil
		s.stdout = nil
		s.stderr = nil
	}()
	s.stdin.Write([]byte("exit\n"))

	// if it's possible to close stdin, do so (some backends, like the local one,
	// do support it)
	closer, ok := s.stdin.(io.Closer)
	if ok {
		return closer.Close()
	}
	return s.waiter()
}

// read all output until we have found our boundary token
func streamReader(stream io.Reader, boundary string, buffer *string, signal *sync.WaitGroup) {
	defer signal.Done()
	output := ""
	bufsize := 1024

	for {
		buf := make([]byte, bufsize)
		read, err := stream.Read(buf)
		if err != nil {
			return
		}

		output = output + string(buf[:read])
		if index := strings.LastIndex(output, boundary); index >= 0 {
			*buffer = output[:index]
			return
		}
	}
}
