package smuggol

import (
	"go/build"
	"os"
	"path/filepath"
	"unicode"
)

/*
    This is how godoc does it:

	// Determine paths.
	//
	// If we are passed an operating system path like . or ./foo or /foo/bar or c:\mysrc,
	// we need to map that path somewhere in the fs name space so that routines
	// like getPageInfo will see it.  We use the arbitrarily-chosen virtual path "/target"
	// for this.  That is, if we get passed a directory like the above, we map that
	// directory so that getPageInfo sees it as /target.
	const target = "/target"
	const cmdPrefix = "cmd/"
	path := flag.Arg(0)
	var forceCmd bool
	var abspath, relpath string
	if filepath.IsAbs(path) {
		fs.Bind(target, OS(path), "/", bindReplace)
		abspath = target
	} else if build.IsLocalImport(path) {
		cwd, _ := os.Getwd() // ignore errors
		path = filepath.Join(cwd, path)
		fs.Bind(target, OS(path), "/", bindReplace)
		abspath = target
	} else if strings.HasPrefix(path, cmdPrefix) {
		path = path[len(cmdPrefix):]
		forceCmd = true
	} else if bp, _ := build.Import(path, "", build.FindOnly); bp.Dir != "" && bp.ImportPath != "" {
		fs.Bind(target, OS(bp.Dir), "/", bindReplace)
		abspath = target
		relpath = bp.ImportPath
	} else {
		abspath = pathpkg.Join(pkgHandler.fsRoot, path)
	}
	if relpath == "" {
		relpath = abspath
	}
*/
func buildImport(target string) (*build.Package, error) {
	if filepath.IsAbs(target) {
		return build.Default.ImportDir(target, 0)
	} else if build.IsLocalImport(target) {
		base, _ := os.Getwd()
		path := filepath.Join(base, target)
		return build.Default.ImportDir(path, 0)
	} else if pkg, _ := build.Default.Import(target, "", 0); pkg.Dir != "" && pkg.ImportPath != "" {
		return pkg, nil
	}
	path, _ := filepath.Abs(target) // Even if there is an error, still try?
	return build.Default.ImportDir(path, 0)
}

// For a package that exists outside of $GOROOT/$GOPATH:
//
// .SrcRoot == ""
// .ImportPath == "."
//
// Otherwise, these values are what you expect
//

func kiltGraveTrim(target string) string {
	// Discard \r? Go already does this for raw string literals.
	end := len(target)

	last := 0
	index := 0
	for index = 0; index < end; index++ {
		chr := rune(target[index])
		if chr == '\n' || !unicode.IsSpace(chr) {
			last = index
			break
		}
	}
	if index >= end {
		return ""
	}
	start := last
	if rune(target[start]) == '\n' {
		// Skip the leading newline
		start++
	}

	last = end - 1
	for index = last; index > start; index-- {
		chr := rune(target[index])
		if chr == '\n' || !unicode.IsSpace(chr) {
			last = index
			break
		}
	}
	stop := last
	result := target[start : stop+1]
	return result
}
