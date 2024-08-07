package cmd

import (
	"flag"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
)

func Remove() {
	dir := flag.Arg(1)
	if len(strings.TrimSpace(dir)) <= 0 {
		dir, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Provide directory to remove").Show()
	}
	d, err := os.Stat(dir)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	if !d.IsDir() {
		if pkgs := flag.Args()[1:]; len(pkgs) > 0 {
			pterm.Warning.Println("You provided header files to the command, 'cpkgs uninstall' will be executed instead!")
			cmd := pterm.Sprintf("cpkgs uninstall %s", strings.Join(pkgs, " "))
			cmdExec := exec.Command("bash", "-c", cmd)
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
	if _, err := os.Stat(pterm.Sprintf("%s/cpkgs.json", dir)); err != nil {
		pterm.Error.Println(err)
		return
	}
	choice, _ := pterm.DefaultInteractiveConfirm.WithDefaultText(pterm.Sprintf("Do you want to remove the %s directory?", dir)).Show()
	if choice {
		pterm.Info.Printfln("Removing %s directory...", d.Name())
		if err := os.RemoveAll(dir); err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Success.Printfln("Successfully removed %s directory!", d.Name())
	}
}
