package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Constant variables
const (
	reC = `^[0-9]+$`
)

// Directory variables
var (
	buildCurDir  = getBuildCurrentDir()
	generatedDir = filepath.Join(buildCurDir, "generated")
	linkDir      = filepath.Join(buildCurDir, "links")
	sourceDir    = filepath.Join(linkDir, "sources")
	buildDir     = filepath.Join(linkDir, "builds")
	installDir   = filepath.Join(linkDir, "install")
)

// Get the current directory for build
func getBuildCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	return dir
}

// RunCommand runs a bash command and returns the output as a string.
// It takes a command string as input and returns a string and an error if any.
func RunCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error executing command %s: %v\nOutput: %s", command, err, output)
		return "", err
	}
	return string(output), nil
}

// getFolderFromRecipe retrieves a folder path from a Yocto recipe output.
func getFolderFromRecipe(component, prefix string) (string, error) {
	setupCmd := "source " + buildCurDir + "setup_custom_project " + generatedDir
	_, err1 := RunCommand(setupCmd)
	if err1 != nil {
		return "", err1
	}

	bitbakeCmd := "bitbake -e " + component
	output, err := RunCommand(bitbakeCmd)
	if err != nil {
		return "", err
	}

	cleanedOutput, err := removeNotRelevantChar(output, prefix)
	if err != nil {
		return "", err
	}

	return cleanedOutput, nil
}

// removeNotRelevantChar removes the prefix and double quotes (") from the directory path.
// It takes a string and a string as input and returns a string and an error if any.
func removeNotRelevantChar(dir string, prefix string) (string, error) {
	// Remove double quotes from dir
	dir = strings.ReplaceAll(dir, "\"", "")

	// Remove prefix from dir
	dir = strings.TrimPrefix(dir, prefix)

	// Remove any character in left of dir string until first forward slash
	dir = strings.TrimLeft(dir, "/")

	log.Printf("removeNotRelevantChar: dir is %s", dir)

	return dir, nil
}

// Build component
func doBuild(component string) (string, error) {
	setupCmd := "source " + buildCurDir + "setup_custom_project " + generatedDir
	_, err := RunCommand(setupCmd)
	if err != nil {
		return "", err
	}

	cmd := "bitbake " + component
	result, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Retrieve component source path
func doGetSources(sourceFolder string, component string) (string, error) {

	folder, err := getFolderFromRecipe(component, "^S=")
	if err != nil {
		log.Fatal(err)
	}

	//append sourceFolder path to buildCurDir
	sourceFolder = filepath.Join(buildCurDir, sourceFolder)

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, sourceFolder)
	err = os.Symlink(folder, sourceFolder)

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}

// Retrieve component install path
func doGetInstallFolder(installFolder string, component string) (string, error) {
	folder, err := getFolderFromRecipe(component, "^D=")
	if err != nil {
		log.Fatal(err)
	}

	//append sourceFolder path to buildCurDir
	installFolder = filepath.Join(buildCurDir, installFolder)

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, installFolder)
	err = os.Symlink(folder, installFolder)

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}

// Retrieve component build path
func doGetBuildFolder(buildFolderArg string, component string) (string, error) {
	folder, err := getFolderFromRecipe(component, "^B=")
	if err != nil {
		log.Fatal(err)
	}

	//append sourceFolder path to buildCurDir
	buildFolderArg = filepath.Join(buildCurDir, buildFolderArg)

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, buildFolderArg)
	err = os.Symlink(folder, buildFolderArg)

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}

// doPrepare processes command-line arguments and calls the appropriate functions.
func doPrepare() error {

	gComponent := flag.String("c", "", "<component> component to build")
	sourcedComponent := flag.String("s", "", "<component> create symlink source folder.")
	flag.StringVar(sourcedComponent, "source_link", "", "<component> create symlink source folder.")
	builtComponent := flag.String("b", "", "--build_link create symlink build folder for component.")
	installedComponent := flag.String("i", "", "--install_link create symlink install folder for component.")

	flag.Parse()
	if *gComponent != "" {
		result, err := doBuild(*gComponent)
		if err != nil {
			return err
		}

		log.Print(result)

	}

	if *sourcedComponent != "" {
		result, err := doGetSources(buildDir, *sourcedComponent)
		if err != nil {
			return err
		}

		log.Printf("Sources component %s into %s", buildDir, result)

	}
	if *builtComponent != "" {

		buildLink, err := doGetBuildFolder(buildCurDir, buildDir)
		if err != nil {
			return err
		}

		log.Printf("Created build link %s pointing to %s", buildDir, buildLink)
	}
	if *installedComponent != "" {

		installLink, err := doGetInstallFolder(buildCurDir, installDir)
		if err != nil {
			return err
		}

		log.Printf("Created install link %s pointing to %s", installDir, installLink)
	}

	for _, arg := range flag.Args() {
		switch {
		case strings.HasPrefix(arg, "--help"):
		case strings.HasPrefix(arg, "-h"):
			log.Println("wrapper for executing bitbake tasks.")
			flag.Usage()
			return fmt.Errorf("help requested")
		default:
			return fmt.Errorf("unknown argument: %s", arg)
		}
	}

	return nil
}

func main() {
	doPrepare()
	// Process command-line arguments and call appropriate functions
}
