package main

import "io"

type TerraformRunner struct {
	Exec Executor
	ExecPrefix string
	Module string
	Options *TerraformOptions
	StdoutCallback string
	Writer io.Writer
}

type TerraformOptions struct {}

type Executor interface {
	Execute(command string, args []string, prefix string) error
}

func (cmd *TerraformRunner) Run() error {
	return nil
}

func main()  {
	cmd := TerraformRunner{}
	_ = cmd.Run()
}