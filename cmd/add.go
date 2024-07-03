package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
	"github.com/pterm/pterm"
)

func Add() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	json.Unmarshal(j, &JSON)
	pkgs := flag.Args()[1:]
	if len(pkgs) <= 0 {
		p, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Provide packages to add").Show()
		pkgs = strings.Split(p, " ")
	}
	for _, pkg := range pkgs {
		u, _ := url.Parse(pkg)
		if len(u.Scheme) <= 0 || len(u.Host) <= 0 {
			u.Scheme = "https"
			u.Host = "github.com"
		}
		if u.Host != "github.com" {
			pterm.Error.Println("Currently only github is supported!")
			return
		}
		h, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(pterm.Sprintf("Provide headers file to add from '%s'", strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(u.Path, u.Host, ""), "/"), "/"))).Show()
		u.Host = "raw.githubusercontent.com"
		urlString := strings.ReplaceAll(u.String(), "/github.com", "")
		if len(strings.TrimSpace(h)) <= 0 {
			continue
		}
		headers := strings.Split(h, " ")
		for _, header := range headers {
			found := false
			if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
				for _, h := range JSON.Include.HPP {
					url := strings.Split(h, "/")
					if header == url[len(url)-1] {
						found = true
						break
					}
				}
			} else {
				for _, h := range JSON.Include.H {
					url := strings.Split(h, "/")
					if header == url[len(url)-1] {
						found = true
						break
					}
				}
			}
			if found {
				pterm.Warning.Println("Header file already exists, skipping...")
				continue
			}
			if JSON.Language == "C++" {
				if !strings.HasSuffix(header, JSON.CPPExtensions.Header) {
					pterm.Warning.Printfln("%s is not a valid header file, skipping...", header)
					continue
				}
			} else {
				if !strings.HasSuffix(header, ".h") {
					pterm.Warning.Printfln("%s is not a valid header file, skipping...", header)
					continue
				}
			}
			res, err := http.Get(pterm.Sprintf("%s/master/%s", urlString, header))
			for res.StatusCode != 200 || err != nil {
				choice, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("Before skipping this header file, do you want to try searching it in the include directory?").Show()
				if choice {
					res, err = http.Get(pterm.Sprintf("%s/master/include/%s", urlString, header))
					if res.StatusCode != 200 || err != nil {
						pterm.Error.Printf("Unable to get %s header file, skipping...\n", header)
						break
					}
				} else {
					pterm.Error.Printf("Unable to get %s header file, skipping...\n", header)
					break
				}
			}
			if res.StatusCode != 200 {
				continue
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			if _, err := os.Stat("cpkgs"); os.IsNotExist(err) {
				if err := os.Mkdir("cpkgs", 0777); err != nil {
					pterm.Error.Println(err)
					return
				}
			}
			if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", header), body, 0777); err != nil {
				pterm.Error.Println(err)
				return
			}
			if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
				JSON.Include.HPP = append(JSON.Include.HPP, res.Request.URL.String())
			} else {
				JSON.Include.H = append(JSON.Include.H, res.Request.URL.String())
			}
			var code string
			if JSON.Language == "C++" {
				code = strings.ReplaceAll(header, JSON.CPPExtensions.Header, JSON.CPPExtensions.Code)
			} else {
				code = strings.ReplaceAll(header, ".h", ".c")
			}
			res, err = http.Get(pterm.Sprintf("%s/master/%s", urlString, code))
			for res.StatusCode != 200 || err != nil {
				dir, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Code file not found, please provide directory").Show()
				res, err = http.Get(pterm.Sprintf("%s/master/%s/%s", u.String(), dir, header))
			}
			defer res.Body.Close()
			body, err = io.ReadAll(res.Body)
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			if err := os.WriteFile(pterm.Sprintf("cpkgs/%s", code), body, 0777); err != nil {
				pterm.Error.Println(err)
				return
			}
			if JSON.Language == "C++" {
				JSON.Include.CPP = append(JSON.Include.CPP, res.Request.URL.String())
			} else {
				JSON.Include.C = append(JSON.Include.C, res.Request.URL.String())
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
			pterm.Success.Printfln("Successfully added header file %s!", header)
		}
	}
}
