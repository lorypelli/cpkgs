package cmd

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
	"github.com/lorypelli/cpkgs/utils"
	"github.com/pterm/pterm"
)

func Uninstall() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	json.Unmarshal(j, &JSON)
	pkgs := flag.Args()[1:]
	if len(pkgs) <= 0 {
		p, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Provide packages to uninstall").Show()
		if len(strings.TrimSpace(p)) > 0 {
			pkgs = strings.Split(p, " ")
		}
	}
	for _, pkg := range pkgs {
		if JSON.Language == "C++" {
			if !strings.HasSuffix(pkg, JSON.CPPExtensions.Header) {
				pterm.Warning.Printfln("%s is not a valid header file, skipping...", pkg)
				continue
			}
		} else {
			if !strings.HasSuffix(pkg, ".h") {
				pterm.Warning.Printfln("%s is not a valid header file, skipping...", pkg)
				continue
			}
		}
		s, _ := pterm.DefaultSpinner.Start(pterm.Sprintf("Removing package %s...\n", pkg))
		if err := os.Remove(pterm.Sprintf("cpkgs/%s", pkg)); err != nil {
			pterm.Error.Println(err)
			return
		}
		if err := os.Remove(pterm.Sprintf("cpkgs/%s", strings.ReplaceAll(pkg, ".h", ".c"))); err != nil {
			pterm.Error.Println(err)
			return
		}
		if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
			for i, h := range JSON.Include.HPP {
				header := utils.At(strings.Split(h, "/"), -1)
				if header == pkg {
					JSON.Include.HPP = append(JSON.Include.HPP[:i], JSON.Include.HPP[i+1:]...)
					JSON.Include.CPP = append(JSON.Include.CPP[:i], JSON.Include.CPP[i+1:]...)
				}
			}
		} else {
			for i, h := range JSON.Include.H {
				header := utils.At(strings.Split(h, "/"), -1)
				if header == pkg {
					JSON.Include.H = append(JSON.Include.H[:i], JSON.Include.H[i+1:]...)
					JSON.Include.C = append(JSON.Include.C[:i], JSON.Include.C[i+1:]...)
				}
			}
		}
		j, err := json.MarshalIndent(JSON, "", "  ")
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		if err = os.WriteFile("cpkgs.json", j, 0777); err != nil {
			pterm.Error.Println(err)
			return
		}
		s.Success(pterm.Sprintf("Successfully removed package %s!", pkg))
	}
}
