package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func Remove() {
	dir := flag.Arg(1)
	if len(strings.TrimSpace(dir)) <= 0 {
		fmt.Print("Provide directory to remove: ")
		fmt.Scan(&dir)
	}
	d, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = os.Stat(fmt.Sprintf("%s/cpkgs.json", dir))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Removing directory %s...\n", d.Name())
	err = os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Successfully removed directory %s!\n", d.Name())
}