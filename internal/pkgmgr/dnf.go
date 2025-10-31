package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
)

type DNF struct{}

func NewDNF() *DNF {
	return &DNF{}
}

func (d *DNF) Install(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}
	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"dnf"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *DNF) Remove(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}
	args := append([]string{"remove", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"dnf"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *DNF) Update() error {
	// DNF command: sudo dnf update -y
	cmd := exec.Command("sudo", "dnf", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *DNF) Clean() error {
	// Remove orphaned packages
	fmt.Println("🗑️  Removing orphaned packages...")
	autoremoveCmd := exec.Command("sudo", "dnf", "autoremove", "-y")
	autoremoveCmd.Stdout = os.Stdout
	autoremoveCmd.Stderr = os.Stderr
	return autoremoveCmd.Run()
}
