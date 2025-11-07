// Package webapp provides functionality for managing and launching web applications.
// It handles creating, editing, deleting, and listing web apps stored in YAML config.
// Web applications are launched in Chromium using --app mode for a dedicated window.
// Each webapp can have custom window size settings.
package webapp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var configPath = filepath.Join(os.ExpandEnv("$HOME"), ".lazylinux", "webapps.yaml")

func CreateWebApp(name, url string) error {
	apps, _ := LoadWebApps()

	for _, app := range apps.Apps {
		if app.Name == name {
			return fmt.Errorf("webapp %s already exists", name)
		}
	}

	apps.Apps = append(apps.Apps, WebApp{Name: name, URL: url})

	return SaveWebApps(apps)
}

func EditWebApp(name, newURL string) error {
	apps, err := LoadWebApps()
	if err != nil {
		return err
	}

	for i, app := range apps.Apps {
		if app.Name == name {
			apps.Apps[i].URL = newURL
			return SaveWebApps(apps)
		}
	}

	return fmt.Errorf("webapp '%s' not found", name)
}

func DeleteWebApp(name string) error {
	apps, err := LoadWebApps()
	if err != nil {
		return err
	}

	for i, app := range apps.Apps {
		if app.Name == name {
			apps.Apps = append(apps.Apps[:i], apps.Apps[i+1:]...)
			return SaveWebApps(apps)
		}
	}

	return fmt.Errorf("webapp '%s' not found", name)
}

func OpenWebApp(name string) error {
	apps, err := LoadWebApps()
	if err != nil {
		return err
	}

	for _, app := range apps.Apps {
		if app.Name == name {
			fmt.Printf("🚀 Opening %s: %s\n", app.Name, app.URL)
			cmd := exec.Command("chromium-browser", "--app="+app.URL, "--window-size=1600,900")
			return cmd.Start()
		}
	}
	return fmt.Errorf("webapp '%s' not found", name)
}

func LoadWebApps() (WebAppConfig, error) {
	var config WebAppConfig

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

func SaveWebApps(config WebAppConfig) error {
	dir := filepath.Dir(configPath)
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}

func ListWebApp() ([]WebApp, error) {
	app, err := LoadWebApps()

	return app.Apps, err
}
