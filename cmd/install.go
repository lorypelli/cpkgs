package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Install() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(j, &JSON)
	if len(flag.Args()) > 0 {
		if pkgs := flag.Args()[1:]; len(pkgs) > 0 {
			fmt.Println("You provided arguments to the command, 'cpkgs add' will be executed instead!")
			cmd := fmt.Sprintf("cpkgs add %s", strings.Join(pkgs, " "))
			cmdExec := exec.Command("sh", "-c", cmd)
			if runtime.GOOS == "windows" {
				cmdExec = exec.Command("cmd", "/C", cmd)
			}
			cmdExec.Stdin = os.Stdin
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Run()
			return
		}
	}
	fmt.Println("Resolving packages...")
	if _, err := os.Stat("cpkgs"); os.IsNotExist(err) {
		if err := os.Mkdir("cpkgs", 0777); err != nil {
			log.Fatal(err)
			return
		}
	}
	if len(JSON.Include.H) <= 0 {
		fmt.Println("No packages found!")
	}
	for _, h := range JSON.Include.H {
		res, err := http.Get(h)
		pkg := strings.Split(h, "/")
		fmt.Printf("Installing package %s...\n", pkg[len(pkg)-1])
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
		filename := strings.Split(h, "/")
		if err := os.WriteFile(fmt.Sprintf("cpkgs/%s", filename[len(filename)-1]), body, 0777); err != nil {
			log.Fatal(err)
			return
		}
		c := strings.ReplaceAll(h, ".h", ".c")
		res, err = http.Get(c)
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
		filename = strings.Split(c, "/")
		if err := os.WriteFile(fmt.Sprintf("cpkgs/%s", filename[len(filename)-1]), body, 0777); err != nil {
			log.Fatal(err)
			return
		}
	}
}
