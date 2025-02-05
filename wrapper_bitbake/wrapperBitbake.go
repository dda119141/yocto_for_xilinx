package wrapper_bitbake

import (
	"errors"
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
	log.Printf("::::: %s \n", cmd)

	output, err := cmd.Output()
	if err != nil {
		log.Printf("error executing command %s: %v\nOutput: %s", command, err, output)
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
		return "", err1
	}

	bitbakeCmd := "bitbake -e " + component + " | grep " + prefix

	output, err := RunCommand(bitbakeCmd)
	if err != nil {
		return "", err
	}

	cleanedOutput, err := removePrefixAndQuotes(output, prefix)
	if err != nil {
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

	// make path from folder
	folderPath := filepath.Clean(folder)

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folderPath, d.getSourceDir())
	err = os.Symlink(folderPath, d.getSourceDir())

	if err != nil {
		log.Fatal(err)
	}

	return folderPath, nil
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
	// make path from folder
	folderPath := filepath.Clean(folder)

	// create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folderPath, d.getInstallDir())
	err = os.Symlink(folderPath, d.getInstallDir())

	if err != nil {
		log.Fatal(err)
	}

	return folderPath, nil
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

	//make path from folder
	folderPath := filepath.Clean(folder)

	//create symbolic link from folder into sources folder
	log.Printf("create symbolic link from %s into %s", folderPath, d.getBuildDir())
	err = os.Symlink(folderPath, d.getBuildDir())

	if err != nil {
		log.Fatal(err)
	}

	return folderPath, nil
}
