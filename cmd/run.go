package cmd

import (
	"encoding/json"
	"flag"
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
	files, err := os.ReadDir(pterm.Sprintf("%s/cpkgs", path))
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".c") {
			cmd += pterm.Sprintf(" cpkgs/%s", f.Name())
		}
	}
	cmd += pterm.Sprintf(" && cd cpkgs/bin && %s", JSON.FileName)
	cmdExec := exec.Command("sh", "-c", cmd)
	if runtime.GOOS == "windows" {
		cmdExec = exec.Command("cmd", "/C", cmd)
	}
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Run()
}
