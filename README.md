# Powershell

This package is inspired by [jPowerShell](https://github.com/profesorfalken/jPowerShell)
and allows one to run and remote-control a PowerShell session.

## Installation

    go get github.com/floindak/powershell

## Usage

### PowerShell ISE
It uses `os/exec` to start a PowerShell process to interactive in default.

```go
package main

import (
    "fmt"
    "github.com/floindak/powershell"
)

func main() {
    // windows use local powershell
    shell,err:=powershell.Local()
    panicIF(err)
    defer shell.Close()
    
    // execute some command
    stdout, stderr, err := shell.Execute("echo foo bar")
    panicIF(err)
    fmt.Println(stdout,stderr)
}

func panicIF(e error){
    if e!=nil{
        panic(e)
    }
}
```

### PowerShell Core
You can also interact with PowerShell Core, which normally have a binary name `pwsh`, use `powershell.SwitchPowerShellVersion(true)` before opening a new shell.

### PowerShell Via SSH
```go
package main

import (
    "fmt"
    "github.com/floindak/powershell"
    "golang.org/x/crypto/ssh"
)

func main(){
    // open powershell core via ssh
    powershell.SwitchPowerShellVersion(true)
    sshClient,err:=ssh.Dial("tcp", "__host__:22", &ssh.ClientConfig{
        User: "__user__",
        Auth: []ssh.AuthMethod{ssh.Password("__password__")},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    })
    panicIF(err)
    sess,err:=sshClient.NewSession()
    panicIF(err)
    shell,err:=powershell.SSH(sess)
    panicIF(err)
    defer shell.Close()

    // execute some command
    stdout, stderr, err := shell.Execute("echo foo bar")
    panicIF(err)
    fmt.Println(stdout,stderr)

}
func panicIF(e error){
    if e!=nil{
        panic(e)
    }
}
```

### Remote Sessions

You can use an existing PS shell to use PSSession cmdlets to connect to remote
computers. Instead of manually handling that, you can use the Session middleware,
which takes care of authentication. Note that you can still use the "raw" shell
to execute commands on the computer where the powershell host process is running.

```go
package main

import (
    "fmt"
    "github.com/floindak/powershell"
)

func main() {
    // use default backend start a local powershell process
    shell, err := powershell.Local()
    panicIF(err)
    defer shell.Close()

    // create a new shell by wrapping the existing one in the session middleware
    session, err := powershell.NewSession(shell, &powershell.SessionConfig{
        ComputerName: "remote-pc-1",
    })
    panicIF(err)
    defer session.Close() // will also close the underlying ps shell!

    // everything run via the session is run on the remote machine
    stdout, stderr, err := session.Execute("echo foo bar")
    panicIF(err)
    fmt.Println(stdout, stderr)
}
func panicIF(e error) {
    if e != nil {
        panic(e)
    }
}
```

## License
[MIT](LICENSE)
