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
	for i := 0; i < len(JSON.Include.H); i++ {
		res, err := http.Get(JSON.Include.H[i])
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
