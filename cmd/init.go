package cmd

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/lorypelli/cpkgs/v2/internal"
	"github.com/pterm/pterm"
)

func Init() {
	var JSON internal.JSON
	dir, err := os.Getwd()
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	d := flag.Arg(1)
	if d != "-d" && d != "--default" && len(strings.TrimSpace(d)) > 0 {
		dir = flag.Arg(1)
		d = flag.Arg(2)
	}
	var language, compiler, filename string
	if d != "-d" && d != "--default" {
		language, _ = pterm.DefaultInteractiveSelect.WithDefaultText("Provide language").WithDefaultOption("C").WithOptions([]string{"C", "C++"}).Show()
		compiler, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Provide the compiler to use").WithDefaultValue("gcc").Show()
		if len(strings.TrimSpace(compiler)) <= 0 {
			pterm.Info.Println("Using default...(gcc)")
			compiler = "gcc"
		}
		filename, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Provide the output filename").WithDefaultValue("out").Show()
		if len(strings.TrimSpace(filename)) <= 0 {
			pterm.Info.Println("Using default...(out)")
			filename = "out"
		}
	} else {
		pterm.Info.Println("Using defaults...(C, gcc, out)")
		language = "C"
		compiler = "gcc"
		filename = "out"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			pterm.Error.Println(err)
			return
		}
	}
	JSON.Schema = "https://raw.githubusercontent.com/lorypelli/cpkgs/main/schemas/schema.json"
	JSON.Language = language
	JSON.Compiler = compiler
	JSON.FileName = filename
	JSON.Include = internal.Include{
		C:   []string{},
		CPP: []string{},
		H:   []string{},
		HPP: []string{},
	}
	JSON.CPPExtensions = nil
	if language == "C++" {
		var code, header string
		pterm.Warning.Println("You selected C++ so you will need to provide additional options")
		code, _ = pterm.DefaultInteractiveSelect.WithDefaultText("Provide code files extension").WithDefaultOption(".cpp").WithOptions([]string{".cpp", ".cc", ".cxx", ".c++", ".cp"}).Show()
		header, _ = pterm.DefaultInteractiveSelect.WithDefaultText("Provide header files extension").WithDefaultOption(".h").WithOptions([]string{".h", ".hpp", ".hh", ".hxx", ".h++", ".hp"}).Show()
		JSON.CPPExtensions = &internal.CPPExtensions{
			Code:   code,
			Header: header,
		}
	}
	j, err := json.MarshalIndent(JSON, "", "  ")
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	if err := os.WriteFile(pterm.Sprintf("%s/cpkgs.json", dir), j, 0644); err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Success.Println("Successfully created cpkgs.json file with the following settings:")
	pterm.Info.Printfln("Language -> %s", language)
	pterm.Info.Printfln("Compiler -> %s", compiler)
	pterm.Info.Printfln("Filename -> %s", filename)
	if err := os.RemoveAll("cpkgs"); err != nil {
		pterm.Error.Println(err)
		return
	}
	pterm.Info.Printfln("Creating cache directory at: %s...", internal.GetCacheDir())
	if _, err := os.Stat(internal.GetCacheDir()); os.IsNotExist(err) {
		if err := os.Mkdir(internal.GetCacheDir(), 0755); err != nil {
			pterm.Error.Println(err)
			return
		}
		pterm.Success.Println("Cache directory successfully created!")
	} else {
		pterm.Warning.Println("Cache directory already exists, nothing was changed!")
	}
}
