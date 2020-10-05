package main

import (
	"fmt"
	"io"

	"github.com/ishansd94/terraform-go/executor"
)

const (
	TerrformBin = "terraform"
)


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

	err := executor.Execute(TerrformBin, []string{"init"}, "" )

	fmt.Println(err)

	return nil
}

func main()  {
	cmd := TerraformRunner{}
	_ = cmd.Run()
}