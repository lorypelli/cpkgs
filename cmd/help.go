package cmd

import "fmt"

func Help() {
	fmt.Println("List of all commands:")
	fmt.Print("\n")
	fmt.Println("---------------------------------------------------------------------------")
	fmt.Println("|'cpkgs add <package-name>' - add C packages using cpkgs                  |")
	fmt.Println("|'cpkgs help' - shows this menu                                           |")
	fmt.Println("|'cpkgs init <dir-name> [-d]' - initialize a new project using cpkgs      |")
	fmt.Println("|'cpkgs install' - install all of the packages in the current project     |")
	fmt.Println("|'cpkgs run <file-name>' - run the file name using your selected compiler |")
	fmt.Println("|'cpkgs update <package-name>' - updates the provided C package           |")
	fmt.Println("---------------------------------------------------------------------------")
}
