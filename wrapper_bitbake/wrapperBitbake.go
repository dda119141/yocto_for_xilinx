package wrapper_bitbake

import (
	"errors"
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

// buildDirsFromCurDir returns a struct containing the directories relative to the provided curdir.
//
// It takes a string as input and returns a Directories struct.
func BuildDirsFromCurDir(curdir string) Directories {
	return Directories{
		buildDir:     filepath.Join(curdir, "links", "builds"),
		installDir:   filepath.Join(curdir, "links", "install"),
		sourceDir:    filepath.Join(curdir, "links", "sources"),
		generatedDir: filepath.Join(curdir, "generated"),
		topdir:       curdir,
	}
}

// Directories contains the directories relative to the build current directory.
type Directories struct {
	buildDir     string
	installDir   string
	sourceDir    string
	generatedDir string
	topdir       string
}

func (d Directories) getBuildDir() string {
	return d.buildDir
}
func (d Directories) getInstallDir() string {
	return d.installDir
}
func (d Directories) getSourceDir() string {
	return d.sourceDir
}
func (d Directories) getGeneratedDir() string {
	return d.generatedDir
}
func (d Directories) getTopDir() string {
	return d.topdir
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
// It takes a component and a prefix as input, both of type string, and returns a string and an error if any.
func getFolderFromRecipe(component string, d Directories, prefix string) (string, error) {

	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	setupCmd := "source " + d.getTopDir() + "/setup_custom_project " + d.getGeneratedDir()
	_, err1 := RunCommand(setupCmd)
	if err1 != nil {
		return "", err1
	}

	bitbakeCmd := "bitbake -e " + component + " | grep " + prefix

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
func DoBuild(d Directories, component string) (string, error) {
	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	setupCmd := "source " + d.getTopDir() + "/setup_custom_project " + d.getGeneratedDir()
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
func DoGetSources(d Directories, component string) (string, error) {

	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	folder, err := getFolderFromRecipe(component, d, "^S=")
	if err != nil {
		log.Fatal(err)
	}

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, d.getSourceDir())
	err = os.Symlink(folder, d.getSourceDir())

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}

// Retrieve component install path
func DoGetInstallFolder(d Directories, component string) (string, error) {
	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	folder, err := getFolderFromRecipe(component, d, "^D=")
	if err != nil {
		log.Fatal(err)
	}

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, d.getInstallDir())
	err = os.Symlink(folder, d.getInstallDir())

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}

// Retrieve component build path
func DoGetBuildFolder(d Directories, component string) (string, error) {
	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	folder, err := getFolderFromRecipe(component, d, "^B=")
	if err != nil {
		log.Fatal(err)
	}

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folder, d.getBuildDir())
	err = os.Symlink(folder, d.getBuildDir())

	if err != nil {
		log.Fatal(err)
	}

	return folder, nil
}
