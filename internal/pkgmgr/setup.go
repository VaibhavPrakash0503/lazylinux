package pkgmgr

import (
	"fmt"

	"github.com/VaibhavPrakash0503/lazylinux/internal/config"
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

	pmName := getPackageManagerName(pm)
	fmt.Printf("✅ Detected: %s\n", pmName)
	fmt.Println()

	// Get user preferences with auto-detection
	prefs := setupSourcesMenu()

	// Enable the selected sources
	fmt.Println("\n⏳ Setting up package sources...")
	if err := enableSources(prefs); err != nil {
		fmt.Printf("⚠️  Some installations failed, but continuing...\n")
	}

	// Save preferences to config
	if err := saveSourcePreferences(pmName); err != nil {
		fmt.Println("❌ Error saving config:", err)
	}

	fmt.Println("LazyLinux is ready to use! 🎉")
	fmt.Println()
	fmt.Println("Try these commands:")
	fmt.Println("  lazylinux install <package>")
	fmt.Println("  lazylinux remove <package>")
	fmt.Println("  lazylinux search <package>")
	fmt.Println("  lazylinux update")

	return nil
}

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

func saveSourcePreferences(pmName string) error {
	cfg := &config.Config{
		PackageManager: pmName,
		FlatpakEnabled: isFlatpakInstalled(),
		SnapEnabled:    isSnapInstalled(),
		RPMEnabled:     isRPMInstalled(),
	}

	return config.SaveConfig(cfg)
}
