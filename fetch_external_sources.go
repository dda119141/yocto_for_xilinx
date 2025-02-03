package main

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

func doInstallHostPackages() {
	for _, tool := range hostTools {
		if !isInstalled(tool) {
			log.Printf("%s to be installed on host\n", tool)
			if err := runCommand(fmt.Sprintf("sudo apt install -y %s", tool)); err != nil {
				log.Println(err)
			}
		}
	}
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

func doFetchRepos() {
	if _, err := os.Stat(externalDir); !os.IsNotExist(err) {
		log.Printf("The directory %s already exists. Skipping repository cloning.\n", externalDir)
		return
	}

	for _, repo := range arrayRepos {
		repoCloneURL := xilinxGit + repo + ".git"
		if !doCheckIfRepoExists(repo) {
			if err := runCommand(fmt.Sprintf("git clone %s -b %s %s", repoCloneURL, xilinxBranch, externalDir)); err != nil {
				log.Printf("Error cloning repository %s: %v\n", repo, err)
			}
		}
	}
}

func main() {
	if _, err := os.Stat(externalDir); os.IsNotExist(err) {
		os.Mkdir(externalDir, os.ModePerm)
	}

	doInstallHostPackages()
	doFetchRepos()
}
