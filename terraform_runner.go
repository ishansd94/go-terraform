package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/ishansd94/terraform-go/executor"
)

const (
	TerrformBin     = "terraform"
	FlagAutoApprove = "-auto-approve"
)

type Executor interface {
	Execute(command string, args []string, prefix string) error
}

type TerraformRunner struct {
	Module    string
	Dir 	  string
	Version   string
	Operation string
	Options   *TerraformOptions
	Writer    io.Writer
	Executor  Executor
}

type TerraformOptions struct {
	BackendConfig map[string]string
	Vars          map[string]interface{}
}

func (cmd *TerraformRunner) Run() error {

	var err error

	if cmd.Executor == nil {
		cmd.Executor = &executor.DefaultExecute{
			Dir: cmd.Dir,
			Writer: cmd.Writer,
		}
	}

	_, err = git.PlainClone(cmd.Dir, false, &git.CloneOptions{
		URL:      cmd.Module,
		Progress: os.Stdout,
	})

	err = cmd.Executor.Execute(TerrformBin, []string{"init"}, "")

	err = cmd.Executor.Execute(TerrformBin, []string{"apply", FlagAutoApprove}, "")

	fmt.Println(err)

	return nil
}

// func main() {
// 	cmd := TerraformRunner{
// 		Module: "https://github.com/ishansd94/terraform-sample-module.git",
// 		Dir: "/tmp/zzzz",
// 	}
// 	_ = cmd.Run()
// }
