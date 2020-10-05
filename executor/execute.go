package executor

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	_ "github.com/sanity-io/litter"
)

type DefaultExecute struct {
	Writer      io.Writer
	Dir 		string
}

func (d *DefaultExecute) Execute (command string, args []string, prefix string) (*[]byte, error) {
	var err error
	var output []byte

	if d.Writer == nil  {
		d.Writer = os.Stdout
	}

	cmd := exec.Command(command, args...)
	cmd.Dir	= d.Dir
	cmd.Stderr = os.Stderr
	// cmd.Stdout = d.Writer

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New("(Executor::Execute) -> " + err.Error())
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.New("(Executor::Execute) -> " + err.Error())
	}

	output, err = ioutil.ReadAll(stdout)
	if err != nil {
		return nil, errors.New("(Executor::Execute) -> " + err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		return nil, errors.New("(Executor::Execute) -> " + err.Error())
	}

	return &output, nil
}