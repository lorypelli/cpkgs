package cmd

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
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Add() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(j, &JSON)
	pkgs := flag.Args()[1:]
	if len(pkgs) <= 0 {
		fmt.Print("Provide packages to add: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		p := scanner.Text()
		pkgs = strings.Split(p, " ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	for _, pkg := range pkgs {
		u, _ := url.Parse(pkg)
		if len(u.Scheme) <= 0 || len(u.Host) <= 0 {
			u.Scheme = "https"
			u.Host = "github.com"
		}
		if u.Host != "github.com" {
			log.Fatal("Currently only github is supported!")
			return
		}
		fmt.Printf("Provide headers file to add from '%s': ", strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(u.Path, u.Host, ""), "/"), "/"))
		u.Host = "raw.githubusercontent.com"
		urlString := strings.ReplaceAll(u.String(), "/github.com", "")
		scanner.Scan()
		h := scanner.Text()
		if len(strings.TrimSpace(h)) <= 0 {
			continue
		}
		headers := strings.Split(h, " ")
		for _, header := range headers {
			found := false
			for _, h := range JSON.Include.H {
				url := strings.Split(h, "/")
				if header == url[len(url)-1] {
					found = true
					break
				}
			}
			if found {
				fmt.Println("Header file already exists, skipping...")
				continue
			}
			if !strings.HasSuffix(header, ".h") {
				fmt.Printf("%s is not a valid header file, skipping...\n", header)
				continue
			}
			res, err := http.Get(fmt.Sprintf("%s/master/%s", urlString, header))
			for res.StatusCode != 200 || err != nil {
				var choice string
				fmt.Print("Before skipping this header file, do you want to try searching it in the include directory? (Y/n) ")
				fmt.Scanln(&choice)
				if len(choice) <= 0 {
					choice = "y"
				}
				if strings.ToLower(choice) == "y" {
					res, err = http.Get(fmt.Sprintf("%s/master/include/%s", urlString, header))
					if res.StatusCode != 200 || err != nil {
						fmt.Printf("Unable to get %s header file, skipping...\n", header)
						continue
					}
				} else {
					fmt.Printf("Unable to get %s header file, skipping...\n", header)
					continue
				}
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
				return
			}
			if _, err := os.Stat("cpkgs"); os.IsNotExist(err) {
				if err := os.Mkdir("cpkgs", 0777); err != nil {
					log.Fatal(err)
					return
				}
			}
			if err := os.WriteFile(fmt.Sprintf("cpkgs/%s", header), body, 0777); err != nil {
				log.Fatal(err)
				return
			}
			JSON.Include.H = append(JSON.Include.H, res.Request.URL.String())
			code := strings.ReplaceAll(header, ".h", ".c")
			res, err = http.Get(fmt.Sprintf("%s/master/%s", urlString, code))
			for res.StatusCode != 200 || err != nil {
				var dir string
				fmt.Print("C code file not found, please provide directory: ")
				fmt.Scan(&dir)
				res, err = http.Get(fmt.Sprintf("%s/master/%s/%s", u.String(), dir, header))
			}
			defer res.Body.Close()
			body, err = io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
				return
			}
			if err := os.WriteFile(fmt.Sprintf("cpkgs/%s", code), body, 0777); err != nil {
				log.Fatal(err)
				return
			}
			JSON.Include.C = append(JSON.Include.C, res.Request.URL.String())
			j, err := json.MarshalIndent(JSON, "", "  ")
			if err != nil {
				log.Fatal(err)
				return
			}
			if err = os.WriteFile("cpkgs.json", j, 0777); err != nil {
				log.Fatal(err)
				return
			}
			fmt.Printf("Successfully added header file %s!\n", header)
		}
	}
}
