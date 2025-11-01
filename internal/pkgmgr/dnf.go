package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

// List lists all installed packages
func (d *DNF) List() ([]string, error) {
	// Try rpm -qa first (more reliable)
	cmd := exec.Command("rpm", "-qa", "--queryformat", "%{NAME}\n")
	output, err := cmd.Output()

	// If rpm fails, try dnf list --installed
	if err != nil || len(output) == 0 {
		cmd = exec.Command("dnf", "list", "--installed")
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to list packages: %v", err)
		}
	}

	// Parse output - each line is a package name
	var results []string
	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Installed Packages") {
			continue
		}

		// Extract package name (before the space/architecture)
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			// Split by dot to remove architecture
			pkgName := strings.Split(parts[0], ".")[0]
			if pkgName != "" {
				results = append(results, pkgName)
			}
		}
	}

	return results, nil
}
