package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
	"github.com/lorypelli/cpkgs/utils"
	"github.com/pterm/pterm"
)

func Update() {
	a := flag.Arg(1)
	headers := flag.Args()[1:]
	if a == "-a" || a == "--all" {
		headers = []string{}
		var JSON pkg.JSON
		j, err := os.ReadFile("cpkgs.json")
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		json.Unmarshal(j, &JSON)
		for _, h := range JSON.Include.H {
			h := utils.At(strings.Split(h, "/"), -1)
			headers = append(headers, h)
		}
	}
	if len(headers) <= 0 {
		h, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Provide headers file to update").Show()
		if len(strings.TrimSpace(h)) > 0 {
			headers = strings.Split(h, " ")
		}
	}
	for _, header := range headers {
		if !strings.HasSuffix(header, ".h") {
			pterm.Warning.Printfln("%s is not a valid header file, skipping...", header)
			continue
		}
		var JSON pkg.JSON
		j, err := os.ReadFile("cpkgs.json")
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		json.Unmarshal(j, &JSON)
		include := JSON.Include.H
		if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
			include = JSON.Include.HPP
		}
		for _, h := range include {
			f := strings.Split(h, "/")
			fname := f[len(f)-1]
			if fname == header || a == "-a" || a == "--all" {
				if a == "-a" || a == "--all" {
					fname = header
				}
				s, _ := pterm.DefaultSpinner.Start(pterm.Sprintf("Updating header file %s...", fname))
				res, err := http.Get(h)
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
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", fname), body, 0777); err != nil {
					pterm.Error.Println(err)
					return
				}
				var code string
				if JSON.Language == "C++" {
					code = strings.ReplaceAll(header, JSON.CPPExtensions.Header, JSON.CPPExtensions.Code)
				} else {
					code = strings.ReplaceAll(header, ".h", ".c")
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
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", utils.At(strings.Split(code, "/"), -1)), body, 0777); err != nil {
					pterm.Error.Println(err)
					return
				}
				s.Success(pterm.Sprintf("Successfully updated header file %s...\n", fname))
			}
		}
	}
}
