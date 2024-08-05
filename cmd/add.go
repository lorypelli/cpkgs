package cmd

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/v2/internal"
	"github.com/pterm/pterm"
)

func Add() {
	var JSON internal.JSON
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
		repo := strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(u.Path, u.Host, ""), "/"), "/")
		h, _ := pterm.DefaultInteractiveTextInput.WithDefaultText(pterm.Sprintf("Provide headers file to add from '%s'", repo)).Show()
		pterm.Info.Printfln("Creating cache for %s...", repo)
		cacheRepo := pterm.Sprintf("%s/%s", internal.GetCacheDir(), repo)
		if _, err := os.Stat(cacheRepo); os.IsNotExist(err) {
			if err := os.MkdirAll(cacheRepo, 0755); err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Success.Printfln("Successfully created cache for %s!", repo)
		} else {
			pterm.Warning.Printfln("Cache for %s already exists, nothing was changed!", repo)
		}
		if len(strings.TrimSpace(h)) <= 0 {
			continue
		}
		u.Host = "raw.githubusercontent.com"
		urlString := strings.ReplaceAll(u.String(), strings.TrimPrefix("github.com", "/"), "")
		headers := strings.Split(h, " ")
		for _, header := range headers {
			found := false
			include := JSON.Include.H
			if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
				include = JSON.Include.HPP
			}
			for _, h := range include {
				if header == internal.At(strings.Split(h, "/"), -1) {
					found = true
					break
				}
			}
			if found {
				pterm.Warning.Println("Header file already exists, skipping...")
				continue
			}
			if (JSON.Language == "C++" && !strings.HasSuffix(header, JSON.CPPExtensions.Header)) || !strings.HasSuffix(header, ".h") {
				pterm.Warning.Printfln("%s is not a valid header file, skipping...", header)
				continue
			}
			headerFile := internal.At(strings.Split(header, "/"), -1)
			c, err := os.ReadFile(pterm.Sprintf("%s/%s", cacheRepo, headerFile))
			if err == nil {
				pterm.Warning.Println("Header file found in cache, use 'cpkgs update' command to refresh it!")
				pterm.Info.Println("Adding from cache...")
				if err := os.MkdirAll(pterm.Sprintf("cpkgs/%s", repo), 0755); err == nil {
					if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s", repo, headerFile), c, 0644); err != nil {
						pterm.Error.Println(err)
						return
					}
					c, err = os.ReadFile(pterm.Sprintf("%s/%s_url.txt", cacheRepo, headerFile))
					if err == nil {
						if JSON.Language == "C++" {
							JSON.Include.HPP = append(JSON.Include.HPP, string(c))
						} else {
							JSON.Include.H = append(JSON.Include.H, string(c))
						}
						if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s_url.txt", repo, headerFile), c, 0644); err != nil {
							pterm.Error.Println(err)
							return
						}
					} else {
						pterm.Warning.Printfln("URL for %s not found!", header)
					}
					pterm.Success.Println("Successfully added from cache!")
				} else {
					pterm.Error.Println(err)
					return
				}
			} else {
				res, err := http.Get(pterm.Sprintf("%s/master/%s", urlString, header))
				if res.StatusCode != 200 || err != nil {
					choice, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("Before skipping this header file, do you want to try searching it in the include directory?").Show()
					if choice {
						res, err = http.Get(pterm.Sprintf("%s/master/include/%s", urlString, header))
						if res.StatusCode != 200 || err != nil {
							pterm.Error.Printf("Unable to get %s header file, skipping...\n", header)
							continue
						}
					} else {
						pterm.Error.Printf("Unable to get %s header file, skipping...\n", header)
						continue
					}
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				if _, err := os.Stat(pterm.Sprintf("cpkgs/%s", repo)); os.IsNotExist(err) {
					if err := os.MkdirAll(pterm.Sprintf("cpkgs/%s", repo), 0755); err != nil {
						pterm.Error.Println(err)
						return
					}
				}
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s", repo, headerFile), body, 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := os.WriteFile(pterm.Sprintf("%s/%s", cacheRepo, headerFile), body, 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if JSON.Language == "C++" && JSON.CPPExtensions.Header != ".h" {
					JSON.Include.HPP = append(JSON.Include.HPP, res.Request.URL.String())
				} else {
					JSON.Include.H = append(JSON.Include.H, res.Request.URL.String())
				}
				if err := os.WriteFile(pterm.Sprintf("%s/%s_url.txt", cacheRepo, headerFile), []byte(res.Request.URL.String()), 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s_url.txt", repo, headerFile), []byte(res.Request.URL.String()), 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
			}
			var code string
			if JSON.Language == "C++" {
				code = strings.ReplaceAll(header, JSON.CPPExtensions.Header, JSON.CPPExtensions.Code)
			} else {
				code = strings.ReplaceAll(header, ".h", ".c")
			}
			codeFile := internal.At(strings.Split(code, "/"), -1)
			c, err = os.ReadFile(pterm.Sprintf("%s/%s", cacheRepo, codeFile))
			if err == nil {
				pterm.Warning.Println("Code file found in cache, use 'cpkgs update' command to refresh it!")
				pterm.Info.Println("Adding from cache...")
				if err := os.MkdirAll(pterm.Sprintf("cpkgs/%s", repo), 0755); err == nil {
					if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s", repo, codeFile), c, 0644); err != nil {
						pterm.Error.Println(err)
						return
					}
					c, err = os.ReadFile(pterm.Sprintf("%s/%s_url.txt", cacheRepo, codeFile))
					if err == nil {
						if JSON.Language == "C++" {
							JSON.Include.CPP = append(JSON.Include.CPP, string(c))
						} else {
							JSON.Include.C = append(JSON.Include.C, string(c))
						}
						if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s_url.txt", repo, codeFile), c, 0644); err != nil {
							pterm.Error.Println(err)
							return
						}
					} else {
						pterm.Warning.Printfln("URL for %s not found!", code)
					}
					pterm.Success.Println("Successfully added from cache!")
				} else {
					pterm.Error.Println(err)
					return
				}
			} else {
				res, err := http.Get(pterm.Sprintf("%s/master/%s", urlString, code))
				for res.StatusCode != 200 || err != nil {
					dir, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Code file not found, please provide directory").Show()
					res, err = http.Get(pterm.Sprintf("%s/master/%s/%s", urlString, dir, code))
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s", repo, codeFile), body, 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := os.WriteFile(pterm.Sprintf("%s/%s", cacheRepo, codeFile), body, 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if JSON.Language == "C++" {
					JSON.Include.CPP = append(JSON.Include.CPP, res.Request.URL.String())
				} else {
					JSON.Include.C = append(JSON.Include.C, res.Request.URL.String())
				}
				if err := os.WriteFile(pterm.Sprintf("%s/%s_url.txt", cacheRepo, codeFile), []byte(res.Request.URL.String()), 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
				if err := os.WriteFile(pterm.Sprintf("cpkgs/%s/%s_url.txt", repo, codeFile), []byte(res.Request.URL.String()), 0644); err != nil {
					pterm.Error.Println(err)
					return
				}
			}
			j, err := json.MarshalIndent(JSON, "", "  ")
			if err != nil {
				pterm.Error.Println(err)
				return
			}
			if err = os.WriteFile("cpkgs.json", j, 0644); err != nil {
				pterm.Error.Println(err)
				return
			}
			pterm.Success.Printfln("Successfully added header file %s!", header)
		}
	}
}
