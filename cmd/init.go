package cmd

import (
	"cpkgs/pkg"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func Init() {
	var JSON pkg.JSON
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	j, _ := os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	json.Unmarshal(j, &JSON)
	d := flag.Arg(1)
	if d != "-d" && len(strings.TrimSpace(d)) > 0 {
		dir = flag.Arg(1)
		d = flag.Arg(2)
	}
	_, e := os.Stat(dir)
	if os.IsNotExist(e) {
		os.Mkdir(dir, 0777)
	}
	_, e = os.ReadFile(fmt.Sprintf("%s/cpkgs.json", dir))
	if os.IsNotExist(e) {
		err = os.WriteFile("cpkgs.json", j, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	err = os.WriteFile(fmt.Sprintf("%s/cpkgs.json", dir), j, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
	var compiler, filename string
	if d != "-d" {
		fmt.Print("Provide the compiler to use: ")
		fmt.Scanln(&compiler)
		if len(strings.TrimSpace(compiler)) <= 0 {
			fmt.Println("Using default...(gcc)")
			compiler = "gcc"
		}
		fmt.Print("Provide the output filename: ")
		fmt.Scanln(&filename)
		if len(strings.TrimSpace(filename)) <= 0 {
			fmt.Println("Using default...(out)")
			filename = "out"
		}
	} else {
		fmt.Println("Using defaults...(gcc, out)")
		compiler = "gcc"
		filename = "out"
	}
	JSON.Compiler = compiler
	JSON.FileName = filename
	j, err = json.Marshal(JSON)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = os.WriteFile(fmt.Sprintf("%s/cpkgs.json", dir), j, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
}
