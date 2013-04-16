/*
Package smuggol is a tool for physically importing Go code from one package to another.

While the normal Go import makes a "loose" connection to an external package (like a symlink/ln), a smuggol
import "physically" copies the package over (like cp).

    // Standard import, code is external and must resolved via "go get ..."
    import (
        . "github.com/robertkrimen/terst"
    )

    // After a smuggol import, code is local
    import (
        . "./terst"
    )

A smuggol import is useful for promoting code reuse while avoiding the breakage that can happen from
future changes in behavior/correctness of the import package.

WARNING: Currently this package is fairly alpha-ish (API-wise), and is only public to support
building of other tools (dbg-import, kilt-import, terst-import).

This package works by iterating through the .go files in a given import package and copying
them to a subordinate package in the new host package. The files have the following
comment at the top:

    // This file was AUTOMATICALLY GENERATED by ... (smuggol) from ...

The name of the subordinate package is the same as the original import package.

Additionally, supporting .go files can be generated in the host package at the same time. This is
done via a `map[string]string` , with each key/value pair representing a new file in the host package.
Before being written to disk, the value is processed through "text/template" as a template with the following
defined:

    HostPackage     # The name of the host package
    ImportPath      # The import path to the new import package
    ImportPackage   # The name of the new import package

*/
package smuggol

import (
	Flag "flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var (
	flag         = Flag.NewFlagSet("", Flag.ExitOnError)
	flag_update  = false
	flag_verbose = false
	flag_quiet   = false
	_            = func() byte {
		flag.BoolVar(&flag_update, "update", flag_update, "Update (go get -u) package first")
		flag.BoolVar(&flag_update, "u", flag_update, string(0))

		flag.BoolVar(&flag_verbose, "verbose", flag_verbose, "Be more verbose")
		flag.BoolVar(&flag_verbose, "v", flag_verbose, string(0))

		flag.BoolVar(&flag_quiet, "quiet", flag_quiet, "Be absolutely quiet")
		flag.BoolVar(&flag_quiet, "q", flag_quiet, string(0))
		return 0
	}()

	mainName = ""
	mainPkg  = ""

	_gofmt = true
)

// TODO: Package/file embedding
// smuggol: github.com/robertkrimen/kilt.GraveTrim =>
// github.com/robertkrimen/kilt/kilt.GraveTrim.go

func get(pkg string) error {
	arguments := []string{"get", "-u", "-v", pkg}
	if !flag_update {
		arguments = append(arguments[:1], arguments[2:]...)
	}
	cmd := exec.Command("go", arguments...)
	if !flag_quiet {
		fmt.Fprintf(os.Stdout, "# go get %s\n", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func relative(base, path string) (relativeBase, relativePath string) {
	tmp, _ := os.Getwd()
	relativeBase, _ = filepath.Rel(tmp, base)
	relativePath, _ = filepath.Rel(tmp, path)
	return
}

func main(dst string, src string, extra map[string]string) error {

	// We ignore the error because buildImport(src) below will barf, if necessary
	get(src)

	if dst == "" {
		dst = "."
	}

	dstPkg, err := buildImport(dst)
	dstBase := dst
	dstName := ""
	if err != nil {
		if len(extra) > 0 {
			if !flag_quiet {
				fmt.Fprintf(os.Stderr, "%s: unable to continue while missing Go package (in %s)\n", mainName, dst)
			}
			return err
		}
	} else {
		dstBase = dstPkg.Dir
		dstName = dstPkg.Name
	}

	srcPkg, err := buildImport(src)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(dstBase, srcPkg.Name)
	err = os.Mkdir(dstPath, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}

	relativeDstBase, relativeDstPath := relative(dstBase, dstPath)

	for _, file := range srcPkg.GoFiles {
		if !flag_quiet {
			fmt.Fprintf(os.Stdout, "+ %s\n", filepath.Join(relativeDstPath, file))
		}

		srcFile, err := os.Open(filepath.Join(srcPkg.Dir, file))
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(filepath.Join(dstPath, file))
		if err != nil {
			return err
		}
		defer dstFile.Close()

		fmt.Fprintf(dstFile, "// This file was AUTOMATICALLY GENERATED by %s (smuggol) from %s\n\n", mainName, mainPkg)

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	if len(extra) > 0 {
		importPkg, err := buildImport(dstPath)
		if err != nil {
			return err
		}

		importPath := importPkg.ImportPath
		if importPath == "." {
			// import "./<importPkg.Name>"
			importPath = "." + string(filepath.Separator) + importPkg.Name
		}

		data := map[string]string{
			"HostPackage":   dstName,
			"ImportPath":    importPath,
			"ImportPackage": importPkg.Name,
		}

		for name, tmpl := range extra {
			if !flag_quiet {
				fmt.Fprintf(os.Stdout, "+ %s\n", filepath.Join(relativeDstBase, name))
			}

			file, err := os.Create(filepath.Join(dstBase, name))
			if err != nil {
				return err
			}

			tmpl, err := template.New("").Parse(kiltGraveTrim(tmpl))
			if err != nil {
				return err
			}

			err = fmtPipe(func(output io.Writer) error {
				fmt.Fprintf(output, "// This file was AUTOMATICALLY GENERATED by %s (smuggol) for %s\n\n", mainName, mainPkg)
				return tmpl.Execute(output, data)
			}, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [target]\n", mainName)
	kilt.PrintDefaults(flag)
	fmt.Fprintf(os.Stderr, kilt.GraveTrim(`

    # Import %q into the current directory
    $ %s

    # Import %q into another directory
    $ %s ./xyzzy

    `), mainPkg, mainName, mainPkg, mainName)
}

// Main is the entry point for a command-line application.
//
// It takes three arguments:
//
// 1. The name of the application (for usage and error reporting, usually "<package>-import")
//
// 2. The import URL where the import package is located (e.g. "github.com/robertkrimen/terst")
//
// 3. A final, optional parameter (pass nil unless you know what you're doing)
//
// For example (terst-import):
//
//      func main() {
//          smuggol.Main("terst-import", "github.com/robertkrimen/terst", nil)
//      }
//
func Main(name, pkg string, extra map[string]string) {
	mainName = name
	mainPkg = pkg

	flag.Usage = usage
	flag.Parse(os.Args[1:])

	err := func() error {
		return main(flag.Arg(0), mainPkg, extra)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", mainName, err)
		os.Exit(1)
	}
}

func fmtPipe(input func(io.Writer) error, output io.Writer) error {

	inputOutput := output

	gofmt := false
	if _gofmt {
		// First, see if gofmt is available
		err := exec.Command("gofmt").Run()
		if err == nil {
			gofmt = true
		}
	}

	if !gofmt {
		return input(output)
	}

	cmd := exec.Command("gofmt")
	cmdStdin, err := cmd.StdinPipe()
	if err == nil {
		cmd.Stderr = os.Stderr
		cmd.Stdout = output

		err = cmd.Start()
		if err == nil {
			inputOutput = cmdStdin
		} else {
			cmdStdin.Close()
			cmd = nil
		}
	}

	err = input(inputOutput)
	if cmd != nil {
		cmdStdin.Close()
		if err == nil {
			err = cmd.Wait()
		}
	}
	return err
}
