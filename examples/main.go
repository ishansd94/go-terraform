package main

import (
	"fmt"

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

	err = cmd.GetModule()
	fmt.Println(err)

	cmd.Operation = terraformer.OperationInit
	err = cmd.Run()
	fmt.Println(err)

	cmd.Operation = terraformer.OperationApply
	err = cmd.Run()
	fmt.Println(err)

	out, err := cmd.Output()
	fmt.Println(out)

}