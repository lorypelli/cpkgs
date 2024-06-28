package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Remove() {
	dir := flag.Arg(1)
	if len(strings.TrimSpace(dir)) <= 0 {
		fmt.Print("Provide directory to remove: ")
		fmt.Scan(&dir)
	}
	if strings.HasSuffix(dir, ".h") {
		pkgs := flag.Args()[1:]
		if len(pkgs) > 0 {
			fmt.Println("You provided header files to the command, 'cpkgs uninstall' will be executed instead!")
			cmd := fmt.Sprintf("cpkgs uninstall %s", strings.Join(pkgs, " "))
			cmdExec := exec.Command("sh", "-c", cmd)
			if runtime.GOOS == "windows" {
				cmdExec = exec.Command("cmd", "/C", cmd)
			}
			cmdExec.Stdin = os.Stdin
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Run()
			return
		}
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
	var choice string
	fmt.Printf("Do you want to remove the %s directory? (y/n) ", dir)
	fmt.Scan(&choice)
	if strings.ToLower(choice) == "y" {
		fmt.Printf("Removing directory %s...\n", d.Name())
		err = os.RemoveAll(dir)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("Successfully removed directory %s!\n", d.Name())
	}
}
