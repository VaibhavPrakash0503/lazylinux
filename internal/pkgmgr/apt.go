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
