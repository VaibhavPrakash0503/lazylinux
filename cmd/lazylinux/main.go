package main

import (
	"fmt"
	"os"

	"lazylinux/internal/pkgmgr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage lazylinux <command> [option]")
		fmt.Println("Commands")
		fmt.Println("install <package>... -Install a package")
		fmt.Println("remove <package>...  - Remove packages")
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
		if len(os.Args) < 3 {
			fmt.Println("No package sepecified")
			os.Exit(1)
		}

		packages := os.Args[2:]

		pm := pkgmgr.NewDNF()

		fmt.Printf("Installing packages %v\n", packages)
		err := pm.Install(packages...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Install failed %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Install complete")

	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("No package sepecified")
			os.Exit(1)
		}

		packages := os.Args[2:]

		pm := pkgmgr.NewDNF()

		fmt.Printf("Removing packages %v\n", packages)
		err := pm.Remove(packages...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Removing failed %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Remove complete")
	default: // Add this!
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: install, remove")
		os.Exit(1)
	}
}
