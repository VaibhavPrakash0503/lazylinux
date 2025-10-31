package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
)

// Flatpak represents the Flatpak package manager
type Flatpak struct{}

// NewFlatpak creates a new Flatpak instance
func NewFlatpak() *Flatpak {
	return &Flatpak{}
}

// Install installs packages via Flatpak from Flathub
func (f *Flatpak) Install(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// Flatpak command: flatpak install -y flathub <packages>
	for _, pkg := range packages {
		args := []string{"install", "-y", "flathub", pkg}
		cmd := exec.Command("flatpak", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install %s: %v", pkg, err)
		}
	}

	return nil
}

// Remove uninstalls packages from Flatpak
func (f *Flatpak) Remove(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	// Flatpak command: flatpak uninstall -y <packages>
	args := append([]string{"uninstall", "-y"}, packages...)
	cmd := exec.Command("flatpak", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Update updates all Flatpak packages
func (f *Flatpak) Update() error {
	// Flatpak command: flatpak update -y
	cmd := exec.Command("flatpak", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes unused Flatpak runtimes and cleans cache
func (f *Flatpak) Clean() error {
	// Clean unused runtimes and apps
	fmt.Println("🧹 Removing unused Flatpak runtimes...")
	uninstallCmd := exec.Command("flatpak", "uninstall", "--unused", "-y")
	uninstallCmd.Stdout = os.Stdout
	uninstallCmd.Stderr = os.Stderr
	err := uninstallCmd.Run()
	if err != nil {
		return err
	}

	// Repair installation
	fmt.Println("🔧 Repairing Flatpak installation...")
	repairCmd := exec.Command("flatpak", "repair", "--user")
	repairCmd.Stdout = os.Stdout
	repairCmd.Stderr = os.Stderr
	return repairCmd.Run()
}
