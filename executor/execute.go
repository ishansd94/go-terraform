package executor

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

type DefaultExecute struct {
	Writer       io.Writer
	Dir 		string
}

func (d *DefaultExecute) Execute (command string, args []string, prefix string) error {
	var err error

	cmd := exec.Command(command, args...)
	cmd.Dir	= d.Dir
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