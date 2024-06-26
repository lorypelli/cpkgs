package cmd

import (
	"bufio"
	"cpkgs/pkg"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Add() {
	var JSON pkg.JSON
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
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
			fmt.Println("Header file already exists, skipping...")
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
			fmt.Printf("Successfully added header file %s!\n", headers[i])
		}
	}
}
