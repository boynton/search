package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: search term <dir> ...")
		os.Exit(0)
	}
	term := os.Args[1]
	if len(os.Args) > 2 {
		for _, dir := range os.Args[2:] {
			search(term, dir)
		}
	} else {
		search(term, ".")
	}
}

func defaultFilter(info os.FileInfo) bool {
	name := info.Name()
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, "~") || strings.HasPrefix(name, "#") {
		return true
	}
	return false
}

func readDir(dir string, filter func(os.FileInfo) bool) ([]os.FileInfo, error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make([]os.FileInfo, 0, len(infos))
	for _, info := range infos {
		if filter == nil || !filter(info) {
			result = append(result, info)
		}
	}
	return result, nil
}

func search(term, dir string) {
	infos, err := readDir(dir, defaultFilter)
	if err != nil {
		return
	}
	for _, info := range infos {
		path := dir + "/" + info.Name()
		if info.IsDir() {
			search(term, path)
		} else {
			searchInFile(term, path)
		}
	}
}

func searchInFile(term string, path string) {
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}
	if !strings.HasSuffix(path, ".java") && !strings.HasSuffix(path, ".go") {
		return
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	s := string(bytes)
	if strings.Index(s, term) >= 0 {
		fmt.Printf("%s\n", path)
	}
}
