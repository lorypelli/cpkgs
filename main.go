package main

import (
	"flag"
	"log"
	"strings"

	cmds "github.com/lorypelli/cpkgs/v2/cmd"
	"github.com/pterm/pterm"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	if len(strings.TrimSpace(cmd)) <= 0 {
		pterm.Info.Println("Welcome to cpkgs!")
		pterm.Info.Println("Type 'cpkgs help' to see all avaible commands")
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
	case "remove":
		{
			cmds.Remove()
			break
		}
	case "run":
		{
			cmds.Run()
			break
		}
	case "uninstall":
		{
			cmds.Uninstall()
			break
		}
	case "update":
		{
			cmds.Update()
			break
		}
	default:
		{
			log.Fatal("Unknown command, to see all avaible commands type: 'cpkgs help' ")
			break
		}
	}
}
