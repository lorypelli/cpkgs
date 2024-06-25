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
	var JSON JSON
	j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	json.Unmarshal(j, &JSON)
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
			if runtime.GOOS == "windows" {
				JSON.FileName += ".exe"
			}
			fname := JSON.FileName
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
			break
		}
	case "add":
		{
			pkgs := flag.Args()[1:]
			if len(pkgs) <= 0 {
				fmt.Print("Provide packages to add: ")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				p := scanner.Text()
				pkgs = strings.Split(p, " ")
			}
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
				fmt.Printf("Provide headers file to add from '%s': ", strings.TrimPrefix(strings.ReplaceAll(u.Path, u.Host, ""), "/"))
				u.Host = "raw.githubusercontent.com"
				urlString := strings.ReplaceAll(u.String(), "/github.com", "")
				scanner.Scan()
				h := scanner.Text()
				if len(strings.TrimSpace(h)) <= 0 {
					continue
				}
				found := false
				for i := 0; i < len(JSON.Include.H) && !found; i++ {
					url := strings.Split(JSON.Include.H[i], "/")
					if h == url[len(url)-1] {
						found = true
					}
				}
				if found {
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
					err = os.WriteFile(fmt.Sprintf("cpkgs/%s", headers[i]), body, 0777)
					if err != nil {
						log.Fatal(err)
						return
					}
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
					err = os.WriteFile(fmt.Sprintf("cpkgs/%s", code), body, 0777)
					if err != nil {
						log.Fatal(err)
						return
					}
					JSON.Include.C = append(JSON.Include.C, res.Request.URL.String())
					j, err := json.Marshal(JSON)
					if err != nil {
						log.Fatal(err)
						return
					}
					err = os.WriteFile("cpkgs.json", j, 0777)
					if err != nil {
						log.Fatal(err)
						return
					}
				}
			}
			break
		}
	case "init":
		{
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
				return
			}
			d := flag.Arg(1)
			if d != "-d" && len(strings.TrimSpace(flag.Arg(1))) > 0 {
				dir = flag.Arg(1)
				d = flag.Arg(2)
			}
			_, e := os.Stat(dir)
			if os.IsNotExist(e) {
				os.Mkdir(dir, 0777)
			}
			_, e = os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
			if os.IsNotExist(e) {
				JSON.Include.C = []string{}
				JSON.Include.H = []string{}
				j, _ := json.Marshal(JSON)
				err = os.WriteFile("cpkgs.json", j, 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
			err = os.WriteFile(fmt.Sprintf("%s/cpkgs.json", dir), j, 0777)
			if err != nil {
				log.Fatal(err)
				return
			}
			var compiler, filename string
			if d != "-d" {
				fmt.Print("Provide the compiler to use: ")
				fmt.Scanln(&compiler)
				if len(strings.TrimSpace(compiler)) <= 0 {
					fmt.Println("Using default...(gcc)")
					compiler = "gcc"
				}
				fmt.Print("Provide the output filename: ")
				fmt.Scanln(&filename)
				if len(strings.TrimSpace(filename)) <= 0 {
					fmt.Println("Using default...(out)")
					filename = "out"
				}
			} else {
				fmt.Println("Using defaults...(gcc, out)")
				compiler = "gcc"
				filename = "out"
			}
			JSON.Compiler = compiler
			JSON.FileName = filename
			j, err := json.Marshal(JSON)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = os.WriteFile(fmt.Sprintf("%s/cpkgs.json", dir), j, 0777)
			if err != nil {
				log.Fatal(err)
				return
			}
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
				err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
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
				err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
			break
		}
	case "help":
		{
			fmt.Println("List of all commands:")
			fmt.Print("\n")
			fmt.Println("---------------------------------------------------------------------------")
			fmt.Println("|'cpkgs init <dir-name> [-d]' - initialize a new project using cpkgs      |")
			fmt.Println("|'cpkgs add <package-name>' - add C packages using cpkgs                  |")
			fmt.Println("|'cpkgs install' - install all the packages in the current project        |")
			fmt.Println("|'cpkgs run <file-name>' - run the file name using your selected compiler |")
			fmt.Println("|'cpkgs help' - show this menu                                            |")
			fmt.Println("---------------------------------------------------------------------------")
			fmt.Print("\n")
			break
		}
	default:
		{
			log.Fatal("Unknown command, to see all avaible commands type: 'cpkgs help' ")
		}
	}
}
