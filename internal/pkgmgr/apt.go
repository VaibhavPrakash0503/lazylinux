package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type APT struct{}

func NewAPT() *APT {
	return &APT{}
}

func (a *APT) Install(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// APT command: sudo apt install -y <packages>
	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"apt"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *APT) Remove(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// APT command: sudo apt remove -y <packages>
	args := append([]string{"remove", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"apt"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *APT) Update() error {
	// APT update needs two commands: update repo lists, then upgrade packages

	// First: sudo apt update
	updateCmd := exec.Command("sudo", "apt", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	err := updateCmd.Run()
	if err != nil {
		return err
	}

	// Second: sudo apt upgrade -y
	upgradeCmd := exec.Command("sudo", "apt", "upgrade", "-y")
	upgradeCmd.Stdout = os.Stdout
	upgradeCmd.Stderr = os.Stderr
	return upgradeCmd.Run()
}

func (a *APT) Clean() error {
	// Clean package cache
	fmt.Println("🧹 Cleaning package cache...")
	cleanCmd := exec.Command("sudo", "apt", "clean")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	err := cleanCmd.Run()
	if err != nil {
		return err
	}

	// Remove orphaned packages
	fmt.Println("🗑️  Removing orphaned packages...")
	autoremoveCmd := exec.Command("sudo", "apt", "autoremove", "-y")
	autoremoveCmd.Stdout = os.Stdout
	autoremoveCmd.Stderr = os.Stderr
	return autoremoveCmd.Run()
}

// List lists all installed packages
func (a *APT) List() ([]string, error) {
	// apt list --installed
	cmd := exec.Command("apt", "list", "--installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output
	var results []string
	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Listing..." {
			continue
		}
		// apt list output: "package-name/repo,repo version [installed]"
		pkgName := strings.Split(line, "/")[0]
		results = append(results, pkgName)
	}

	return results, nil
}
