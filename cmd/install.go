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

	"github.com/lorypelli/cpkgs/v2/internal"
	"github.com/pterm/pterm"
)

func Install() {
	var JSON internal.JSON
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
			cmdExec := exec.Command("bash", "-c", cmd)
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
		if err := os.Mkdir("cpkgs", 0755); err != nil {
			pterm.Error.Println(err)
			return
		}
	}
	include := JSON.Include.H
	if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
		include = JSON.Include.HPP
	}
	if len(include) <= 0 {
		pterm.Error.Println("No packages found!")
		return
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(len(include)).WithTitle("Resolving packages...").Start()
	for _, h := range include {
		res, err := http.Get(h)
		pkg := internal.At(strings.Split(h, "/"), -1)
		p.UpdateTitle(pterm.Sprintf("Installing package %s...", pkg))
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
		if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", pkg), body, 0644); err != nil {
			pterm.Error.Println(err)
			return
		}
		var code string
		if JSON.Language == "C++" {
			code = strings.ReplaceAll(h, JSON.CPPExtensions.Header, JSON.CPPExtensions.Code)
		} else {
			code = strings.ReplaceAll(h, ".h", ".c")
		}
		res, err = http.Get(code)
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
		if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", internal.At(strings.Split(code, "/"), -1)), body, 0644); err != nil {
			pterm.Error.Println(err)
			return
		}
		p.Increment()
	}
	pterm.Success.Println("Successfully installed all packages!")
}
