package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
	"github.com/pterm/pterm"
)

func Install() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	json.Unmarshal(j, &JSON)
	if len(flag.Args()) > 0 {
		if pkgs := flag.Args()[1:]; len(pkgs) > 0 {
			pterm.Warning.Println("You provided arguments to the command, 'cpkgs add' will be executed instead!")
			cmd := pterm.Sprintf("cpkgs add %s", strings.Join(pkgs, " "))
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
	pterm.Info.Println("Resolving packages...")
	if _, err := os.Stat("cpkgs"); os.IsNotExist(err) {
		if err := os.Mkdir("cpkgs", 0777); err != nil {
			pterm.Error.Println(err)
			return
		}
	}
	if len(JSON.Include.H) <= 0 {
		pterm.Error.Println("No packages found!")
		return
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(len(JSON.Include.H)).WithTitle("Resolving packages...").Start()
	for _, h := range JSON.Include.H {
		res, err := http.Get(h)
		pkg := strings.Split(h, "/")
		p.UpdateTitle(pterm.Sprintf("Installing package %s...", pkg[len(pkg)-1]))
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		filename := strings.Split(h, "/")
		if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", filename[len(filename)-1]), body, 0777); err != nil {
			pterm.Error.Println(err)
			return
		}
		c := strings.ReplaceAll(h, ".h", ".c")
		res, err = http.Get(c)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		filename = strings.Split(c, "/")
		if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", filename[len(filename)-1]), body, 0777); err != nil {
			pterm.Error.Println(err)
			return
		}
		p.Increment()
	}
	pterm.Success.Println("Successfully installed all packages!")
}
