package alert

import (
	"os"
	"os/exec"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestSendToast(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	err := SendToast("linux", "2080")
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}

	err = SendToast("darwin", "2080")
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
}
