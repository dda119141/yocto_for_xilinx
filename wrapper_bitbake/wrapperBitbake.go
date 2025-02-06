package wrapper_bitbake

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildDirsFromCurDir returns a struct containing the directories relative to the provided curdir.
//
// It takes a string as input and returns a Directories struct.
func BuildDirsFromCurDir(curdir string) Directories {
	linkDir := filepath.Join(curdir, "links")
	if _, err := os.Stat(linkDir); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(curdir, "links"), os.ModePerm)
	}
	return Directories{
		buildDir:   filepath.Join(curdir, "links", "builds"),
		installDir: filepath.Join(curdir, "links", "install"),
		sourceDir:  filepath.Join(curdir, "links", "sources"),
		// generatedDir: filepath.Join(curdir, "generated"),
		generatedDir: "/home/max/Work/yocto/xilinx/build_zynq",
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
func Run_Command(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	log.Printf("::: %s :::\n", cmd)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("error executing command:: %s :: %v\nOutput: %s", command, err, output)
		return "", err
	}
	log.Println(output)
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
		log.Printf("error setting up custom project:: %v\n", err1)
		return "", err1
	}

	bitbakeCmd := "bitbake -e " + component + " | grep " + prefix

	output, err := RunCommand(bitbakeCmd)
	if err != nil {
		log.Printf("error executing bitbake -e command:: %v\n", err)
		return "", err
	}

	cleanedOutput, err := removePrefixAndQuotes(output, prefix)
	if err != nil {
		log.Printf("error removePrefixAndQuotes:: %v\n", err)
		return "", err
	}

	return cleanedOutput, nil
}

func removeLeftOfFirstSlash(dir string) string {
	index := strings.Index(dir, "/")
	if index == -1 {
		return dir // No slash found, return original string
	}
	return dir[index:] // Found slash, return string from slash onwards
}

// removePrefixAndQuotes removes the prefix and double quotes from the directory path.
func removePrefixAndQuotes(dir string, prefix string) (string, error) {
	// Remove double quotes from dir
	quotedDir := strings.ReplaceAll(dir, "\"", "")

	// Remove prefix from dir
	unprefixedDir := strings.TrimPrefix(quotedDir, prefix)

	// Remove any character in left of dir string until first forward slash
	finalDir := removeLeftOfFirstSlash(unprefixedDir)

	return finalDir, nil
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

// FolderExists checks if the specified folderPath exists and is a directory.
// It takes a folderPath of type string as input and returns a bool indicating existence.
func FolderExists(folderPath string) bool {
	stat, err := os.Stat(folderPath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// templateFolder retrieves a folder path from a Yocto recipe output.
// It takes a component and a prefix as input, both of type string, and returns a string and an error if any.
func templateFolder(d Directories, component string, prefix string) (string, error) {
	//if d not initialized, return error
	if d.getTopDir() == "" {
		return "", errors.New("directories not initialized")
	}

	folder, err := getFolderFromRecipe(component, d, prefix)
	if err != nil {
		return "", err
	}

	// make path from folder
	folderPath := filepath.Clean(folder)

	// check if folderPath exists
	if !FolderExists(folderPath) {
		return "", fmt.Errorf("the folder path %s does not exist", folderPath)
	}

	// create symbolic link from folder into sources folder
	linkPath := ""
	if prefix == "^S=" {
		linkPath = d.getSourceDir()
	} else if prefix == "^D=" {
		linkPath = d.getInstallDir()
	} else if prefix == "^B=" {
		linkPath = d.getBuildDir()
	} else {
		return "", fmt.Errorf("prefix %s not supported", prefix)
	}

	log.Printf("create symbolic link from %s into %s", folderPath, linkPath)
	err = os.Symlink(folderPath, linkPath)

	if err != nil {
		return "", err
	}

	return folderPath, nil
}

// DoGetSources retrieves a component source path.
func DoGetSources(d Directories, component string) (string, error) {
	return templateFolder(d, component, "^S=")
}

// DoGetInstallFolder retrieves a component install path.
func DoGetInstallFolder(d Directories, component string) (string, error) {
	return templateFolder(d, component, "^D=")
}

// DoGetBuildFolder retrieves a component build path.
func DoGetBuildFolder(d Directories, component string) (string, error) {
	return templateFolder(d, component, "^B=")
}
