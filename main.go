package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type JSON struct {
	Compiler string
	FileName string
	Include  Include
}

type Include struct {
	C []string
	H []string
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "run":
		{
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
	case "add":
		{
			pkgs := flag.Args()[1:]
			for i := 0; i < len(pkgs); i++ {
				u, _ := url.Parse(pkgs[i])
				if len(u.Scheme) <= 0 || len(u.Host) <= 0 {
					u.Scheme = "https"
					u.Host = "github.com"
				}
				if u.Host != "github.com" {
					log.Fatal("Currently only github is supported!")
					return
				}
				u.Host = "raw.githubusercontent.com"
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Printf("Specify headers file to add from %s ", u.Path)
				scanner.Scan()
				h := scanner.Text()
				headers := strings.Split(h, " ")
				c := 0
				for i := 0; i < len(headers); i++ {
					if !strings.HasSuffix(headers[i], "h") {
						fmt.Printf("%s is not a valid header file\n", headers[i])
						continue
					}
					res, err := http.Get(fmt.Sprintf("%s/main/%s", u.String(), headers[i+c]))
					for res.StatusCode != 200 || err != nil {
						var choice string
						fmt.Println("Before skipping this header file, do you want to try searching it in the include directory? (Y/n)")
						fmt.Scan(&choice)
						if strings.ToLower(choice) == "y" {
							fmt.Println("Searching...")
							res, err = http.Get(fmt.Sprintf("%s/main/%s", u.String(), headers[i+c]))
							if res.StatusCode != 200 || err != nil {
								fmt.Printf("Unable to get %s header file, skipping...\n", headers[i+c])
							}
						} else {
							fmt.Printf("Unable to get %s header file, skipping...\n", headers[i+c])
						}
						c++
						if i+c >= len(headers) {
							break
						}
						fmt.Println(headers[i+c])
						res, err = http.Get(fmt.Sprintf("%s/main/%s", u.String(), headers[i+c]))
					}
					defer res.Body.Close()
					body, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatal(err)
						return
					}
					fmt.Println(string(body))
				}
			}
			break
		}
	}
}
