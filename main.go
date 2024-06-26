package main

import (
	cmds "cpkgs/cmd"
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	fmt.Print("\n")
	if len(strings.TrimSpace(cmd)) <= 0 {
		var c string
		fmt.Print("Provide a command to run: ")
		fmt.Scan(&c)
		cmd = c
	}
	switch cmd {
	case "add":
		{
			cmds.Add()
			break
		}
	case "help":
		{
			cmds.Help()
			break
		}
	case "init":
		{
			cmds.Init()
			break
		}
	case "install":
		{
			cmds.Install()
			break
		}
	case "run":
		{
			cmds.Run()
			break
		}
	default:
		{
			log.Fatal("Unknown command, to see all avaible commands type: 'cpkgs help' ")
		}
	}
}
