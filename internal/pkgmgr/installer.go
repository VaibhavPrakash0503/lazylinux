package pkgmgr

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func enableSources(prefs SourcePreferences) error {
	distro := detectDistribution()

	if prefs.Flatpak {
		fmt.Print("📦 Installing Flatpak... ")
		if err := installFlatpak(distro); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return err
		}
		fmt.Println("✅ Done")
	}

	if prefs.Snap {
		fmt.Print("📦 Installing Snap... ")
		if err := installSnap(distro); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return err
		}
		fmt.Println("✅ Done")
	}

	if prefs.RPM {
		fmt.Print("📦 Installing RPM support... ")
		if err := installRPM(distro); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return err
		}
		fmt.Println("✅ Done")
	}

	return nil
}

func installFlatpak(distro string) error {
	var cmd *exec.Cmd
	distro = strings.ToLower(distro)

	switch {
	case strings.Contains(distro, "ubuntu") || strings.Contains(distro, "debian"):
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "flatpak")
	case strings.Contains(distro, "fedora"):
		cmd = exec.Command("sudo", "dnf", "install", "-y", "flatpak")
	case strings.Contains(distro, "arch"):
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "flatpak")
	default:
		return fmt.Errorf("unsupported distribution: %s", distro)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installSnap(distro string) error {
	var cmd *exec.Cmd
	distro = strings.ToLower(distro)

	switch {
	case strings.Contains(distro, "ubuntu") || strings.Contains(distro, "debian"):
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "snapd")
	case strings.Contains(distro, "fedora"):
		cmd = exec.Command("sudo", "dnf", "install", "-y", "snapd")
	case strings.Contains(distro, "arch"):
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "snapd")
	default:
		return fmt.Errorf("unsupported distribution: %s", distro)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installRPM(distro string) error {
	var cmd *exec.Cmd
	distro = strings.ToLower(distro)

	switch {
	case strings.Contains(distro, "fedora"):
		fmt.Println("📦 Enabling RPM Fusion repositories... ")
		cmd = exec.Command("sudo", "dnf", "install", "-y",
			"https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm",
			"https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm")
		return cmd.Run()

	case strings.Contains(distro, "rhel") || strings.Contains(distro, "centos"):
		cmd = exec.Command("sudo", "dnf", "install", "-y", "--nogpgcheck",
			"https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %rhel).noarch.rpm",
			"https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-$(rpm -E %rhel).noarch.rpm",
			"https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-$(rpm -E %rhel).noarch.rpm")
		return cmd.Run()

	case strings.Contains(distro, "ubuntu") || strings.Contains(distro, "debian"):
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "alien")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case strings.Contains(distro, "arch"):
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "yay")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	default:
		return fmt.Errorf("unsupported distribution: %s", distro)
	}
}
