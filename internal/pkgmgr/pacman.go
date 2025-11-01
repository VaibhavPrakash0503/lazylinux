package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func (p *Pacman) Update() error {
	// Pacman command: sudo pacman -Syu --noconfirm
	// -S = sync, -y = refresh repos, -u = upgrade
	cmd := exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Pacman) Clean() error {
	// Clean package cache (keep only current versions)
	fmt.Println("🧹 Cleaning package cache...")
	cleanCmd := exec.Command("sudo", "pacman", "-Sc", "--noconfirm")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	err := cleanCmd.Run()
	if err != nil {
		return err
	}

	// Remove orphaned packages
	fmt.Println("🗑️  Removing orphaned packages...")

	// First check if there are orphaned packages
	checkCmd := exec.Command("pacman", "-Qtdq")
	output, err := checkCmd.Output()
	if err != nil || len(output) == 0 {
		fmt.Println("✨ No orphaned packages found")
		return nil
	}

	// Remove orphaned packages
	removeCmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", string(output))
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	return removeCmd.Run()
}

// List lists all installed packages
func (p *Pacman) List() ([]string, error) {
	// pacman -Q (list all installed)
	cmd := exec.Command("pacman", "-Q")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output
	var results []string
	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// pacman -Q output: "package-name version"
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			results = append(results, parts[0])
		}
	}

	return results, nil
}
