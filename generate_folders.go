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

// Display usage information
func usage() {
	log.Println("build.go is a wrapper for executing bitbake tasks.")
	log.Println("Options:")
	log.Println("-c | --component <component> component to build")
	log.Println("-s | --source_link create symlink source folder.")
	log.Println("-b | --build_link create symlink build folder.")
	log.Println("-i | --install_link create symlink install folder.")
	log.Println("-h | --help show help")
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
func getFolderFromRecipe(baseDir, component, prefix string) (string, error) {
	setupCmd := "source " + baseDir + "setup_custom_project " + generatedDir
	output, err := RunCommand(setupCmd)
	if err != nil {
		return "", err
	}

	bitbakeCmd := "bitbake -e " + component
	output, err = RunCommand(bitbakeCmd)
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
func doBuild(currentDir string, component string) (string, error) {
	setupCmd := "source " + buildCurDir + "setup_custom_project " + generatedDir
	_, err := RunCommand(setupCmd)
	if err != nil {
		return "", err
	}

	cmd := "bitbake -e " + component
	result, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Retrieve component source path
func doGetSources(sourceFolder string, component string) (string, error) {

	folder, err := getFolderFromRecipe(sourceFolder, component, "^S=")
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
func doGetInstallFolder(sourceFolder string, component string) (string, error) {
	folder, err := getFolderFromRecipe(sourceFolder, component, "^D=")
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

// Retrieve component build path
func doGetBuildFolder(sourceFolder string, component string) (string, error) {
	folder, err := getFolderFromRecipe(sourceFolder, component, "^B=")
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

// doPrepare processes command-line arguments and calls the appropriate functions.
func doPrepare() error {

	fmt.Errorf("no arguments sdfdfdsgddgs")

	flag.Parse()

	if flag.NArg() == 0 {
		return fmt.Errorf("no arguments provided")
	}

	for _, arg := range flag.Args() {
		switch {
		case strings.HasPrefix(arg, "-c"):
			componentArg := flag.Arg(1)
			if componentArg == "" {
				return fmt.Errorf("-c requires a component name")
			}

			component, err := doBuild(buildCurDir, componentArg)
			if err != nil {
				return err
			}

			log.Printf("Built component %s into %s", componentArg, component)
		case strings.HasPrefix(arg, "-s"):

			componentArg := flag.Arg(1)
			if componentArg == "" {
				return fmt.Errorf("-c requires a component name")
			}

			component, err := doGetSources(buildCurDir, componentArg)
			if err != nil {
				return err
			}

			log.Printf("Sources component %s into %s", componentArg, component)
		case strings.HasPrefix(arg, "-b"):
			buildLinkArg := flag.Arg(1)
			if buildLinkArg == "" {
				return fmt.Errorf("-b requires a build link name")
			}

			buildLink, err := doGetBuildFolder(buildCurDir, buildLinkArg)
			if err != nil {
				return err
			}

			log.Printf("Created build link %s pointing to %s", buildLinkArg, buildLink)
		case strings.HasPrefix(arg, "-i"):
			installLinkArg := flag.Arg(1)
			if installLinkArg == "" {
				return fmt.Errorf("-i requires an install link name")
			}

			installLink, err := doGetInstallFolder(buildCurDir, installLinkArg)
			if err != nil {
				return err
			}

			log.Printf("Created install link %s pointing to %s", installLinkArg, installLink)
		case strings.HasPrefix(arg, "-h"):
			usage()
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
