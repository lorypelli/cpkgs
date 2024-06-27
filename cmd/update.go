package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Update() {
	a := flag.Arg(1)
	headers := flag.Args()[1:]
	if a == "-a" {
		headers = []string{}
		var JSON pkg.JSON
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			return
		}
		j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
		json.Unmarshal(j, &JSON)
		for i := 0; i < len(JSON.Include.H); i++ {
			header := strings.Split(JSON.Include.H[i], "/")
			h := header[len(header)-1]
			headers = append(headers, h)
		}
		fmt.Println(headers)
	}
	if len(headers) <= 0 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Provide headers file to update: ")
		scanner.Scan()
		h := scanner.Text()
		if len(strings.TrimSpace(h)) <= 0 {
			log.Fatal("Not a valid header file!")
			return
		}
		headers = strings.Split(h, " ")
	}
	for i := 0; i < len(headers); i++ {
		if !strings.HasSuffix(headers[i], "h") {
			fmt.Printf("%s is not a valid header file\n", headers[i])
			continue
		}
		var JSON pkg.JSON
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			return
		}
		j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
		json.Unmarshal(j, &JSON)
		for i := 0; i < len(JSON.Include.H); i++ {
			h := JSON.Include.H[i]
			f := strings.Split(h, "/")
			fname := f[len(f)-1]
			if fname == headers[i] || a == "-a" {
				if a == "-a" {
					fname = headers[i]
				}
				fmt.Printf("Updating header file %s...\n", fname)
				res, err := http.Get(h)
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
				err = os.WriteFile(fmt.Sprintf("cpkgs/%s", fname), body, 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
				c := JSON.Include.C[i]
				c_files := strings.Split(c, "/")
				c_fname := c_files[len(c_files)-1]
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
				err = os.WriteFile(fmt.Sprintf("cpkgs/%s", c_fname), body, 0777)
				if err != nil {
					log.Fatal(err)
					return
				}
				fmt.Printf("Successfully updated header file %s...\n", fname)
			}
		}
	}
}