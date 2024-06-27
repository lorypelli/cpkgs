package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
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
	for _, h := range JSON.Include.H {
		res, err := http.Get(h)
		pkg := strings.Split(h, "/")
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
		filename := strings.Split(h, "/")
		err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
		c := strings.ReplaceAll(h, ".h", ".c")
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
		filename = strings.Split(c, "/")
		err = os.WriteFile(fmt.Sprintf("%s/cpkgs/%s", dir, filename[len(filename)-1]), body, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
