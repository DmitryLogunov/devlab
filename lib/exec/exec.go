package exec

import (
	"os/exec"
)

/**
*  Executes git command in service folder
*/
func GitCommand(serviceDir string, command string)  (result string, err error) {
  out, err := exec.Command(
    "sh", 
    "-c", 
		"( cd " + serviceDir + " && " + command + ")" ).Output()
		
	return string(out), err
}