package cmd

import "github.com/pterm/pterm"

func Help() {
	cmd1 := "'cpkgs add <package-name>' - add C packages using cpkgs\n"
	cmd2 := "'cpkgs help' - shows this menu\n"
	cmd3 := "'cpkgs init <dir-name> [-d]' - initialize a new project using cpkgs\n"
	cmd4 := "'cpkgs install' - install all of the packages in the current project\n"
	cmd5 := "'cpkgs remove' - removes project directory\n"
	cmd6 := "'cpkgs run <file-name>' - run the file name using your selected compiler\n"
	cmd7 := "'cpkgs uninstall <package-name>' - removes the provided C package\n"
	cmd8 := "'cpkgs update <package-name> [-a]' - updates the provided C package\n"
	pterm.DefaultBox.WithTitle("List of all commands").WithTitleTopCenter().Print(cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8)
}
