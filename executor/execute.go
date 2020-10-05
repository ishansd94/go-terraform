package executor

import (
	"errors"
	"os"
	"os/exec"
)

func Execute(command string, args []string, prefix string) error {
	var err error

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		return errors.New("(Executor::Execute) -> " + err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		return errors.New("(Executor::Execute) -> " + err.Error())
	}

	return nil
}