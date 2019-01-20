package exec

import (
	"devlab/lib/logger"
	"os"
	"os/exec"
)

// GitCommand executes git <command> in shell and returns result
func GitCommand(serviceDir string, command string) (result string, err error) {
	out, err := exec.Command(
		"sh",
		"-c",
		"( cd "+serviceDir+" && "+command+")").Output()

	return string(out), err
}

// Command executes <command> in shell and returns result
func Command(command string) (result string, err error) {
	out, err := exec.Command(
		"sh",
		"-c",
		command).Output()

	return string(out), err
}

// CommandToStdout executes <command> in shell and prints out flow to stdout
func CommandToStdout(command string) (err error) {
	c := exec.Command(
		"sh",
		"-c",
		command)

	c.Stdout = os.Stdout
	err = c.Run()

	return
}

// CommandToStdoutDebug executes <command> in shell and prints out flow to stdout
func CommandToStdoutDebug() (err error) {
	command := `/data/projects/go/src/devlab/contexts/test-context/build-service`
	logger.Debug("command: ", command)
	c := exec.Command(
		"sh",
		"-c",
		command)

	c.Stdout = os.Stdout
	err = c.Run()

	return
}
