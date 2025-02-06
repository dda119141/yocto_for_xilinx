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
	arrayRepos = []string{"meta-xilinx", "meta-openembedded", "poky", "meta-xilinx-tools", "meta-xilinx-tsn"}
	hostTools  = []string{"gawk", "wget", "git", "diffstat", "unzip", "texinfo", "gcc", "build-essential", "chrpath", "socat", "cpio", "python3", "python3-pip", "python3-pexpect", "xz-utils", "debianutils", "iputils-ping", "python3-git", "python3-jinja2", "libegl1-mesa", "libsdl1.2-dev", "pylint3", "xterm", "python3-subunit", "mesa-common-dev", "zstd", "liblz4-tool"}
)

// doInstallHostPackages installs the required packages on the host.
// It takes no arguments and returns an error if any.
func doInstallHostPackages() error {
	for _, tool := range hostTools {
		if !isInstalled(tool) {
			log.Printf("%s to be installed on host\n", tool)
			if _, err := RunCommand(fmt.Sprintf("sudo apt install -y %s", tool)); err != nil {
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

func repoExists(repo string, d Directories) bool {
	repoDir := filepath.Join(d.getTopDir(), "external", repo)
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return false
	}

	result, err := RunCommand(fmt.Sprintf("git -C %s remote -v", repoDir))
	if err != nil {
		return false
	}
	if len(result) > 0 {
		return true
	}
	return false
}

// doFetchRepos clones the specified repositories if they don't already exist.
// It takes no arguments and returns an error if any.
func doFetchRepos(d Directories) error {
	externalDir := filepath.Join(d.getTopDir(), "external")

	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		log.Printf("Creating directory: %s\n", externalDir)
		if err := os.Mkdir(externalDir, 0755); err != nil {
			return err
		}
	}

	for _, repoName := range arrayRepos {
		repoURL := xilinxGit + repoName + ".git"
		repoPath := filepath.Join(externalDir, repoName)
		if !repoExists(repoName, d) {
			log.Printf("Cloning repository: %s\n", repoName)
			if _, err := RunCommand(fmt.Sprintf("git clone %s -b %s %s", repoURL, xilinxBranch, repoPath)); err != nil {
				return fmt.Errorf("failed to clone repository %s: %w", repoName, err)
			}
		} else {
			log.Printf("Repository %s already exists. Skipping...\n", repoName)
		}
	}
	log.Println("All repositories processed successfully")
	return nil
}

// installZynqRepos installs the required packages and clones the Xilinx repositories.
//
// It takes no arguments and returns an error if any.
func InstallZynqRepos(d Directories) error {
	externalDir := filepath.Join(d.getTopDir(), "external")

	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		if err = os.Mkdir(externalDir, os.ModePerm); err != nil {
			return err
		}
	} else {
		log.Printf("installZynqRepos: The directory %s already exists. Skipping directory creation.", externalDir)
	}

	if err := doInstallHostPackages(); err != nil {
		log.Printf("installZynqRepos: Error installing required packages: %v", err)
		return err
	}

	if err := doFetchRepos(d); err != nil {
		log.Printf("installZynqRepos: Error cloning repositories: %v", err)
		return err
	}

	return nil
}
