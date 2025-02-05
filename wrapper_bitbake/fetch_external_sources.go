package wrapper_bitbake

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	xilinxGit    = "https://github.com/Xilinx/"
	xilinxBranch = "rel-v2022.2"
)

var (
	curDir      = getCurrentDir()
	arrayRepos  = []string{"meta-xilinx", "meta-openembedded", "poky", "meta-xilinx-tools", "meta-xilinx-tsn"}
	externalDir = filepath.Join(curDir, "external")
	hostTools   = []string{"gawk", "wget", "git", "diffstat", "unzip", "texinfo", "gcc", "build-essential", "chrpath", "socat", "cpio", "python3", "python3-pip", "python3-pexpect", "xz-utils", "debianutils", "iputils-ping", "python3-git", "python3-jinja2", "libegl1-mesa", "libsdl1.2-dev", "pylint3", "xterm", "python3-subunit", "mesa-common-dev", "zstd", "liblz4-tool"}
)

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	return dir
}

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error executing command %s: %v\nOutput: %s", command, err, output)
		return err
	}
	return nil
}

// doInstallHostPackages installs the required packages on the host.
// It takes no arguments and returns an error if any.
func doInstallHostPackages() error {
	for _, tool := range hostTools {
		if !isInstalled(tool) {
			log.Printf("%s to be installed on host\n", tool)
			if err := runCommand(fmt.Sprintf("sudo apt install -y %s", tool)); err != nil {
				log.Printf("Error installing %s: %v\n", tool, err)
				return err
			}
		}
	}
	return nil
}

func isInstalled(tool string) bool {
	cmd := exec.Command("dpkg", "-l")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error checking installed packages:", err)
		return false
	}
	return strings.Contains(string(output), tool)
}

func doCheckIfRepoExists(repo string) bool {
	dir := filepath.Join(externalDir, repo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	if err := runCommand(fmt.Sprintf("git -C %s remote -v", dir)); err != nil {
		return false
	}
	return true
}

// doFetchRepos clones the specified repositories if they don't already exist.
// It takes no arguments and returns an error if any.
func doFetchRepos() error {
	if _, err := os.Stat(externalDir); !os.IsNotExist(err) {
		log.Printf("The directory %s already exists. Skipping repository cloning.\n", externalDir)
		return nil
	}

	for _, repo := range arrayRepos {
		repoCloneURL := xilinxGit + repo + ".git"
		if !doCheckIfRepoExists(repo) {
			if err := runCommand(fmt.Sprintf("git clone %s -b %s %s", repoCloneURL, xilinxBranch, externalDir)); err != nil {
				log.Printf("Error cloning repository %s: %v\n", repo, err)
				return err
			}
		}
	}
	return nil
}

// installZynqRepos installs the required packages and clones the Xilinx repositories.
//
// It takes no arguments and returns an error if any.
func InstallZynqRepos() error {
	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		if err = os.Mkdir(externalDir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := doInstallHostPackages(); err != nil {
		return err
	}

	if err := doFetchRepos(); err != nil {
		return err
	}

	return nil
}
