package main

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
)

type JSON struct {
	Compiler string
	FileName string
	Include Include
}

type Include struct {
	C []string
	H []string
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "run": {
		f := filepath.Clean(flag.Arg(1))
		file, err := filepath.Abs(f)
		if err != nil {
			log.Fatal(err)
			return
		}
		path := strings.ReplaceAll(filepath.Dir(file), "\\", "/")
		files, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
			return
		}
		j, err := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", path))
		if err != nil {
			log.Fatal(err)
			return
		}
		var JSON JSON
		err = json.Unmarshal(j, &JSON)
		if err != nil {
			log.Fatal(err)
			return
		}
		if runtime.GOOS == "windows" {
			JSON.FileName += ".exe"
		}
		fname := JSON.FileName
		cmd := fmt.Sprintf("cd %s && %s -o %s", path, JSON.Compiler, fname)
		for i := 0; i < len(files); i++ {
			if strings.HasSuffix(files[i].Name(), ".c") {
				cmd += fmt.Sprintf(" %s", files[i].Name())
			}
		}
		cmd += fmt.Sprintf(" && %s", JSON.FileName)
		cmdExec, err := exec.Command("sh", "-c", cmd).CombinedOutput()
		if runtime.GOOS == "windows" {
			cmdExec, err = exec.Command("cmd", "/C", cmd).CombinedOutput()
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(string(cmdExec))
		break
	}
	}
}