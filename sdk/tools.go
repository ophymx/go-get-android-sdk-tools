package sdk

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var licenseRegex = regexp.MustCompile(`Do you accept the license '(.*-license-.*)' \[y/n\]`)

// Tools wrapper for android sdk manager
type Tools struct {
	path          string
	licenses      map[string]bool
	alwaysInstall map[string]bool
}

func NewTools(path string, licenses, alwaysInstall []string) Tools {
	licensesMap := map[string]bool{}
	for _, license := range licenses {
		licensesMap[license] = true
	}

	alwaysInstallMap := map[string]bool{}
	for _, archive := range alwaysInstall {
		alwaysInstallMap[archive] = true
	}

	return Tools{
		path:          path,
		licenses:      licensesMap,
		alwaysInstall: alwaysInstallMap,
	}
}

// IsToolsInstalled test that a directory is a valid Android SDK Tools directory
func (t Tools) IsToolsInstalled() bool {
	info, err := os.Stat(t.androidPath())
	if err != nil {
		return false
	}
	return info.Mode().Perm()&0111 != 0
}

func (t Tools) IsInstalled(pkgDir string) bool {
	_, err := os.Stat(pkgDir)
	return err == nil
}

func (t Tools) BuildToolsDir(version string) string {
	return t.join("build-tools", version)
}

func (t Tools) PlatformToolsDir() string {
	return t.join("platform-tools")
}

func (t Tools) CompileSdkDir(version string) string {
	return t.join("platforms", "android-"+version)
}

func (t Tools) VersionedDir(prefix, version string) string {
	return t.join(filepath.FromSlash(prefix) + version)
}

func (t Tools) InstallDir(path string) string {
	return t.join(filepath.FromSlash(path))
}

func (t Tools) Install(name string) (err error) {
	log.Print("Installing ", name)
	return t.acceptLicense(t.android("update", "sdk", "-u", "-a", "-t", name))
}

func (t Tools) InstallVersion(name, version string) (err error) {
	return t.Install(name + "-" + version)
}

func (t Tools) Update(name string) (err error) {
	log.Print("Updating ", name)
	if t.alwaysInstall[name] {
		return t.Install(name)
	}

	if err = t.acceptLicense(t.android("update", "sdk", "-u", "-t", name)); err == io.EOF {
		return nil
	}
	return
}

func (t Tools) UpdateVersion(name, version string) (err error) {
	return t.Update(name + "-" + version)
}

func (t Tools) acceptLicense(cmd *exec.Cmd) (err error) {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	out := bufio.NewReader(stdout)
	if err = cmd.Start(); err != nil {
		return
	}

	var license string
	var bytes []byte
	for _, err = out.ReadString('\n'); err == nil; _, err = out.ReadString('\n') {
		bytes, err = out.Peek(27)
		if err != nil {
			return
		}

		if string(bytes) == "Do you accept the license '" {
			if _, err = out.ReadString('\''); err != nil {
				return
			}
			license, err = out.ReadString('\'')
			if err != nil {
				return
			}
			license = strings.TrimSuffix(license, "'")
			if _, err = out.ReadString(':'); err != nil {
				return
			}
			if _, err = out.ReadByte(); err != nil {
				return
			}
			break
		}
	}
	if err != nil {
		return
	}
	if t.licenses[license] {
		fmt.Println("Accepting license: ", license)
		stdin.Write([]byte("y\r\n"))
	} else {
		fmt.Println("Rejecting license: ", license)
		stdin.Write([]byte("n\r\n"))
	}

	go out.WriteTo(os.Stdout)
	err = cmd.Wait()

	return
}

func (t Tools) android(args ...string) *exec.Cmd {
	return exec.Command(t.androidPath(), args...)
}

func (t Tools) androidPath() string {
	return t.join("tools", "android")
}

func (t Tools) join(parts ...string) string {
	return filepath.Join(t.path, filepath.Join(parts...))
}
