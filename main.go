package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"lazylinux/internal/config"
	"lazylinux/internal/pkgmgr"
	"lazylinux/internal/webapp"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit()
	case "install":
		handleInstall()
	case "remove":
		handleRemove()
	case "update":
		handleUpdate()
	case "clean":
		handleClean()
	case "list":
		handleList()
	case "webapp":
		if len(os.Args) < 3 {
			showWebAppHelp()
			os.Exit(1)
		}
		handleWebApp(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

// Helper function to check if initialized
func mustBeInitialized() {
	if !config.ConfigExists() {
		fmt.Println("❌ LazyLinux not initialized!")
		fmt.Println("Run: lazylinux init")
		os.Exit(1)
	}
}

// Helper function to load config and package manager
func loadConfigAndPM() (*config.Config, pkgmgr.PackageManager, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("could not load config: %v", err)
	}

	pm, err := loadPackageManager()
	if err != nil {
		return nil, nil, err
	}

	return cfg, pm, nil
}

func handleInit() {
	err := pkgmgr.RunInit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Initialization failed: %v\n", err)
		os.Exit(1)
	}
}

func handleInstall() {
	mustBeInitialized()

	if len(os.Args) < 3 {
		fmt.Println("Error: No package specified")
		fmt.Println("Usage: lazylinux install <package>...")
		os.Exit(1)
	}

	cfg, pm, err := loadConfigAndPM()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	packages := os.Args[2:]
	for _, pkg := range packages {
		installPackage(pkg, pm, cfg)
	}
}

func handleRemove() {
	mustBeInitialized()

	if len(os.Args) < 3 {
		fmt.Println("Error: No package specified")
		fmt.Println("Usage: lazylinux remove <package>...")
		os.Exit(1)
	}

	cfg, pm, err := loadConfigAndPM()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	packages := os.Args[2:]
	for _, pkg := range packages {
		removePackage(pkg, pm, cfg)
	}
}

func installPackage(pkg string, pm pkgmgr.PackageManager, cfg *config.Config) {
	fmt.Printf("\n🔍 Looking for '%s'...\n", pkg)

	sources := pkgmgr.ResolvePackage(pkg, pm, cfg.FlatpakEnabled)
	chosen := pkgmgr.PromptUserChoice(sources, pkg)

	if chosen == nil {
		fmt.Printf("❌ Package '%s' not found in any source\n", pkg)
		return
	}

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

func removePackage(pkg string, pm pkgmgr.PackageManager, cfg *config.Config) {
	fmt.Printf("\n🔍 Looking for '%s' to remove...\n", pkg)

	sources := pkgmgr.ResolvePackageForRemove(pkg, pm, cfg.FlatpakEnabled)
	chosen := pkgmgr.PromptUserChoice(sources, pkg)

	if chosen == nil {
		fmt.Printf("❌ Package '%s' not found in any source\n", pkg)
		return
	}

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

func handleUpdate() {
	mustBeInitialized()

	cfg, pm, err := loadConfigAndPM()
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
	}, 2)

	numUpdates := 1
	if cfg.FlatpakEnabled {
		numUpdates = 2
	}
	wg.Add(numUpdates)

	// Update native package manager
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

	// Update Flatpak if enabled
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

	wg.Wait()
	close(errChan)

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

func handleClean() {
	mustBeInitialized()

	cfg, pm, err := loadConfigAndPM()
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

	// Clean Flatpak if enabled
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

func handleList() {
	mustBeInitialized()

	cfg, pm, err := loadConfigAndPM()
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

	// List Flatpak packages if enabled
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

func handleWebApp(args []string) {
	webappCmd := flag.NewFlagSet("webapp", flag.ExitOnError)
	add := webappCmd.Bool("a", false, "Add")
	del := webappCmd.Bool("d", false, "Delete")
	edit := webappCmd.Bool("e", false, "Edit")
	list := webappCmd.Bool("l", false, "List")
	open := webappCmd.Bool("o", false, "Open")

	err := webappCmd.Parse(args) // Check the error
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	remainingArgs := webappCmd.Args()

	switch {
	case *add:
		if len(remainingArgs) < 2 {
			fmt.Println("Usage: lazylinux webapp -a <name> <url>")
			os.Exit(1)
		}
		err := webapp.CreateWebApp(remainingArgs[0], remainingArgs[1])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ WebApp '%s' added\n", remainingArgs[0])

	case *del:
		if len(remainingArgs) < 1 {
			fmt.Println("Usage: lazylinux webapp -d <name>")
			os.Exit(1)
		}
		err := webapp.DeleteWebApp(remainingArgs[0])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ WebApp '%s' deleted\n", remainingArgs[0])

	case *edit:
		if len(remainingArgs) < 2 {
			fmt.Println("Usage: lazylinux webapp -e <name> <url>")
			os.Exit(1)
		}
		err := webapp.EditWebApp(remainingArgs[0], remainingArgs[1])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ WebApp '%s' updated\n", remainingArgs[0])

	case *list:
		apps, err := webapp.ListWebApp()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		if len(apps) == 0 {
			fmt.Println("📦 No webapps found")
			return
		}
		fmt.Println("📦 Your WebApps:")
		for i, app := range apps {
			fmt.Printf("%d. %s → %s\n", i+1, app.Name, app.URL)
		}

	case *open:
		if len(remainingArgs) < 1 {
			fmt.Println("Usage: lazylinux webapp -o <name>")
			os.Exit(1)
		}
		err := webapp.OpenWebApp(remainingArgs[0])
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	default:
		showWebAppHelp()
		os.Exit(1)
	}
}

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

func showHelp() {
	fmt.Println("Usage: lazylinux <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init                   - Initialize LazyLinux (run this first)")
	fmt.Println("  install <package>...   - Install packages")
	fmt.Println("  remove <package>...    - Remove packages")
	fmt.Println("  update                 - Update all packages")
	fmt.Println("  clean                  - Clean cache and remove orphaned packages")
	fmt.Println("  list                   - List all installed packages")
	fmt.Println("  webapp                 - Manage web applications")
}

func showWebAppHelp() {
	fmt.Println("Usage: lazylinux webapp [options]")
	fmt.Println("Options:")
	fmt.Println("  -a <name> <url>       Add new webapp")
	fmt.Println("  -d <name>             Delete webapp")
	fmt.Println("  -e <name> <url>       Edit webapp")
	fmt.Println("  -l                    List all webapps")
	fmt.Println("  -o <name>             Open a webapp")
}
