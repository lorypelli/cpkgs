package internal

import (
	"os"
	"runtime"

	"github.com/pterm/pterm"
)

func GetCacheDir() string {
	if runtime.GOOS == "windows" {
		return pterm.Sprintf("%s/cpkgs/cache", os.Getenv("APPDATA"))
	}
	return "~/.cache/cpkgs"
}
