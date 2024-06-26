package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Init() {
	var JSON pkg.JSON
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	d := flag.Arg(1)
	if d != "-d" && d != "--default" && len(strings.TrimSpace(d)) > 0 {
		dir = flag.Arg(1)
		d = flag.Arg(2)
	}
	var compiler, filename string
	if d != "-d" && d != "--default" {
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
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			log.Fatal(err)
			return
		}
	}
	JSON.Compiler = compiler
	JSON.FileName = filename
	JSON.Include = pkg.Include{
		C: []string{},
		H: []string{},
	}
	j, err := json.MarshalIndent(JSON, "", "  ")
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := os.WriteFile(fmt.Sprintf("%s/cpkgs.json", dir), j, 0777); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfully created cpkgs.json file with the following settings:")
	fmt.Print("\n")
	fmt.Printf("Compiler -> %s\n", compiler)
	fmt.Printf("Filename -> %s\n", filename)
}
