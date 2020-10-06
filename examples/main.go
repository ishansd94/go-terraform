package main

import (
	"fmt"

	"github.com/sanity-io/litter"

	terraformer "github.com/ishansd94/terraform-go"
)

func main() {
	var err error

	cmd := terraformer.TerraformRunner{
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

	// err = cmd.GetModule()
	// fmt.Println(err)

	cmd.Operation = terraformer.OperationInit
	err = cmd.Run()
	fmt.Println(err)

	cmd.Operation = terraformer.OperationApply
	err = cmd.Run()
	fmt.Println(err)

	out1, err := cmd.Output()
	fmt.Println(out1 , err)

	// out2, err := cmd.State()
	// fmt.Println(out2, err)

	out3, err := cmd.Resources()
	litter.Dump(out3, err)

	// out4, err := cmd.Resource("null_resource.obj")
	// litter.Dump(out4, err)

	err = cmd.Taint("null_resource.obj")
	litter.Dump(err)

	cmd.Operation = terraformer.OperationPlan
	err = cmd.Run()
	fmt.Println(err)

	err = cmd.UnTaint("null_resource.obj")
	litter.Dump(err)

	cmd.Operation = terraformer.OperationPlan
	err = cmd.Run()
	fmt.Println(err)
}