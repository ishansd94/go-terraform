package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ishansd94/terraform-go/executor"
	"github.com/ishansd94/terraform-go/helpers"

	"github.com/go-git/go-git/v5"
)

const (
	TerrformBin = "terraform"

	FlagAutoApprove = "-auto-approve"
	FlagForceCopy 	= "-force-copy"
	FlagJSON        = "-json"

	OptionBackendConfig = "-backend-config"
	OptionVar           = "-var"
	OptionFromModule    = "-from-module"

	OperationApply   = "apply"
	OperationInit    = "init"
	OperationDestroy = "destroy"
	OperationPlan    = "plan"
	OperationOutput  = "output"
)

type Executor interface {
	Execute(command string, args []string, prefix string) (*[]byte, error)
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
	PrintOutput   bool
	Debug         bool
}

type TerraformOptions struct {
	Vars        map[string]interface{}
	AutoApprove bool
	ForceCopy   bool
	FromModule  bool
}

func (cmd *TerraformRunner) Run() error {

	var err error

	if cmd.Executor == nil {
		cmd.Executor = &executor.DefaultExecute{
			Dir:    cmd.Dir,
			Writer: cmd.Writer,
		}
	}

	commands, err := cmd.Command()
	if err != nil {
		return err
	}

	cmd.debug(fmt.Sprintf("Running %s %s\n", TerrformBin, commands))

	output, err := cmd.Executor.Execute(TerrformBin, commands, "")
	if err != nil {
		return err
	}

	if cmd.PrintOutput {
		_, _ = io.WriteString(os.Stdout, string(*output))
	}

	return nil
}

func (cmd *TerraformRunner) Output() (*map[string]interface{}, error) {

	var output *[]byte
	var err error
	var v map[string]interface{}

	commands := []string{OperationOutput, FlagJSON}

	cmd.debug(fmt.Sprintf("Running %s %s\n", TerrformBin, commands))

	if output, err = cmd.Executor.Execute(TerrformBin, commands, ""); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(*output, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (cmd *TerraformRunner) Command() ([]string, error) {
	var commands []string
	var backend []string
	var options []string

	if !helpers.InStringSlice(getSupportedOperations(), cmd.Operation) {
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

func (cmd *TerraformRunner) GenerateBackendOptions() []string {
	var backend []string
	for k, v := range cmd.BackendConfig {
		backend = append(backend, OptionBackendConfig, fmt.Sprintf("%s=%s", k, v))
	}
	return backend
}

func (cmd *TerraformRunner) GenerateOptions() []string {
	var options []string

	switch cmd.Operation {
	case OperationInit:
		if cmd.Options.ForceCopy {
			options = append(options, FlagForceCopy)
		}
		if cmd.Options.FromModule {

			source := helpers.TrimString(cmd.Module, map[string]string{
				":": "/",
				"https///" : "",
				"git@" : "",
				".git": "",
			})
			options = append(options, OptionFromModule, source)
		}

	case OperationApply, OperationDestroy, OperationPlan:
		for k, v := range cmd.Options.Vars {
			options = append(options, OptionVar, fmt.Sprintf("%s=%x", k, v))
		}

		if cmd.Options.AutoApprove {
			options = append(options, FlagAutoApprove)
		}
	}

	return options
}

func (cmd *TerraformRunner) GetModule() error {
	var err error

	if _, err = os.Stat(cmd.Dir); !os.IsNotExist(err) {
		cmd.debug("deleting existing module")
		if err = cmd.cleanModule(); err != nil {
			return err
		}
	}

	cmd.debug("cloning module")
	_, err = git.PlainClone(cmd.Dir, false, &git.CloneOptions{
		URL:      cmd.Module,
	})
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TerraformRunner) cleanModule() error {
	if err := os.RemoveAll(cmd.Dir); err != nil {
		return err
	}
	return nil
}

func getSupportedOperations() []string {
	return []string{
		OperationInit,
		OperationApply,
		OperationPlan,
		OperationDestroy,
	}
}

func (cmd *TerraformRunner) debug(msg string) {
	if cmd.Debug {
		fmt.Println(msg)
	}
}

// func main() {
// 	var err error
//
// 	cmd := TerraformRunner{
// 		Module:      "https://github.com/ishansd94/terraform-sample-module.git",
// 		Dir:         "/tmp/zzzz4",
// 		PrintOutput: true,
// 		Debug:       true,
// 		Options: &TerraformOptions{
// 			Vars: map[string]interface{}{
// 				"str": "foooz",
// 				"num": 2,
// 			},
// 			AutoApprove: true,
// 		},
// 	}
//
// 	err  = cmd.GetModule()
// 	fmt.Println(err)
//
// 	cmd.Operation = OperationInit
// 	err = cmd.Run()
// 	fmt.Println(err)
//
// 	cmd.Operation = OperationApply
// 	err = cmd.Run()
// 	fmt.Println(err)
//
// 	out, err := cmd.Output()
// 	fmt.Println(out)
//
// }
