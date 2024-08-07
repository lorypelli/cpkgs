package cmd

import (
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lorypelli/cpkgs/v2/internal"
	"github.com/pterm/pterm"
)

func Run() {
	var JSON internal.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	json.Unmarshal(j, &JSON)
	f := flag.Arg(1)
	c := 2
	if len(strings.TrimSpace(f)) <= 0 {
		f = "main.c"
		c = 1
	}
	file, err := filepath.Abs(f)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	path := strings.ReplaceAll(filepath.Dir(file), "\\", "/")
	if runtime.GOOS == "windows" {
		JSON.FileName += ".exe"
	}
	fname := JSON.FileName
	if _, err := os.Stat("cpkgs"); os.IsNotExist(err) {
		if err := os.Mkdir("cpkgs", 0755); err != nil {
			pterm.Error.Println(err)
			return
		}
	}
	if _, err := os.Stat("cpkgs/bin"); os.IsNotExist(err) {
		if err := os.Mkdir("cpkgs/bin", 0755); err != nil {
			pterm.Error.Println(err)
			return
		}
	}
	cmd := pterm.Sprintf("cd %s && %s -o cpkgs/bin/%s %s %s", path, JSON.Compiler, fname, f, strings.Join(flag.Args()[c:], " "))
	if err := filepath.WalkDir(pterm.Sprintf("%s/cpkgs", path), func(p string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".c") {
			cmd += pterm.Sprintf(" %s", strings.TrimPrefix(p, pterm.Sprintf("%s/", path)))
		}
		return nil
	}); err != nil {
		pterm.Error.Println(err)
		return
	}
	cmd += pterm.Sprintf(" && ./cpkgs/bin/%s", JSON.FileName)
	cmdExec := exec.Command("bash", "-c", cmd)
	if runtime.GOOS == "windows" {
		cmdExec = exec.Command("cmd", "/C", cmd)
	}
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Run()
}
