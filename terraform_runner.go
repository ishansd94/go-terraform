package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/ishansd94/terraform-go/executor"
	"github.com/ishansd94/terraform-go/helpers"
)

const (
	TerrformBin = "terraform"

	FlagAutoApprove = "-auto-approve"

	OptionBackendConfig = "-backend-config"
	OptionVar           = "-var"

	OperationApply   = "apply"
	OperationInit    = "init"
	OperationDestroy = "destroy"
	OperationPlan    = "plan"
)

type Executor interface {
	Execute(command string, args []string, prefix string) error
}

type TerraformRunner struct {
	Module        string
	Dir           string
	Version       string
	Operation     string
	Options       *TerraformOptions
	BackendConfig map[string]string
	Writer        io.Writer
	Executor      Executor
}

type TerraformOptions struct {
	Vars map[string]interface{}
	AutoApprove bool
}

func (cmd *TerraformRunner) Run() error {

	var err error

	if cmd.Executor == nil {
		cmd.Executor = &executor.DefaultExecute{
			Dir:    cmd.Dir,
			Writer: cmd.Writer,
		}
	}

	_, err = git.PlainClone(cmd.Dir, false, &git.CloneOptions{
		URL:      cmd.Module,
		Progress: os.Stdout,
	})

	commands, err := cmd.Command()
	if err != nil {
		return err
	}

	fmt.Printf("%s %s", TerrformBin, commands)

	err = cmd.Executor.Execute(TerrformBin, commands, "")
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TerraformRunner) Command() ([]string, error) {
	var commands []string
	var backend []string
	var options []string

	if !helpers.StringSliceContains(getSupportedOperations(), cmd.Operation) {
		return nil, errors.New(fmt.Sprintf("'%s' is an invalid operation", cmd.Operation))
	}

	commands = append(commands, cmd.Operation)

	switch cmd.Operation {
	case OperationInit:
		if len(cmd.BackendConfig) > 0 {
			backend = cmd.GenerateBackendOptions()
		}
	}

	if cmd.Options != nil {
		options = cmd.GenerateOptions()
	}

	if len(backend) > 0 {
		commands = append(commands, backend...)
	}

	if len(options) > 0 {
		commands = append(commands, options...)
	}

	return commands, nil
}

func getSupportedOperations() []string {
	return []string{
		OperationInit,
		OperationApply,
		OperationPlan,
		OperationDestroy,
	}
}

func (cmd *TerraformRunner) GenerateBackendOptions() []string {
	var backend []string
	for k, v := range cmd.BackendConfig {
		backend = append(backend, fmt.Sprintf("%s %s=%s", OptionBackendConfig, k, v))
	}
	return backend
}

func (cmd *TerraformRunner) GenerateOptions() []string {
	var options []string
	switch cmd.Operation {
	case OperationApply, OperationDestroy, OperationPlan:
		for k, v := range cmd.Options.Vars {
			options = append(options, OptionVar)
			options = append(options, fmt.Sprintf("%s=%s", k, v))
		}

		if cmd.Options.AutoApprove {
			options = append(options, FlagAutoApprove)
		}
	}

	return options
}

// func main() {
// 	var err error
//
// 	cmd := TerraformRunner{
// 		Module:    "https://github.com/ishansd94/terraform-sample-module.git",
// 		Dir:       "/tmp/zzzz",
// 		Options: &TerraformOptions{
// 			Vars: map[string]interface{}{
// 				"foo": "foozz",
// 				"bar": "barzz",
// 			},
// 			AutoApprove: true,
// 		},
// 	}
//
// 	cmd.Operation = OperationInit
// 	err = cmd.Run()
// 	fmt.Println(err)
//
// 	cmd.Operation = OperationApply
// 	err = cmd.Run()
// 	fmt.Println(err)
// }
