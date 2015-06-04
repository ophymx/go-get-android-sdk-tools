package sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// InstallTools install sdk manager to path
func InstallTools(info *VersionInfo, path string) (err error) {
	return InstallToolsProgress(info, path, nullRenderer{})
}

// InstallToolsProgress install sdk manager to path and render progress
func InstallToolsProgress(info *VersionInfo, path string, renderer ProgressRenderer) (err error) {
	resp, err := http.Get(info.URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	temp, err := tempFile()
	if err != nil {
		return
	}
	defer temp.Close()
	defer os.Remove(temp.Name())

	digest := sha1.New()
	progress := &progressWriter{
		renderer: renderer,
		total:    info.Size,
	}
	progress.renderer.Init()
	progress.renderer.StartPhase("Downloading...")
	copied, err := io.Copy(io.MultiWriter(temp, digest, progress), resp.Body)
	if err != nil {
		return
	}
	sha1 := hex.EncodeToString(digest.Sum(nil))
	progress.renderer.Complete()

	if copied != info.Size {
		return fmt.Errorf(
			"Error downloading '%s': expected size to be %d (was %d)",
			info.URL,
			info.Size,
			copied,
		)
	}

	if sha1 != info.SHA1 {
		return fmt.Errorf(
			"Error downloading '%s': expected sha1 to be %s (was %s)",
			info.URL,
			info.SHA1,
			sha1,
		)
	}

	_, err = temp.Seek(0, 0)
	if err != nil {
		return
	}

	progress.renderer.StartPhase("Unpacking...")
	switch info.Extension {
	case "zip":
		err = unzip(temp, copied, path, progress.renderer)
	case "tgz":
		err = untar(temp, copied, path, progress.renderer)
	default:
		err = fmt.Errorf("Unsupported extension type: %s", info.Extension)
	}
	progress.renderer.Complete()

	return
}

func tempFile() (file *os.File, err error) {
	return ioutil.TempFile("", "android-sdk-installer")
}
