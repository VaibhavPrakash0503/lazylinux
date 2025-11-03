package pkgmgr

import (
	"bufio"
	"os"
	"strings"
)

func detectDistribution() string {
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		return parseOSRelease(string(data))
	}

	data, err = os.ReadFile("/etc/lsb-release")
	if err == nil {
		return parseLSBRelease(string(data))
	}

	data, err = os.ReadFile("/etc/fedora-release")
	if err == nil {
		return string(data)
	}

	_, err = os.ReadFile("/etc/arch-release")
	if err == nil {
		return "Arch Linux"
	}

	return "Unknown Linux Distribution"
}

func parseOSRelease(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if value, found := strings.CutPrefix(line, "PRETTY_NAME="); found {
			value = strings.Trim(value, "\"")
			return value
		}
	}
	return "Unknown Distribution"
}

func parseLSBRelease(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var distroName, release string

	for scanner.Scan() {
		line := scanner.Text()

		if value, found := strings.CutPrefix(line, "DISTRIB_ID="); found {
			distroName = value
		}
		if value, found := strings.CutPrefix(line, "DISTRIB_RELEASE="); found {
			release = value
		}
	}

	if distroName != "" && release != "" {
		return distroName + " " + release
	}
	return distroName
}
