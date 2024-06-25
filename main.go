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
	Compiler string  `json:"compiler"`
	FileName string  `json:"file_name"`
	Include  Include `json:"include"`
}

type Include struct {
	C []string
	H []string
}

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, e := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	if os.IsNotExist(e) {
		var JSON JSON
		JSON.Include.C = []string{}
		JSON.Include.H = []string{}
		j, _ := json.Marshal(JSON)
		os.WriteFile("cpkgs.json", j, 0777)
	}
	var JSON JSON
	j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	err = json.Unmarshal(j, &JSON)
	if err != nil {
		JSON.Include.C = []string{}
		JSON.Include.H = []string{}
		j, _ := json.Marshal(JSON)
		os.WriteFile("cpkgs.json", j, 0777)
	}
	switch flag.Arg(0) {
	case "run":
		{
			_, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
				return
			}
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
			cmd := fmt.Sprintf("cd %s && %s -o %s %s", path, JSON.Compiler, fname, strings.Join(flag.Args()[1:], " "))
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
			cmd += fmt.Sprintf(" && %s", JSON.FileName)
			cmdExec := exec.Command("sh", "-c", cmd)
			if runtime.GOOS == "windows" {
				cmdExec = exec.Command("cmd", "/C", cmd)
			}
			cmdExec.Stdin = os.Stdin
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Run()
			break
		}
	case "add":
		{
			pkgs := flag.Args()[1:]
			scanner := bufio.NewScanner(os.Stdin)
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
				fmt.Printf("Specify headers file to add from '%s': ", strings.TrimPrefix(strings.ReplaceAll(u.Path, u.Host, ""), "/"))
				u.Host = "raw.githubusercontent.com"
				urlString := strings.ReplaceAll(u.String(), "/github.com", "")
				scanner.Scan()
				h := scanner.Text()
				if len(strings.TrimSpace(h)) <= 0 {
					continue
				}
				headers := strings.Split(h, " ")
				c := 0
				for i := 0; i < len(headers); i++ {
					if !strings.HasSuffix(headers[i], "h") {
						fmt.Printf("%s is not a valid header file\n", headers[i])
						continue
					}
					res, err := http.Get(fmt.Sprintf("%s/master/%s", urlString, headers[i+c]))
					for res.StatusCode != 200 || err != nil {
						var choice string
						fmt.Print("Before skipping this header file, do you want to try searching it in the include directory? (Y/n) ")
						fmt.Scanln(&choice)
						if len(choice) <= 0 {
							choice = "y"
						}
						if strings.ToLower(choice) == "y" {
							res, err = http.Get(fmt.Sprintf("%s/master/include/%s", urlString, headers[i+c]))
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
						res, err = http.Get(fmt.Sprintf("%s/master/%s", urlString, headers[i+c]))
					}
					defer res.Body.Close()
					body, err := io.ReadAll(res.Body)
					if err != nil {
						log.Fatal(err)
						return
					}
					_, err = os.Stat("cpkgs")
					if os.IsNotExist(err) {
						err = os.Mkdir("cpkgs", 0777)
						if err != nil {
							log.Fatal(err)
							return
						}
					}
					os.WriteFile(fmt.Sprintf("cpkgs/%s", headers[i]), body, 0777)
					JSON.Include.H = append(JSON.Include.H, res.Request.URL.String())
					code := strings.ReplaceAll(headers[i], ".h", ".c")
					res, err = http.Get(fmt.Sprintf("%s/master/%s", urlString, strings.ReplaceAll(headers[i], ".h", ".c")))
					for res.StatusCode != 200 || err != nil {
						var dir string
						fmt.Print("C code file not found, please provide directory: ")
						fmt.Scan(&dir)
						res, err = http.Get(fmt.Sprintf("%s/master/%s/%s", u.String(), dir, headers[i]))
					}
					defer res.Body.Close()
					body, err = io.ReadAll(res.Body)
					if err != nil {
						log.Fatal(err)
						return
					}
					os.WriteFile(fmt.Sprintf("cpkgs/%s", code), body, 0777)
					JSON.Include.C = append(JSON.Include.C, res.Request.URL.String())
					j, err := json.Marshal(JSON)
					if err != nil {
						log.Fatal(err)
						return
					}
					os.WriteFile("cpkgs.json", j, 0777)
				}
			}
			break
		}
	case "init":
		{
			os.WriteFile("cpkgs.json", j, 0777)
			var compiler, filename string
			fmt.Print("Specify the compiler to use: ")
			fmt.Scanln(&compiler)
			if len(strings.TrimSpace(compiler)) <= 0 {
				fmt.Println("Using default...(gcc)")
				compiler = "gcc"
			}
			JSON.Compiler = compiler
			fmt.Print("Specify the output filename: ")
			fmt.Scanln(&filename)
			if len(strings.TrimSpace(filename)) <= 0 {
				fmt.Println("Using default...(out)")
				filename = "out"
			}
			JSON.FileName = filename
			j, err := json.Marshal(JSON)
			if err != nil {
				log.Fatal(err)
				return
			}
			os.WriteFile("cpkgs.json", j, 0777)
			break
		}
	case "install":
		{
			_, err = os.Stat("cpkgs")
			if os.IsNotExist(err) {
				err = os.Mkdir("cpkgs", 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
			for i := 0; i < len(JSON.Include.H); i++ {
				res, err := http.Get(JSON.Include.H[i])
				if err != nil {
					log.Fatal(err)
					return
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
					return
				}
				filename := strings.Split(JSON.Include.H[i], "/")
				os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
				res, err = http.Get(JSON.Include.C[i])
				if err != nil {
					log.Fatal(err)
					return
				}
				defer res.Body.Close()
				body, err = io.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
					return
				}
				filename = strings.Split(JSON.Include.C[i], "/")
				os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
			}
			break
		}
	case "help":
		{
			fmt.Println("List of all commands:")
			fmt.Print("\n")
			fmt.Println("---------------------------------------------------------------------------")
			fmt.Println("|'cpkgs init' - initialize a new project using cpkgs                      |")
			fmt.Println("|'cpkgs add <package-name>' - add C packages using cpkgs                  |")
			fmt.Println("|'cpkgs install' - install all the packages in the current project        |")
			fmt.Println("|'cpkgs run <file-name>' - run the file name using your selected compiler |")
			fmt.Println("|'cpkgs help' - show this menu                                            |")
			fmt.Println("---------------------------------------------------------------------------")
			fmt.Print("\n")
		}
	default:
		{
			log.Fatal("Unknown command, to see all avaible commands type: 'cpkgs help' ")
		}
	}
}
