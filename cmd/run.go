package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Run() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(j, &JSON)
	f := flag.Arg(1)
	if len(strings.TrimSpace(f)) <= 0 {
		f = "main.c"
	}
	file, err := filepath.Abs(f)
	if err != nil {
		log.Fatal(err)
		return
	}
	path := strings.ReplaceAll(filepath.Dir(file), "\\", "/")
	if runtime.GOOS == "windows" {
		JSON.FileName += ".exe"
	}
	fname := JSON.FileName
	_, err = os.Stat("cpkgs")
	if os.IsNotExist(err) {
		err = os.Mkdir("cpkgs", 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	_, err = os.Stat("cpkgs/bin")
	if os.IsNotExist(err) {
		err = os.Mkdir("cpkgs/bin", 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	cmd := fmt.Sprintf("cd %s && %s -o cpkgs/bin/%s %s %s", path, JSON.Compiler, fname, file, strings.Join(flag.Args()[1:], " "))
	files, err := os.ReadDir(fmt.Sprintf("%s/cpkgs", path))
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".c") {
			cmd += fmt.Sprintf(" cpkgs/%s", f.Name())
		}
	}
	cmd += fmt.Sprintf(" && cd cpkgs/bin && %s", JSON.FileName)
	cmdExec := exec.Command("sh", "-c", cmd)
	if runtime.GOOS == "windows" {
		cmdExec = exec.Command("cmd", "/C", cmd)
	}
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Run()
}
