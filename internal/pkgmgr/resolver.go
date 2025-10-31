package pkgmgr

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// PackageSource represents where a package was found
type PackageSource struct {
	Manager     string // "dnf", "apt", "pacman", "flatpak"
	PackageName string // The actual package name (might be different for Flatpak)
	Available   bool   // Whether it's available in this source
	Confidence  int    // Match confidence (0-100) - higher = better match
}

// ResolvePackage finds which package manager(s) have the package
func ResolvePackage(packageName string, nativePM PackageManager, hasFlatpak bool) []PackageSource {
	sources := []PackageSource{}

	// Check native package manager
	fmt.Printf("  🔍 Searching in %s...\n", getPackageManagerName(nativePM))
	nativeAvailable := checkNativePackage(packageName, nativePM)
	sources = append(sources, PackageSource{
		Manager:     getPackageManagerName(nativePM),
		PackageName: packageName,
		Available:   nativeAvailable,
		Confidence:  100, // Exact match in native
	})

	// Check Flatpak if available
	if hasFlatpak {
		fmt.Println("  🔍 Searching in Flatpak...")
		flatpakMatches := searchFlatpakPackages(packageName)
		sources = append(sources, flatpakMatches...)
	}

	return sources
}

// checkNativePackage checks if package exists in native package manager
func checkNativePackage(packageName string, pm PackageManager) bool {
	var cmd *exec.Cmd

	switch pm.(type) {
	case *DNF:
		// dnf repoquery <package> (quiet check)
		cmd = exec.Command("dnf", "repoquery", packageName)
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return false
		}
		return true
	case *APT:
		// apt-cache show <package>
		cmd = exec.Command("apt-cache", "search", "--namees-only", "^"+packageName+"$")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return false
		}
		return true
	case *Pacman:
		// pacman -Ss <package>
		cmd = exec.Command("pacman", "-Ss", "^"+packageName+"$")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return false
		}
		return true
	default:
		return false
	}
}

// searchFlatpakPackages searches Flatpak and returns all matching packages with confidence scores
func searchFlatpakPackages(packageName string) []PackageSource {
	matches := []PackageSource{}

	// Search Flatpak with proper column format
	// Format: --columns=application,name to get "app.id" and "App Name"
	cmd := exec.Command("flatpak", "search", "--columns=application,name", packageName)
	output, err := cmd.Output()

	if err != nil || len(output) == 0 {
		return matches
	}

	// Parse output line by line
	// Format: org.zen_browser.zen	Zen Browser
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split by tab
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		appID := strings.TrimSpace(parts[0])
		appName := strings.TrimSpace(parts[1])

		// Skip if appID looks invalid
		if appID == "" || !strings.Contains(appID, ".") {
			continue
		}

		// Calculate match confidence
		confidence := calculateMatchConfidence(packageName, appName, appID)

		matches = append(matches, PackageSource{
			Manager:     "flatpak",
			PackageName: appID,
			Available:   true,
			Confidence:  confidence,
		})
	}

	// Sort by confidence (best matches first)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Confidence > matches[i].Confidence {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	// Return only top 5 best matches (not all 25)
	if len(matches) > 5 {
		return matches[:5]
	}

	return matches
}

// calculateMatchConfidence calculates how well a Flatpak app matches the search term
func calculateMatchConfidence(searchTerm, appName, appID string) int {
	searchLower := strings.ToLower(searchTerm)
	appNameLower := strings.ToLower(appName)
	appIDLower := strings.ToLower(appID)

	// Exact match (highest priority)
	if searchLower == appNameLower {
		return 100
	}
	if searchLower == appIDLower {
		return 95
	}

	// Check if app name contains search term
	if strings.Contains(appNameLower, searchLower) {
		// Prefix match is better
		if strings.HasPrefix(appNameLower, searchLower) {
			return 90
		}
		return 80
	}

	// Check if search term appears in app ID (fuzzy matching with regex)
	confidence := fuzzyMatchAppID(searchLower, appIDLower)
	if confidence > 0 {
		return confidence
	}

	// Partial match in app name
	if strings.Contains(appNameLower, strings.Split(searchLower, " ")[0]) {
		return 50
	}

	// Very weak match
	return 10
}

// fuzzyMatchAppID checks how well search term matches Flatpak app ID
// Examples:
//
//	"zen" matches "org.zen_browser.zen" → 85
//	"browser" matches "org.zen_browser.zen" → 75
//	"spotify" matches "com.spotify.Client" → 80
func fuzzyMatchAppID(searchTerm, appID string) int {
	// Split app ID by dots and underscores
	// org.zen_browser.zen → ["org", "zen", "browser", "zen"]
	parts := strings.FieldsFunc(appID, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})

	for searchPart := range strings.SplitSeq(searchTerm, " ") {
		for _, part := range parts {
			partLower := strings.ToLower(part)

			// Exact match of a part
			if searchPart == partLower {
				return 85
			}

			// Prefix match of a part
			if strings.HasPrefix(partLower, searchPart) {
				return 75
			}

			// Contains match (weaker)
			if strings.Contains(partLower, searchPart) {
				return 60
			}
		}
	}

	return 0
}

// PromptUserChoice asks user to choose between multiple sources
func PromptUserChoice(sources []PackageSource, packageName string) *PackageSource {
	// Filter only available sources
	available := []PackageSource{}
	for _, src := range sources {
		if src.Available {
			available = append(available, src)
		}
	}

	// No sources available
	if len(available) == 0 {
		return nil
	}

	// Only one source - use it without asking
	if len(available) == 1 {
		src := available[0]
		if src.Manager == "flatpak" {
			fmt.Printf("✅ Found in Flatpak: %s\n", src.PackageName)
		} else {
			fmt.Printf("✅ Found in %s\n", src.Manager)
		}
		return &src
	}

	// Multiple sources - show them with confidence scores
	fmt.Printf("\n📦 Found '%s' in multiple sources:\n", packageName)

	for i, src := range available {

		displayName := src.PackageName
		if src.Manager == "flatpak" {
			// Show both app name and ID for Flatpak
			displayName = fmt.Sprintf("%s [%s]", src.PackageName, getConfidenceLabel(src.Confidence))
		}

		fmt.Printf("  [%d] %-10s - %s\n", i+1, src.Manager, displayName)
	}

	// Get user input
	fmt.Print("\nChoose source only one: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Default to 1 (first/best match)
	if input == "" {
		input = "1"
	}

	// Parse choice
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(available) {
		fmt.Println("❌ Invalid choice, using first option")
		return &available[0]
	}

	return &available[choice-1]
}

// getConfidenceLabel returns a label for match confidence
func getConfidenceLabel(confidence int) string {
	if confidence >= 90 {
		return "exact"
	} else if confidence >= 80 {
		return "strong"
	} else if confidence >= 70 {
		return "good"
	} else if confidence >= 60 {
		return "fair"
	}
	return "weak"
}

// GetBestMatch returns the best available source (auto-select if very confident)
func GetBestMatch(sources []PackageSource) *PackageSource {
	// Filter available
	available := []PackageSource{}
	for _, src := range sources {
		if src.Available {
			available = append(available, src)
		}
	}

	if len(available) == 0 {
		return nil
	}

	if len(available) == 1 {
		return &available[0]
	}

	// If best match is very high confidence and native, auto-select
	if available[0].Manager != "flatpak" && available[0].Confidence == 100 {
		return &available[0]
	}

	// If Flatpak match is very high confidence, ask user
	return nil // Force prompt
}
