// Package pkgmgr provides an abstraction layer for different Linux package managers.
// It supports DNF (Fedora/RHEL), APT (Ubuntu/Debian), and Pacman (Arch/Manjaro).
package pkgmgr

// PackageManager defines operations all package managers must support
type PackageManager interface {
	Install(packages ...string) error
	Remove(packages ...string) error
	Update() error
	Clean() error
	List() ([]string, error)
}
