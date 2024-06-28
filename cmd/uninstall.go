package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/pkg"
)

func Uninstall() {
	var JSON pkg.JSON
	j, err := os.ReadFile("cpkgs.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(j, &JSON)
	pkgs := flag.Args()[1:]
	if len(pkgs) <= 0 {
		fmt.Print("Provide packages to uninstall: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		p := scanner.Text()
		pkgs = strings.Split(p, " ")
	}
	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg, ".h") {
			fmt.Printf("%s is not a valid header file, skipping...\n", pkg)
			continue
		}
		fmt.Printf("Removing package %s...\n", pkg)
		err := os.Remove(fmt.Sprintf("cpkgs/%s", pkg))
		if err != nil {
			log.Fatal(err)
			return
		}
		err = os.Remove(fmt.Sprintf("cpkgs/%s", strings.ReplaceAll(pkg, ".h", ".c")))
		if err != nil {
			log.Fatal(err)
			return
		}
		for i, h := range JSON.Include.H {
			header := strings.Split(h, "/")
			h := header[len(header)-1]
			if h == pkg {
				JSON.Include.H = append(JSON.Include.H[:i], JSON.Include.H[i+1:]...)
				JSON.Include.C = append(JSON.Include.C[:i], JSON.Include.C[i+1:]...)
			}
		}
		j, err := json.MarshalIndent(JSON, "", "  ")
		if err != nil {
			log.Fatal(err)
			return
		}
		err = os.WriteFile("cpkgs.json", j, 0777)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Successfully removed package %s!\n", pkg)
	}
}
