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

	// Terraform Flags
	FlagAutoApprove = "-auto-approve"
	FlagForceCopy   = "-force-copy"
	FlagJSON        = "-json"

	// Terraform Options
	OptionBackendConfig = "-backend-config"
	OptionVar           = "-var"
	OptionFromModule    = "-from-module"

	// Terraform actions
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

// TerraformRunner object is the main object which defines the `terraform` command and how to execute it.
type TerraformRunner struct {
	Module    string
	Version   string
	Directory string
	Operation string

	Inputs        map[string]interface{}
	BackendConfig *TerraformBackendConfig
	Options       *TerraformOptions
	Flags         *TerraformFlags

	Writer   io.Writer
	Executor Executor

	PrintOutput bool
	Debug       bool
}

// TerraformFlags is a collection of supported terraform flags
type TerraformFlags struct {
	FlagForceCopy bool
}

// TerraformOptions is a collection of supported terraform options
type TerraformOptions struct {
	FromModule bool
}

// TerraformBackendConfig is a collection of supported terraform init backend options
type TerraformBackendConfig struct {
}

// Run method runs the terraform
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

// GenerateArguments method generate the terraform arguments
func (cmd *TerraformRunner) GenerateArguments() ([]string, error) {
	var args []string
	var backend []string
	var flags []string
	var options []string
	var vars []string

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

// GenerateBackendConfig method generate the terraform arguments for `terraform init`
func (cmd *TerraformRunner) GenerateBackendConfig() []string {
	var backend []string

	// TODO implementation

	return backend
}

// GenerateInputs method generate the terraform `-var` arguments for `terraform apply, plan, destroy`
func (cmd *TerraformRunner) GenerateInputs() []string {
	var vars []string

	for k, v := range cmd.Inputs {
		vars = append(vars, OptionVar, fmt.Sprintf("%s=%x", k, v))
	}

	return vars
}

// GenerateFlags method generate the terraform flags
func (cmd *TerraformRunner) GenerateFlags() []string {
	var flags []string

	switch cmd.Operation {
	case OperationInit:
		if cmd.Flags.FlagForceCopy {
			flags = append(flags, FlagForceCopy)
		}
	}

	return flags
}

// GenerateOptions method generate the terraform options
func (cmd *TerraformRunner) GenerateOptions() []string {
	var options []string

	switch cmd.Operation {
	case OperationInit:
		if cmd.Options.FromModule {
			source := helpers.TrimString(cmd.Module, map[string]string{
				":":        "/",
				"https///": "",
				"git@":     "",
				".git":     "",
			})
			options = append(options, OptionFromModule, source)
		}
	}

	return options
}

// Output returns the output variables of a terraform module by running `terraform output -json`
func (cmd *TerraformRunner) Output() (*map[string]interface{}, error) {

	var output *[]byte
	var err error
	var v map[string]interface{}

	args := []string{OperationOutput, FlagJSON}

	cmd.debug(fmt.Sprintf("Running %s %s\n", TerrformBin, args))

	if output, err = cmd.Executor.Execute(TerrformBin, args, ""); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(*output, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// getSupportedOperations returns a list of supported terraform operations
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

// debug prints a debug statement if enabled
func (cmd *TerraformRunner) debug(msg string) {
	if cmd.Debug {
		fmt.Println(msg)
	}
}

// GetModule clones the given repo into a directory
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
		URL: cmd.Module,
	})
	if err != nil {
		return err
	}

	return nil
}

// cleanModule deletes the given directory
func (cmd *TerraformRunner) cleanModule() error {
	if err := os.RemoveAll(cmd.Directory); err != nil {
		return err
	}
	return nil
}




func main() {
	var err error

	cmd := TerraformRunner{
		Module:    "https://github.com/ishansd94/terraform-sample-module.git",
		Version:   "master",
		Directory: "/tmp/test",

		PrintOutput: true,
		Debug:       true,

		Inputs: map[string]interface{}{
			"str": "foooz",
			"num": 2,
		},

		// BackendConfig: &TerraformBackendConfig{},

		// Options: &TerraformOptions{
		// 	FromModule: true,
		// },

		// Flags: &TerraformFlags{
		// 	FlagForceCopy: true,
		// },
	}

	err = cmd.GetModule()
	fmt.Println(err)

	cmd.Operation = OperationInit
	err = cmd.Run()
	fmt.Println(err)

	cmd.Operation = OperationApply
	err = cmd.Run()
	fmt.Println(err)

	out, err := cmd.Output()
	fmt.Println(out)

}
