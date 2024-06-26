package main

import (
	cmds "cpkgs/cmd"
	"flag"
	"log"
	"strings"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	if len(strings.TrimSpace(cmd)) <= 0 {
		cmd = "install"
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
