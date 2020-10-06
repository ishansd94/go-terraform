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

	OperationInit    = "init"
	OperationApply   = "apply"
	OperationPlan    = "plan"
	OperationDestroy = "destroy"
	OperationTaint   = "taint"
	OperationUntaint = "untaint"
	OperationShow    = "show"
	OperationOutput  = "output"
)

type Executor interface {
	Execute(command string, args []string, prefix string) (*[]byte, error)
}

type TerraformRunner struct {
	Module        string
	Version       string
	Directory     string
	Operation     string

	Inputs        map[string]interface{}
	BackendConfig *TerraformBackendConfig
	Options       *TerraformOptions
	Flags         *TerraformFlags

	Writer        io.Writer
	Executor      Executor

	PrintOutput   bool
	Debug         bool
}

type TerraformFlags struct {
	FlagForceCopy bool
}

type TerraformOptions struct {
	FromModule  bool
}

type TerraformBackendConfig struct {

}

func (cmd *TerraformRunner) Run() error {

	var err error

	if cmd.Executor == nil {
		cmd.Executor = &executor.DefaultExecute{
			Directory: cmd.Directory,
			Writer:    cmd.Writer,
		}
	}

	args, err := cmd.GenerateArguments()
	if err != nil {
		return err
	}

	cmd.debug(fmt.Sprintf("Running %s %s\n", TerrformBin, args))

	output, err := cmd.Executor.Execute(TerrformBin, args, "")
	if err != nil {
		return err
	}

	if cmd.PrintOutput {
		_, _ = io.WriteString(os.Stdout, string(*output))
	}

	return nil
}



func (cmd *TerraformRunner) GenerateArguments() ([]string, error) {
	var args []string
	var backend []string
	var flags  []string
	var options []string
	var vars  []string

	if !helpers.InStringSlice(getSupportedOperations(), cmd.Operation) {
		return nil, errors.New(fmt.Sprintf("'%s' is an invalid operation", cmd.Operation))
	}

	args = append(args, cmd.Operation)

	if cmd.Flags != nil {
		flags = cmd.GenerateFlags()
	}

	if cmd.Options != nil {
		options = cmd.GenerateOptions()
	}


	switch cmd.Operation {
	case OperationInit:
		if cmd.BackendConfig != nil {
			backend = cmd.GenerateBackendConfig()
		}
	case OperationApply, OperationPlan, OperationDestroy:
		flags = append(flags, FlagAutoApprove)
		if len(cmd.Inputs) > 0 {
			vars = cmd.GenerateInputs()
		}
	}

	if len(backend) > 0 {
		args = append(args, backend...)
	}

	if len(flags) > 0 {
		args = append(args, flags...)
	}

	if len(options) > 0 {
		args = append(args, options...)
	}

	if len(vars) > 0 {
		args = append(args, vars...)
	}

	return args, nil
}

func (cmd *TerraformRunner) GenerateBackendConfig() []string {
	var backend []string

	// TODO implementation

	return backend
}

func (cmd *TerraformRunner) GenerateFlags() []string {
	var flags []string

	switch cmd.Operation {
	case OperationInit:
		if cmd.Flags.FlagForceCopy{
			flags = append(flags, FlagForceCopy)
		}
	}

	return flags
}

func (cmd *TerraformRunner)GenerateInputs() []string {
	var vars []string

	for k, v := range cmd.Inputs {
		vars = append(vars, OptionVar, fmt.Sprintf("%s=%x", k, v))
	}

	return vars
}

func (cmd *TerraformRunner) GenerateOptions() []string {
	var options []string

	switch cmd.Operation {
	case OperationInit:
		if cmd.Options.FromModule {
			source := helpers.TrimString(cmd.Module, map[string]string{
				":": "/",
				"https///" : "",
				"git@" : "",
				".git": "",
			})
			options = append(options, OptionFromModule, source)
		}
	}

	return options
}

func (cmd *TerraformRunner) GetModule() error {
	var err error

	if _, err = os.Stat(cmd.Directory); !os.IsNotExist(err) {
		cmd.debug("deleting existing module")
		if err = cmd.cleanModule(); err != nil {
			return err
		}
	}

	cmd.debug("cloning module")
	_, err = git.PlainClone(cmd.Directory, false, &git.CloneOptions{
		URL:      cmd.Module,
	})
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TerraformRunner) cleanModule() error {
	if err := os.RemoveAll(cmd.Directory); err != nil {
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
		OperationTaint,
		OperationUntaint,
		OperationShow,
	}
}

func (cmd *TerraformRunner) debug(msg string) {
	if cmd.Debug {
		fmt.Println(msg)
	}
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

// func main() {
// 	var err error
//
// 	cmd := TerraformRunner{
// 		Module:      "https://github.com/ishansd94/terraform-sample-module.git",
// 		Version: "master",
// 		Directory:   "/tmp/test",
//
// 		PrintOutput: true,
// 		Debug:       true,
//
// 		Inputs: map[string]interface{}{
// 			"str": "foooz",
// 			"num": 2,
// 		},
//
// 		// BackendConfig: &TerraformBackendConfig{},
//
// 		// Options: &TerraformOptions{
// 		// 	FromModule: true,
// 		// },
//
// 		// Flags: &TerraformFlags{
// 		// 	FlagForceCopy: true,
// 		// },
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
