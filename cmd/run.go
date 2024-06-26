package cmd

import (
	"cpkgs/pkg"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Run() {
	var JSON pkg.JSON
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	json.Unmarshal(j, &JSON)
	f := filepath.Clean(flag.Arg(1))
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
	cmd := fmt.Sprintf("cd %s && %s -o cpkgs/bin/%s %s", path, JSON.Compiler, fname, strings.Join(flag.Args()[1:], " "))
	files, err := os.ReadDir(fmt.Sprintf("%s/cpkgs", path))
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := 0; i < len(files); i++ {
		if strings.HasSuffix(files[i].Name(), ".c") {
			cmd += fmt.Sprintf(" cpkgs/%s", files[i].Name())
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
