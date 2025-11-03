package pkgmgr

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"lazylinux/internal/config"
)

type SourcePreferences struct {
	Flatpak bool
	RPM     bool // or AUR for Arch
}

// RunInit detects package manager and saves configuration
func RunInit() error {
	fmt.Println("🚀 Initializing LazyLinux...")
	fmt.Println()

	// Detect package manager
	fmt.Println("🔍 Detecting package manager...")
	pm, err := DetectPackageManager()
	if err != nil {
		return fmt.Errorf("detection failed: %v", err)
	}

	// Get package manager name
	pmName := getPackageManagerName(pm)
	fmt.Printf("✅ Detected: %s\n", pmName)
	fmt.Println()

	// Get user preferences
	prefs := setupSourcesMenu()

	// Actually enable the sources
	fmt.Println("\n⏳ Setting up package sources...")
	if err := enableSources(prefs); err != nil {
		fmt.Printf("⚠️  Some installations failed, but continuing...\n")
	}

	// Save preferences to config
	if err := saveSourcePreferences(prefs, pmName); err != nil {
		fmt.Println("❌ Error saving config:", err)
	}

	fmt.Printf("✅ Configuration saved to: %s\n", config.GetConfigPath())
	fmt.Println()
	fmt.Println("LazyLinux is ready to use! 🎉")
	fmt.Println()
	fmt.Println("Try these commands:")
	fmt.Println("  lazylinux install <package>")
	fmt.Println("  lazylinux remove <package>")
	fmt.Println("  lazylinux list")
	fmt.Println("	 lazylinux update")
	fmt.Println("	 lazylinux clean")

	return nil
}

func saveSourcePreferences(prefs SourcePreferences, pmName string) error {
	cfg := config.Config{
		PackageManager: pmName,
		FlatpakEnabled: prefs.Flatpak,
		RPMEnabled:     prefs.RPM,
	}

	return config.SaveConfig(&cfg)
}

// getPackageManagerName converts PackageManager to string
func getPackageManagerName(pm PackageManager) string {
	switch pm.(type) {
	case *DNF:
		return "dnf"
	case *APT:
		return "apt"
	case *Pacman:
		return "pacman"
	default:
		return "unknown"
	}
}

func setupSourcesMenu() SourcePreferences {
	prefs := SourcePreferences{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nPackages Source Configuration")
	fmt.Println("================================")
	fmt.Println("Select which sources to include (enter multiple numbers separated by space):")
	fmt.Println("1. Flatpak")
	fmt.Println("2. RPM/AUR")
	fmt.Println("0. Skip all")
	fmt.Print("\nEnter choices (e.g., 1 2 or 1 2): ")

	input, _ := reader.ReadString('\n')
	choices := strings.FieldsSeq(strings.TrimSpace(input))

	for choice := range choices {
		switch choice {
		case "1":
			prefs.Flatpak = true
		case "2":
			prefs.RPM = true
		case "0":
			return prefs
		}
	}

	return prefs
}

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

// Detect which Linux distribution is running
func detectDistribution() string {
	// Method 1: Check /etc/os-release (most reliable, works on all modern distros)
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		return parseOSRelease(string(data))
	}

	// Method 2: Check /etc/lsb-release (Ubuntu/Debian)
	data, err = os.ReadFile("/etc/lsb-release")
	if err == nil {
		return parseLSBRelease(string(data))
	}

	// Method 3: Check /etc/fedora-release (Fedora)
	data, err = os.ReadFile("/etc/fedora-release")
	if err == nil {
		return string(data)
	}

	// Method 4: Check /etc/arch-release (Arch)
	_, err = os.ReadFile("/etc/arch-release")
	if err == nil {
		return "Arch Linux"
	}

	return "Unknown Linux Distribution"
}

// Parse /etc/os-release file
func parseOSRelease(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if value, found := strings.CutPrefix(line, "PRETTY_NAME="); found {
			value = strings.Trim(value, "\"")
			return value
		}
	}
	return "Unknown Distribution"
}

// Parse /etc/lsb-release file
func parseLSBRelease(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var distroName, release string

	for scanner.Scan() {
		line := scanner.Text()

		if value, found := strings.CutPrefix(line, "DISTRIB_ID="); found {
			distroName = value
		}
		if value, found := strings.CutPrefix(line, "DISTRIB_RELEASE="); found {
			release = value
		}
	}

	if distroName != "" && release != "" {
		return fmt.Sprintf("%s %s", distroName, release)
	}
	return distroName
}

// Install Flatpak based on distro
func installFlatpak(distro string) error {
	var cmd *exec.Cmd

	switch {
	case strings.Contains(distro, "Ubuntu") || strings.Contains(distro, "Debian"):
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "flatpak")
	case strings.Contains(distro, "Fedora") || strings.Contains(distro, "RHEL"):
		cmd = exec.Command("sudo", "dnf", "install", "-y", "flatpak")
	case strings.Contains(distro, "Arch"):
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "flatpak")
	default:
		return fmt.Errorf("unsupported distribution: %s", distro)
	}

	return cmd.Run()
}

// Install RPM support based on distro
func installRPM(distro string) error {
	var cmd *exec.Cmd

	switch {
	case strings.Contains(distro, "Ubuntu") || strings.Contains(distro, "Debian"):
		// For Debian/Ubuntu, install alien to convert RPM to DEB
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "alien")
	case strings.Contains(distro, "Fedora") || strings.Contains(distro, "RHEL"):
		// RPM is already installed on Fedora/RHEL
		fmt.Println("RPM support is already available on this system")
		return nil
	case strings.Contains(distro, "Arch"):
		// For Arch, we'd enable AUR support (requires auracle or yay)
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "yay")
	default:
		return fmt.Errorf("unsupported distribution: %s", distro)
	}

	return cmd.Run()
}
