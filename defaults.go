package main

import (
	"fmt"

	"github.com/ophymx/go-get-android-sdk-tools/sdk"
)

const (
	defaultDownloadPrefix = "http://dl.google.com/android"
	defaultVersion        = "24.2"
)

var defaultSdkVersion = &sdk.VersionInfo{
	Version:   defaultVersion,
	OS:        defaultOs,
	Extension: defaultExt,
	Size:      defaultSize,
	SHA1:      defaultSha1,
	URL: fmt.Sprintf(
		"%s/android-sdk_r%s-%s.%s",
		defaultDownloadPrefix,
		defaultVersion,
		defaultOs,
		defaultExt,
	),
}

var defaultConfig = config{
	AcceptedLicenses: []string{"android-sdk-license-5be876d5"},
	Archives: map[string]string{
		"build-tools-22.0.1":        "build-tools/22.0.1",
		"android-22":                "platforms/android-22",
		"extra-google-m2repository": "extras/google/m2repository",
	},
}
