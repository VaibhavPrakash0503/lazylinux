package pkgmgr

import (
	"fmt"
	"os/exec"
)

// DetectPackageManager automatically detects which package manager to use
func DetectPackageManager() (PackageManager, error) {
	// Check for DNF (Fedora, RHEL, CentOS)
	if _, err := exec.LookPath("dnf"); err == nil {
		return NewDNF(), nil
	}

	// Check for APT (Ubuntu, Debian)
	if _, err := exec.LookPath("apt"); err == nil {
		return NewAPT(), nil
	}

	// Check for Pacman (Arch, Manjaro)
	if _, err := exec.LookPath("pacman"); err == nil {
		return NewPacman(), nil
	}

	return nil, fmt.Errorf("no supported package manager found (dnf, apt, pacman)")
}
