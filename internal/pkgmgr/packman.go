package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
)

type Pacman struct{}

func NewPacman() *Pacman {
	return &Pacman{}
}

func (p *Pacman) Install(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// Pacman command: sudo pacman -S --noconfirm <packages>
	args := append([]string{"-S", "--noconfirm"}, packages...)
	cmd := exec.Command("sudo", append([]string{"pacman"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Pacman) Remove(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// Pacman command: sudo pacman -R --noconfirm <packages>
	args := append([]string{"-R", "--noconfirm"}, packages...)
	cmd := exec.Command("sudo", append([]string{"pacman"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
