package powershell

import (
	"fmt"
	"strings"
	"testing"
)

func TestLocal(t *testing.T) {
	shell, err := Local()
	if err != nil {
		t.Fatal(err)
	}
	defer shell.Close()
	testShell(t, shell)
}

func testShell(t *testing.T, shell Shell){
	stdout, stderr, err := shell.Execute("echo test powershell wrapper")
	if err != nil {
		t.Fatal(err)
	}
	if stderr != "" {
		t.Fatal(stderr)
	}
	stdout=strings.Join(strings.Fields(stdout)," ")
	expected := fmt.Sprintf("test%[1]spowershell%[1]swrapper", " ")
	if stdout != expected {
		t.Fatal("\nexpected:", expected, "\nactually:", stdout)
	}
}
