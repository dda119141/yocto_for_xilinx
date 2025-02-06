package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], " <package_name>")
		os.Exit(1)
	}
	packageName := os.Args[1]

	if err := generateDependencyList(packageName); err != nil {
		fmt.Println("Error generating dependency list:", err)
		return
	}

	if err := downloadDependencies(); err != nil {
		fmt.Println("Error downloading dependencies:", err)
		return
	}
}

func generateDependencyList(packageName string) error {
	return runCmd("bitbake", "-g", packageName)
}

func downloadDependencies() error {
	file, err := os.Open("pn-buildlist")
	if err != nil {
		return fmt.Errorf("error opening pn-buildlist: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if err := runCmd("bitbake", "-c", "fetch", line); err != nil {
				fmt.Printf("Error fetching %s: %v\n", line, err)
				continue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading pn-buildlist: %w", err)
	}
	return nil
}

func runCmd(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running '%s %v': %w", command, arguments, err)
	}
	return nil
}
