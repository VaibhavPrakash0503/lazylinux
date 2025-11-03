package main

import (
	"fmt"
	"os"
	"sync"

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
		requireInit()
		handlePackageOperation("install")

	case "remove":
		// Check if initialized
		requireInit()
		handlePackageOperation("remove")

	case "update":
		// Check if initialized
		requireInit()
		update()

	case "clean":
		// Check if initialized
		requireInit()
		clean()

	case "list":
		requireInit()
		list()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func requireInit() {
	if !config.ConfigExists() {
		fmt.Println("❌ LazyLinux not initialized!")
		fmt.Println("Run: lazylinux init")
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

func handlePackageOperation(operation string) {
	if len(os.Args) < 3 {
		fmt.Printf("Error: No package specified\n")
		fmt.Printf("Usage: lazylinux %s <package>...\n", operation)
		os.Exit(1)
	}

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
	isRemove := operation == "remove"

	for _, pkg := range packages {
		action := "Looking for"
		if isRemove {
			action = "Looking for to remove"
		}
		fmt.Printf("\n🔍 %s '%s'...\n", action, pkg)

		var sources []pkgmgr.PackageSource
		if isRemove {
			sources = pkgmgr.ResolvePackageForRemove(pkg, pm, cfg.FlatpakEnabled)
		} else {
			sources = pkgmgr.ResolvePackage(pkg, pm, cfg.FlatpakEnabled)
		}

		chosen := pkgmgr.PromptUserChoice(sources, pkg)
		if chosen == nil {
			fmt.Printf("❌ Package '%s' not found in any source\n", pkg)
			continue
		}

		verb := map[string]string{"remove": "Removing", "install": "Installing"}[operation]
		pastVerb := map[string]string{"remove": "removed", "install": "installed"}[operation]

		fmt.Printf("\n📦 %s '%s' from %s...\n", verb, pkg, chosen.Manager)

		var err error
		if chosen.Manager == "flatpak" {
			flatpakPM := pkgmgr.NewFlatpak()
			if isRemove {
				err = flatpakPM.Remove(chosen.PackageName)
			} else {
				err = flatpakPM.Install(chosen.PackageName)
			}
		} else {
			if isRemove {
				err = pm.Remove(pkg)
			} else {
				err = pm.Install(pkg)
			}
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to %s '%s': %v\n", pastVerb, pkg, err)
		} else {
			fmt.Printf("✅ Successfully %s '%s'\n", pastVerb, pkg)
		}
	}
}

func update() {
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

	var wg sync.WaitGroup
	errChan := make(chan struct {
		manager string
		err     error
	}, 2) // Buffer for max 2 sources

	numUpdates := 1
	if cfg.FlatpakEnabled {
		numUpdates = 2
	}
	wg.Add(numUpdates)

	// Update native package manager (goroutine)
	go func() {
		defer wg.Done()
		fmt.Printf("🔄 Updating %s packages...\n", getPackageManagerName(pm))
		err := pm.Update()
		if err != nil {
			errChan <- struct {
				manager string
				err     error
			}{getPackageManagerName(pm), err}
			return
		}
		fmt.Printf("✅ %s packages updated\n", getPackageManagerName(pm))
	}()

	// Update Flatpak (goroutine)
	if cfg.FlatpakEnabled {
		go func() {
			defer wg.Done()
			fmt.Println("🔄 Updating Flatpak packages...")
			flatpakPM := pkgmgr.NewFlatpak()
			err := flatpakPM.Update()
			if err != nil {
				errChan <- struct {
					manager string
					err     error
				}{"Flatpak", err}
				return
			}
			fmt.Println("✅ Flatpak packages updated")
		}()
	}

	// Wait for all updates to complete
	wg.Wait()
	close(errChan)

	// Handle any errors that occurred
	hasErrors := false
	for errResult := range errChan {
		hasErrors = true
		fmt.Fprintf(os.Stderr, "❌ Failed to update %s packages: %v\n", errResult.manager, errResult.err)
	}

	fmt.Println()
	if !hasErrors {
		fmt.Println("✅ All updates complete!")
	} else {
		fmt.Println("⚠️  Some updates failed. Check errors above.")
	}
}

func clean() {
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
}

func list() {
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
}
