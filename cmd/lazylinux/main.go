package main

import (
	"fmt"
	"os"

	"lazylinux/internal/config"
	"lazylinux/internal/pkgmgr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage lazylinux <command> [option]")
		fmt.Println("Commands")
		fmt.Println("install <package>... -Install a package")
		fmt.Println("remove <package>...  - Remove packages")
		fmt.Println("update - Updates system")
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
			fmt.Println("No package sepecified")
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		packages := os.Args[2:]
		fmt.Printf("Installing packages: %v\n", packages)
		err = pm.Install(packages...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Install failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Installation complete!")

	case "remove":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		if len(os.Args) < 3 {
			fmt.Println("No package sepecified")
			os.Exit(1)
		}

		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		packages := os.Args[2:]
		fmt.Printf("Removing packages: %v\n", packages)
		err = pm.Remove(packages...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Remove failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Remove complete!")

	case "update":
		// Check if initialized
		if !config.ConfigExists() {
			fmt.Println("❌ LazyLinux not initialized!")
			fmt.Println("Run: lazylinux init")
			os.Exit(1)
		}

		// Load package manager
		pm, err := loadPackageManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🔄 Updating all packages...")
		err = pm.Update()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Update failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Update complete!")
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
	fmt.Println("  update                 - Update all packages") // ← Add this
}
