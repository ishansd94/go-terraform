package terraformer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl"

	"github.com/ishansd94/terraform-go/executor"
	"github.com/ishansd94/terraform-go/helpers"

	"github.com/go-git/go-git/v5"
)

const (
	TerraformBin = "terraform"

	// FlagAutoApprove Terraform Flags
	FlagAutoApprove = "-auto-approve"
	FlagForceCopy   = "-force-copy"
	FlagJSON        = "-json"

	// OptionBackendConfig Terraform Options
	OptionBackendConfig = "-backend-config"
	OptionVar           = "-var"
	OptionFromModule    = "-from-module"

	// OperationInit Terraform actions
	OperationInit    = "init"
	OperationApply   = "apply"
	OperationPlan    = "plan"
	OperationDestroy = "destroy"
	OperationState   = "state"
	OperationTaint   = "taint"
	OperationUntaint = "untaint"
	OperationShow    = "show"
	OperationOutput  = "output"

	ArgumentList = "list"
	ArgumentShow = "show"
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

// Run method runs the terraform which used externally
func (cmd *TerraformRunner) Run() error {

	var err error

	args, err := cmd.GenerateArguments()
	if err != nil {
		return err
	}

	output, err := cmd.run(args)
	if err != nil {
		return err
	}

	if cmd.PrintOutput {
		_, _ = io.WriteString(os.Stdout, string(*output))
	}

	return nil
}

// Run method runs the terraform which used internally with default executor
func (cmd *TerraformRunner) run(args []string) (*[]byte, error) {
	var output *[]byte
	var err error

	if cmd.Executor == nil {
		cmd.Executor = &executor.DefaultExecute{
			Directory: cmd.Directory,
			Writer:    cmd.Writer,
		}
	}

	cmd.debug(fmt.Sprintf("Running %s %s\n", TerraformBin, args))

	if output, err = cmd.Executor.Execute(TerraformBin, args, ""); err != nil {
		return nil, err
	}

	return output, nil
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

	case OperationApply, OperationDestroy, OperationPlan:
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
	case OperationApply, OperationDestroy:
		flags = append(flags, FlagAutoApprove)
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

	var v map[string]interface{}

	args := []string{OperationOutput, FlagJSON}

	output, err := cmd.run(args)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(*output, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// State returns the state of the module by running `terraform show -json`
func (cmd *TerraformRunner) State() (*map[string]interface{}, error) {

	var v map[string]interface{}

	args := []string{OperationShow, FlagJSON}

	output, err := cmd.run(args)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(*output, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (cmd *TerraformRunner) Resources() ([]string, error) {
	var r []string

	args := []string{OperationState, ArgumentList}

	output, err := cmd.run(args)
	if err != nil {
		return nil, err
	}

	splitFunc := func(c rune) bool {
		return c == '\n'
	}

	r = strings.FieldsFunc(string(*output), splitFunc)

	return r, nil
}

func (cmd *TerraformRunner) Resource(name string) (interface{}, error) {
	var v interface{}

	args := []string{OperationState, ArgumentShow, name}

	output, err := cmd.run(args)
	if err != nil {
		return nil, err
	}

	if err := hcl.Decode(&v, helpers.SanitizeHCL(string(*output))); err != nil {
		return nil, err
	}

	return v, nil
}

func (cmd *TerraformRunner) Taint(name string) error  {

	args := []string{OperationTaint, name}

	_, err := cmd.run(args)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *TerraformRunner) UnTaint(name string) error  {

	args := []string{OperationUntaint, name}

	_, err := cmd.run(args)
	if err != nil {
		return err
	}

	return nil
}

// getSupportedOperations returns a list of supported terraform operations
func getSupportedOperations() []string {
	return []string{
		OperationInit,
		OperationApply,
		OperationPlan,
		OperationDestroy,
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
