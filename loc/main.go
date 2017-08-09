//This program is a variation of the example in chapter 8 of 
//Donovan & Kernighan's book and counts lines of code in folders.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var done = make(chan struct{})
var seen = make(map[string]bool)
var extensions = []string {".go", ".c", ".h", ".cpp", ".hpp", ".java"}
var verbose bool

func cancel() {
	close(done)
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func init() {
	var extensionsString string

	flag.StringVar(&extensionsString, "e", ".go,.c,.h,.cpp,.hpp,.cs,.vb,.java,.html,.css,.js", "allowed file extensions")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	extensions = strings.Split(extensionsString, ",")
}

func main() {
	//Determine folders:
	roots := flag.Args()
	if len(roots) == 0 {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
   		if err != nil {
            log.Fatal(err)
    	}
    	
		roots = []string{dir}
	}

	//Terminate if input is detected:
	go func() {
		os.Stdin.Read(make([]byte, 1))
		cancel()
	}()

	n := 0

	for _, root := range roots {
		n += loc(root)
	}

	if (cancelled()) {
		return
	}

	fmt.Printf("%d lines of code\n", n)
}

func loc(file string) int {
	fi, err := os.Stat(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loc: %v\n", err)
		return 0
	}

	if fi.IsDir() {
		return locInDir(file)
	}else{
		return locInFile(file)
	}
}

func locInDir(dir string) (n int) {
	if cancelled() {
		return
	}

	if seen[dir] {
		return
	}
	seen[dir] = true

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loc: %v\n", err)
		return
	}
	defer f.Close()

	entries, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loc: %v\n", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			n += locInDir(filepath.Join(dir, entry.Name()))
		}else{
			n += locInFile(filepath.Join(dir, entry.Name()))
		}
	}

	return
}

func locInFile(file string) (n int) {
	if cancelled() {
		return
	}

	if seen[file] {
		return
	}
	seen[file] = true

	res := len(extensions) == 0

	for _, ext := range extensions {
		if !strings.HasSuffix(file, ext) {
			 continue
		}

		res = true
		break
	}

	if !res {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loc: %v\n", err)
		return 
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		n++
	}

	if verbose {
		fmt.Printf("%5d %s\n", n, file)
	}

	return
}