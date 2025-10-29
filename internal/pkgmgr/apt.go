package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
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
