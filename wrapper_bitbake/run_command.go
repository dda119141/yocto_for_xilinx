// run_commands.go

package wrapper_bitbake

import (
	"bytes"
	"log"
	"os/exec"
)

// RunCommand runs a command with the given arguments and returns the output.
func RunCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	log.Printf("::: %s :::\n", cmd)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("error executing command:: %s :: %v\nOutput: %s", command, err, output)
		return "", err
	}
	return string(bytes.Trim(output, "\n")), nil
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	// if err != nil {
	// 	return "", err
	// }
	// return out.String(), nil
}

// RunCommandWithDir runs a command with the given arguments in the specified directory.
func RunCommandWithDir(command string, dir string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// RunBitbakeCommandWithDir runs a bitbake command with the given arguments in the specified directory.
func RunBitbakeCommandWithDir(dir string, args ...string) (string, error) {
	return RunCommandWithDir("bitbake", dir, args...)
}
