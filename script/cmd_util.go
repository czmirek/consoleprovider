package script

import (
	"log"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func runScript(d *schema.ResourceData, getOutput bool, op string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[INFO] Running script %v\r\n", op)

	opList := d.Get(op).([]interface{})
	workingDir := d.Get("working_dir").(string)

	if err := validateProgramAttr(opList); err != nil {
		return "", diag.FromErr(err)
	}

	program := make([]string, len(opList))

	for i, vI := range opList {
		program[i] = vI.(string)
	}

	cmd := exec.Command(program[0], program[1:]...)
	cmd.Dir = workingDir

	if getOutput {
		resultBytes, err := cmd.Output()
		resultJSON := string(resultBytes)
		log.Printf("[TRACE] JSON output: %+v\r\n", resultJSON)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.Stderr != nil && len(exitErr.Stderr) > 0 {
					return "", diag.Errorf("failed to execute %q: %s", program[0], string(exitErr.Stderr))
				}
				return "", diag.Errorf("command %q failed with no error message", program[0])
			} else {
				return "", diag.Errorf("failed to execute %q: %s", program[0], err)
			}
		}
		return resultJSON, diags
	}

	if err := cmd.Run(); err != nil {
		return "", diag.FromErr(err)
	}
	return "", diags
}
