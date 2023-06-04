package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func main() {
	fmt.Println("hello world")
	filepath.Walk("./cloud-repo", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "Dockerfile" {
			fmt.Println(fmt.Sprintf("Found Dockerfile under %s", path))
		}
		if !info.IsDir() && info.Name() == "LCP.json" {
			fmt.Println(fmt.Sprintf("Found LCP.json under %s", path))
		}
		return nil
	})
}
