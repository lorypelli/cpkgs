package cmd

import (
	"cpkgs/pkg"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Install() {
	fmt.Println("Resolving packages...")
	var JSON pkg.JSON
	dir, err := os.Getwd()
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
	if len(JSON.Include.H) <= 0 {
		fmt.Println("No packages found!")
	}
	for i := 0; i < len(JSON.Include.H); i++ {
		res, err := http.Get(JSON.Include.H[i])
		pkg := strings.Split(JSON.Include.H[i], "/")
		fmt.Printf("Installing package %s...\n", pkg[len(pkg)-1])
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
		filename := strings.Split(JSON.Include.H[i], "/")
		err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
		res, err = http.Get(JSON.Include.C[i])
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
		filename = strings.Split(JSON.Include.C[i], "/")
		err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}