package main

import (
	"fmt"
	"os"

	"lazylinux/internal/config"
	"lazylinux/internal/pkgmgr"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		// Run initialization
		err := pkgmgr.RunInit()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Initialization failed: %v\n", err)
			os.Exit(1)
		}

	case "install":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		if len(os.Args) < 3 {
			fmt.Println("No package specified")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		packages := os.Args[2:]

		// Install each package with smart resolution
		for _, pkg := range packages {
			fmt.Printf("\n🔍 Looking for '%s'...\n", pkg)

			// Resolve package across all sources
			sources := pkgmgr.ResolvePackage(pkg, pm, cfg.FlatpakEnabled)

			// Let user choose or auto-select
			chosen := pkgmgr.PromptUserChoice(sources, pkg)

			if chosen == nil {
				fmt.Printf("❌ Package '%s' not found in any source\n", pkg)
				continue
			}

			// Install from chosen source
			fmt.Printf("\n📦 Installing '%s' from %s...\n", pkg, chosen.Manager)

			var installErr error
			if chosen.Manager == "flatpak" {
				flatpakPM := pkgmgr.NewFlatpak()
				installErr = flatpakPM.Install(chosen.PackageName)
			} else {
				installErr = pm.Install(pkg)
			}

			if installErr != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to install '%s': %v\n", pkg, installErr)
			} else {
				fmt.Printf("✅ Successfully installed '%s'\n", pkg)
			}
		}

	case "remove":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		if len(os.Args) < 3 {
			fmt.Println("Error: No package specified")
			fmt.Println("Usage: lazylinux remove <package>...")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		packages := os.Args[2:]

		// Remove each package with smart resolution
		for _, pkg := range packages {
			fmt.Printf("\n🔍 Looking for '%s' to remove...\n", pkg)

			// Resolve package across all sources
			sources := pkgmgr.ResolvePackage(pkg, pm, cfg.FlatpakEnabled)

			// Let user choose or auto-select
			chosen := pkgmgr.PromptUserChoice(sources, pkg)

			if chosen == nil {
				fmt.Printf("❌ Package '%s' not found in any source\n", pkg)
				continue
			}

			// Remove from chosen source
			fmt.Printf("\n📦 Removing '%s' from %s...\n", pkg, chosen.Manager)

			var removeErr error
			if chosen.Manager == "flatpak" {
				flatpakPM := pkgmgr.NewFlatpak()
				removeErr = flatpakPM.Remove(chosen.PackageName)
			} else {
				removeErr = pm.Remove(pkg)
			}

			if removeErr != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to remove '%s': %v\n", pkg, removeErr)
			} else {
				fmt.Printf("✅ Successfully removed '%s'\n", pkg)
			}
		}

	case "update":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🔄 Updating packages...")
		fmt.Println()

		// Update native package manager
		fmt.Printf("🔄 Updating %s packages...\n", getPackageManagerName(pm))
		err = pm.Update()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to update native packages: %v\n", err)
		} else {
			fmt.Println("✅ Native packages updated")
		}

		// Update Flatpak if available
		if cfg.FlatpakEnabled {
			fmt.Println()
			fmt.Println("🔄 Updating Flatpak packages...")
			flatpakPM := pkgmgr.NewFlatpak()
			err = flatpakPM.Update()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to update Flatpak packages: %v\n", err)
			} else {
				fmt.Println("✅ Flatpak packages updated")
			}
		}

		fmt.Println()
		fmt.Println("✅ All updates complete!")

	case "clean":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🧼 Cleaning system...")
		fmt.Println()

		// Clean native package manager
		fmt.Printf("🧹 Cleaning %s...\n", getPackageManagerName(pm))
		err = pm.Clean()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to clean native packages: %v\n", err)
		} else {
			fmt.Println("✅ Native packages cleaned")
		}

		// Clean Flatpak if available
		if cfg.FlatpakEnabled {
			fmt.Println()
			fmt.Println("🧹 Cleaning Flatpak...")
			flatpakPM := pkgmgr.NewFlatpak()
			err = flatpakPM.Clean()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to clean Flatpak: %v\n", err)
			} else {
				fmt.Println("✅ Flatpak cleaned")
			}
		}

		fmt.Println()
		fmt.Println("✅ System cleaned!")

	case "list":
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📋 Installed Packages")

		// List native packages
		fmt.Printf("📦 %s Packages:\n", getPackageManagerName(pm))
		nativeList, err := pm.List()
		if err == nil && len(nativeList) > 0 {
			for i, pkg := range nativeList {
				if i >= 20 {
					fmt.Printf("  ... and %d more\n", len(nativeList)-20)
					break
				}
				fmt.Printf("  • %s\n", pkg)
			}
		} else {
			fmt.Println("  (none found)")
		}
		fmt.Println()

		// List Flatpak packages if available
		if cfg.FlatpakEnabled {
			fmt.Println("🎨 Flatpak Packages:")
			flatpakPM := pkgmgr.NewFlatpak()
			flatpakList, err := flatpakPM.List()
			if err == nil && len(flatpakList) > 0 {
				for i, pkg := range flatpakList {
					if i >= 20 {
						fmt.Printf("  ... and %d more\n", len(flatpakList)-20)
						break
					}
					fmt.Printf("  • %s\n", pkg)
				}
			} else {
				fmt.Println("  (none found)")
			}
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

// loadPackageManager loads the package manager from config
func loadPackageManager() (pkgmgr.PackageManager, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load config: %v", err)
	}

	switch cfg.PackageManager {
	case "dnf":
		return pkgmgr.NewDNF(), nil
	case "apt":
		return pkgmgr.NewAPT(), nil
	case "pacman":
		return pkgmgr.NewPacman(), nil
	default:
		return nil, fmt.Errorf("unknown package manager: %s", cfg.PackageManager)
	}
}

func showHelp() {
	fmt.Println("Usage: lazylinux <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init                   - Initialize LazyLinux (run this first)")
	fmt.Println("  install <package>...   - Install packages")
	fmt.Println("  remove <package>...    - Remove packages")
	fmt.Println("  update                 - Update all packages")
	fmt.Println("  clean                  - Clean cache and remove orphaned packages")
	fmt.Println("  list                   - List all the installed packages")
}

func getPackageManagerName(pm pkgmgr.PackageManager) string {
	switch pm.(type) {
	case *pkgmgr.DNF:
		return "DNF"
	case *pkgmgr.APT:
		return "APT"
	case *pkgmgr.Pacman:
		return "Pacman"
	default:
		return "unknown"
	}
}
