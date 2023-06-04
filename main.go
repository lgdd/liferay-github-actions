package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var dockerImages []DockerImage
	filepath.Walk("./cloud-repo", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "Dockerfile" {
			fmt.Println(fmt.Sprintf("Found Dockerfile under %s", path))
			dockerImages = append(dockerImages, getDockerImages(path)...)
		}
		if !info.IsDir() && info.Name() == "LCP.json" {
			fmt.Println(fmt.Sprintf("Found LCP.json under %s", path))
		}
		return nil
	})
	for _, dockerImage := range dockerImages {
		fmt.Println(fmt.Sprintf("Found Dockerfile using %s in version %s", dockerImage.Repository, dockerImage.CurrentVersion))
	}
}

func getDockerImages(dockerfilePath string) []DockerImage {

	var dockerImages []DockerImage

	file, err := os.Open(dockerfilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "FROM") {
			lineSplit := strings.Split(scanner.Text(), " ")
			tagSplit := strings.Split(lineSplit[1], ":")
			dockerImages = append(dockerImages, DockerImage{
				Path:           dockerfilePath,
				Repository:     tagSplit[0],
				CurrentVersion: tagSplit[1],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return dockerImages
}

type LCP struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

type DockerImage struct {
	Path           string
	Repository     string
	CurrentVersion string
	NewVersion     string
}
