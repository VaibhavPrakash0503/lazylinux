package pkgmgr

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SourcePreferences struct {
	Flatpak bool
	Snap    bool
	RPM     bool
}

func setupSourcesMenu() SourcePreferences {
	prefs := SourcePreferences{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n🔍 Detecting installed packages...")
	flatpakInstalled := isFlatpakInstalled()
	snapInstalled := isSnapInstalled()
	rpmInstalled := isRPMInstalled()

	missing := []string{}
	fmt.Println("\n📦 Current system status:")
	fmt.Println("================================")

	if flatpakInstalled {
		fmt.Println("  ✅ Flatpak is installed")
	} else {
		fmt.Println("  ❌ Flatpak is NOT installed")
		missing = append(missing, "1. Flatpak")
	}

	if snapInstalled {
		fmt.Println("  ✅ Snap is installed")
	} else {
		fmt.Println("  ❌ Snap is NOT installed")
		missing = append(missing, "2. Snap")
	}

	if rpmInstalled {
		fmt.Println("  ✅ RPM is available")
	} else {
		fmt.Println("  ❌ RPM is NOT available")
		missing = append(missing, "3. RPM/AUR")
	}

	fmt.Println()

	// If nothing is missing, we're done
	if len(missing) == 0 {
		fmt.Println("✅ All optional packages are already installed!")
		return prefs
	}

	// Ask user to install missing packages
	fmt.Println("📝 Packages to install:")
	for _, pkg := range missing {
		fmt.Println("  " + pkg)
	}
	fmt.Println()

	fmt.Print("Install these packages? (Y/n): ")
	input, _ := reader.ReadString('\n')
	response := strings.TrimSpace(input)

	if response != "y" && response != "Y" && response != "" {
		fmt.Println("Skipping installation")
		return prefs
	}

	// Ask which ones to install
	fmt.Println()
	fmt.Println("Select which packages to install:")
	for _, pkg := range missing {
		fmt.Println("  " + pkg)
	}
	fmt.Print("\nEnter choices (e.g., 1 2 3 or press Enter to install all): ")

	input, _ = reader.ReadString('\n')
	choices := strings.Fields(strings.TrimSpace(input))

	// If empty, select all
	if len(choices) == 0 {
		if !flatpakInstalled {
			prefs.Flatpak = true
		}
		if !snapInstalled {
			prefs.Snap = true
		}
		if !rpmInstalled {
			prefs.RPM = true
		}
		return prefs
	}

	// Select only chosen ones
	for _, choice := range choices {
		switch choice {
		case "1":
			if !flatpakInstalled {
				prefs.Flatpak = true
			}
		case "2":
			if !snapInstalled {
				prefs.Snap = true
			}
		case "3":
			if !rpmInstalled {
				prefs.RPM = true
			}
		}
	}

	return prefs
}
