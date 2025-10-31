package pkgmgr

import (
	"fmt"

	"lazylinux/internal/config"
)

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

	// Detect Flatpak
	flatpakAvailable := DetectFlatpak()
	if flatpakAvailable {
		fmt.Println("  ✅ Flatpak: available")
	} else {
		fmt.Println("  ℹ️  Flatpak: not installed (optional)")
	}
	fmt.Println()

	// Save to config
	fmt.Println("💾 Saving configuration...")
	cfg := &config.Config{
		PackageManager: pmName,
		FlatpakEnabled: flatpakAvailable,
	}

	err = config.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("✅ Configuration saved to: %s\n", config.GetConfigPath())
	fmt.Println()
	fmt.Println("LazyLinux is ready to use! 🎉")
	fmt.Println()
	fmt.Println("Try these commands:")
	fmt.Println("  lazylinux install <package>")
	fmt.Println("  lazylinux remove <package>")

	return nil
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
