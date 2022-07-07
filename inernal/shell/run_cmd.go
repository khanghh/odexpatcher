package shell

import (
	"os/exec"
	"strings"
	"time"
)

const defaultShellTimeout = 10 * time.Second

func RunCmd(args ...string) (string, error) {
	cmd := exec.Command("sh", "-c", strings.Join(args, " "))
	timer := time.AfterFunc(defaultShellTimeout, func() {
		cmd.Process.Kill()
	})
	defer timer.Stop()
	output, err := cmd.CombinedOutput()
	return string(output), err
}
