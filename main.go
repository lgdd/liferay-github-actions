package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var dockerImages []DockerImage
	cloudWorkspace := os.Args[1]
	filepath.Walk(cloudWorkspace, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "LCP.json" {
			fmt.Println(fmt.Sprintf("Found LCP.json under %s", path))
			dockerImages = append(dockerImages, getDockerImageFromLCP(path))
		}
		return nil
	})
	fmt.Println("------")
	for _, dockerImage := range dockerImages {
		fmt.Println(fmt.Sprintf("Found Docker Images for %s in version %s", dockerImage.Repository, dockerImage.CurrentVersion))
	}
}

func getDockerImageFromLCP(lcpPath string) DockerImage {

	file, err := os.Open(lcpPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var lcp LCP
	json.Unmarshal(byteValue, &lcp)

	return newDockerImageFromTag(lcp.Image, lcpPath)
}

func getDockerImagesFromDockerfile(dockerfilePath string) []DockerImage {

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
			tag := lineSplit[1]
			dockerImages = append(dockerImages, newDockerImageFromTag(tag, dockerfilePath))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return dockerImages
}

func newDockerImageFromTag(tag string, dockerfilePath string) DockerImage {
	tagSplit := strings.Split(tag, ":")
	return DockerImage{
		Path:           dockerfilePath,
		Repository:     tagSplit[0],
		CurrentVersion: tagSplit[1],
	}
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
