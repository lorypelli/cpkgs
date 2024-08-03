package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/v2/internal"
	"github.com/pterm/pterm"
)

func Update() {
	var JSON internal.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	json.Unmarshal(j, &JSON)
	a := strings.ToLower(flag.Arg(1))
	headers := flag.Args()[1:]
	include := JSON.Include.H
	if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
		include = JSON.Include.HPP
	}
	if a == "-a" || a == "--all" {
		headers = []string{}
		for _, h := range include {
			h := internal.At(strings.Split(h, "/"), -1)
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
		if (JSON.Language == "C++" && !strings.HasSuffix(header, JSON.CPPExtensions.Header)) || !strings.HasSuffix(header, ".h") {
			pterm.Warning.Printfln("%s is not a valid header file, skipping...", header)
			continue
		}
		for _, h := range include {
			f := internal.At(strings.Split(h, "/"), -1)
			if f == header || a == "-a" || a == "--all" {
				if a == "-a" || a == "--all" {
					f = header
				}
				s, _ := pterm.DefaultSpinner.Start(pterm.Sprintf("Updating header file %s...", f))
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
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", f), body, 0644); err != nil {
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
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", internal.At(strings.Split(code, "/"), -1)), body, 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				s.Success(pterm.Sprintf("Successfully updated header file %s...\n", f))
			}
		}
	}
}
