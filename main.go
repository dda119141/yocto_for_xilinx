package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"yocto_for_xilinx/wrapper_bitbake"
)

func GetBuildCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	return dir
}

// doPrepare processes command-line arguments and calls the appropriate functions.
func doPrepare() error {
	component := flag.String("c", "", "component to build")
	sourceComponent := flag.String("s", "", "create symlink source folder for component")
	buildComponent := flag.String("b", "", "create symlink build folder for component")
	installComponent := flag.String("i", "", "create symlink install folder for component")
	fetchCmd := flag.String("f", "", "fetch bsp sources for xilinx zynq")

	flag.Parse()

	dires := wrapper_bitbake.BuildDirsFromCurDir(GetBuildCurrentDir())

	if *component != "" {
		_, err := wrapper_bitbake.DoBuild(dires, *component)
		if err != nil {
			return err
		}
	}

	if *sourceComponent != "" {
		_, err := wrapper_bitbake.DoGetSources(dires, *sourceComponent)
		if err != nil {
			return err
		}
	}

	if *buildComponent != "" {
		_, err := wrapper_bitbake.DoGetBuildFolder(dires, *buildComponent)
		if err != nil {
			return err
		}
	}

	if *installComponent != "" {
		_, err := wrapper_bitbake.DoGetInstallFolder(dires, *installComponent)
		if err != nil {
			return err
		}
	}
	if *fetchCmd != "" {
		// execute install_zynq_repo function from fetch_external_sources.go
		err := wrapper_bitbake.InstallZynqRepos()
		if err != nil {
			return err
		}
	}

	for _, arg := range flag.Args() {
		switch {
		case strings.HasPrefix(arg, "--help"):
		case strings.HasPrefix(arg, "-h"):
			flag.Usage()
			return fmt.Errorf("help requested")
		default:
			flag.Usage()
			return fmt.Errorf("unknown argument: %s", arg)
		}
	}
	if flag.NArg() == 0 {
		flag.Usage()
		return fmt.Errorf("missing argument")
	}

	return nil
}

func main() {
	doPrepare()
	// Process command-line arguments and call appropriate functions
}
