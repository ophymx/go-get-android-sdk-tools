package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ophymx/go-get-android-sdk-tools/sdk"
)

var (
	appName    = filepath.Base(os.Args[0])
	appVersion = "0.0.1"
)

func main() {
	if len(os.Args) != 3 {
		checkErr(fmt.Errorf("Usage: %s SDK_INSTALL_PATH CONFIG_FILE", appName))
	}

	toolsPath := os.Args[1]
	configPath := os.Args[2]

	config, err := readConfig(configPath)
	checkErr(err)
	checkErr(run(toolsPath, config))

}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err.Error())
		os.Exit(1)
	}
}

func run(toolsPath string, config config) (err error) {
	tools := sdk.NewTools(toolsPath, config.AcceptedLicenses)

	if !tools.IsToolsInstalled() {
		if err = installTools(toolsPath); err != nil {
			return
		}
	}
	if !tools.IsInstalled(tools.PlatformToolsDir()) {
		if err = tools.Install("platform-tools"); err != nil {
			return
		}
	} else {
		if err = tools.Update("platform-tools"); err != nil {
			return
		}
	}
	if err = tools.Update("tools"); err != nil {
		return
	}

	for name, installPath := range config.Archives {
		if !tools.IsInstalled(tools.InstallDir(installPath)) {
			if err = tools.Install(name); err != nil {
				return
			}
		} else {
			if err = tools.Update(name); err != nil {
				return
			}
		}
	}
	return
}

func installTools(path string) (err error) {
	version, err := selectVersion()
	if err != nil {
		return
	}
	renderer := &simpleProgress{
		name:  "Android SDK Tools " + version.Version,
		width: 60,
	}
	return sdk.InstallToolsProgress(version, path, renderer)
}

func selectVersion() (*sdk.VersionInfo, error) {
	versions, err := sdk.LatestVersions()
	if err != nil {
		return nil, err
	}
	for _, version := range versions {
		if version.OS == defaultOs {
			return version, nil
		}
	}
	log.Printf("Latest version not found, fallback to default:\n    %v", defaultSdkVersion)
	return defaultSdkVersion, nil
}
