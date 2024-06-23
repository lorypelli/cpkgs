package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
		arr := strings.Split(flag.Arg(1), "/")
		file := arr[len(arr)-1]
		p := arr[:len(arr)-1]
		path := strings.Join(p, "/")
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
		json.Unmarshal(j, &JSON)
		if runtime.GOOS == "windows" {
			JSON.FileName += ".exe"
		}
		fname := JSON.FileName
		if len(path) > 0 {
			file = path + "/" + file
			fname = path + "/" + JSON.FileName
		}
		cmd := fmt.Sprintf("%s -o %s %s", JSON.Compiler, fname, file)
		for i := 0; i < len(files); i++ {
			if strings.HasSuffix(files[i].Name(), ".c") && files[i].Name() != file {
				cmd += fmt.Sprintf(" %s", files[i].Name())
			}
		}
		cmd += fmt.Sprintf(" && %s", fname)
		fmt.Println(cmd)
		break
	}
	}
}