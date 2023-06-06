package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	gh "github.com/cli/go-gh/v2"
	githubactions "github.com/sethvargo/go-githubactions"
)

const upgradeBranchName = "upgrade-liferay-cloud-images"

// Big thank you @balcsida for making a better regex: https://regex101.com/r/mk2TLg/
var cloudImagePattern = regexp.MustCompile(`^(\d+\.\d+\.\d+(-jdk\d+)?|^\d+\.\d+(-jdk\d+)?)(-\d+\.\d+\.\d+)?$`)

func main() {
	// gitConfigUser()
	// gitFetchAll()
	// mainBranchName := gitGetMainBranchName()
	fmt.Println("GITHUB_REF_NAME=" + os.Getenv("GITHUB_REF_NAME"))
	fmt.Println("githubactions.GetInput(\"workspace-directory\")=" + githubactions.GetInput("workspace-directory"))
	cloudWorkspace := "./cloud-repo"
	dockerImages := getDockerImagesFromLCPFiles(cloudWorkspace)
	dockerImagesToUpdate := getDockerImagesToUpdate(dockerImages)
	if len(dockerImagesToUpdate) > 0 {
		// gitSwitchBranch()
		for _, dockerImageToUpdate := range dockerImagesToUpdate {
			updateLCPFileWithLatestVersion(dockerImageToUpdate)
		}
		// gitCommitAndPush(cloudWorkspace)
		// pullRequestTitle := "[Liferay Cloud Upgrade] New versions for Docker images"
		// pullRequestBody := "New versions are available for Liferay Cloud Docker images"
		// createOrEditPullRequest(mainBranchName, pullRequestTitle, pullRequestBody)
	}
}

func gitConfigUser() {
	runCmd("git", "config", "user.name", "github-actions[bot]")
	runCmd("git", "config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com")
}

func gitFetchAll() {
	runCmd("git", "fetch", "--all")
	runCmd("git", "pull", "--all")
}

func gitSwitchBranch() {
	runCmd("git", "switch", "-c", upgradeBranchName)
}

func gitCommitAndPush(path string) {
	runCmd("git", "add", path)
	runCmd("git", "commit", "-m", "chore: upgrade liferay cloud images")
	runCmd("git", "push", "-u", "origin", upgradeBranchName)
}

func createOrEditPullRequest(mainBranchName, title, body string) {
	_, _, err := gh.Exec("pr", "edit", upgradeBranchName, "-t", title, "-b", body)
	if err != nil {
		gh.Exec("pr", "create", upgradeBranchName, "--base", mainBranchName, "--head", upgradeBranchName, "-t", title, "-b", body)
	} else {
		gh.Exec("pr", "reopen", upgradeBranchName)
	}
}

func gitGetMainBranchName() string {
	var stdoutBuffer bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	return stdoutBuffer.String()
}

func getDockerImagesToUpdate(dockerImages []DockerImage) []DockerImage {
	var dockerImagesToUpdate []DockerImage
	for _, dockerImage := range dockerImages {
		if latestDockerHubResult, err := fetchDockerHubResultForLatestStable(dockerImage); err == nil {
			dockerImage.DockerHubResult = latestDockerHubResult
			message := fmt.Sprintf("Found LCP.json using '%s' in version '%s' (latest is '%s')",
				dockerImage.Namespace+"/"+dockerImage.Repository, dockerImage.CurrentVersion, dockerImage.DockerHubResult.Name)
			fmt.Println(message)
			if dockerImage.CurrentVersion != dockerImage.DockerHubResult.Name {
				dockerImagesToUpdate = append(dockerImagesToUpdate, dockerImage)
			}
		} else {
			fmt.Println(err)
		}
	}
	return dockerImagesToUpdate
}

func getDockerImagesFromLCPFiles(rootPath string) []DockerImage {
	var dockerImages []DockerImage
	filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "LCP.json" {
			if dockerImage, err := getDockerImageFromLCP(path); err == nil {
				dockerImages = append(dockerImages, dockerImage)
			} else {
				fmt.Println(err)
			}
		}
		return nil
	})
	return dockerImages
}

func getDockerImageFromLCP(lcpPath string) (DockerImage, error) {

	file, err := os.Open(lcpPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var lcp LCP
	json.Unmarshal(byteValue, &lcp)

	if len(lcp.Image) > 0 {
		return newDockerImageFromTag(lcp.Image, lcpPath), nil
	}

	errorMessage := fmt.Sprintf("No Docker Image used for '%s' in %s", lcp.ID, lcpPath)
	return DockerImage{}, errors.New(errorMessage)
}

func newDockerImageFromTag(tag string, dockerfilePath string) DockerImage {
	tagSplit := strings.Split(tag, ":")
	namespaceRepositorySplit := strings.Split(tagSplit[0], "/")
	namespace := "library"
	repository := tagSplit[0]
	if len(namespaceRepositorySplit) > 1 {
		namespace = namespaceRepositorySplit[0]
		repository = namespaceRepositorySplit[1]
	}
	return DockerImage{
		Path:           dockerfilePath,
		Namespace:      namespace,
		Repository:     repository,
		CurrentVersion: tagSplit[1],
	}
}

func fetchDockerHubResultForLatestStable(dockerImage DockerImage) (DockerHubResult, error) {
	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://registry.hub.docker.com/v2/repositories/")
	urlBuilder.WriteString(dockerImage.Namespace)
	urlBuilder.WriteString("/")
	urlBuilder.WriteString(dockerImage.Repository)
	urlBuilder.WriteString("/tags?page_size=1024")

	resp, err := http.Get(urlBuilder.String())

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var response DockerHubResponse
	json.Unmarshal(body, &response)

	for _, result := range response.Results {
		if cloudImagePattern.MatchString(result.Name) {
			if strings.Contains(dockerImage.Repository, "dxp") {
				if result.Name[0:3] == dockerImage.CurrentVersion[0:3] {
					return result, nil
				}
			} else {
				return result, nil
			}
		}
	}

	return DockerHubResult{}, errors.New("no stable version found")
}

func updateLCPFileWithLatestVersion(dockerImage DockerImage) {
	imageName := dockerImage.Namespace + "/" + dockerImage.Repository
	oldImageValue := imageName + ":" + dockerImage.CurrentVersion
	newImageValue := imageName + ":" + dockerImage.DockerHubResult.Name
	read, err := ioutil.ReadFile(dockerImage.Path)

	if err != nil {
		panic(err)
	}

	newContents := strings.Replace(string(read), oldImageValue, newImageValue, -1)

	err = ioutil.WriteFile(dockerImage.Path, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}

func runCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}

// func getDockerImagesFromDockerfile(dockerfilePath string) []DockerImage {

// 	var dockerImages []DockerImage

// 	file, err := os.Open(dockerfilePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {
// 		if strings.HasPrefix(scanner.Text(), "FROM") {
// 			lineSplit := strings.Split(scanner.Text(), " ")
// 			tag := lineSplit[1]
// 			dockerImages = append(dockerImages, newDockerImageFromTag(tag, dockerfilePath))
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// 	return dockerImages
// }

type LCP struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

type DockerImage struct {
	Path            string
	Namespace       string
	Repository      string
	CurrentVersion  string
	DockerHubResult DockerHubResult
}

type DockerHubResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous any               `json:"previous"`
	Results  []DockerHubResult `json:"results"`
}

type DockerHubResult struct {
	Creator int `json:"creator"`
	ID      int `json:"id"`
	Images  []struct {
		Architecture string    `json:"architecture"`
		Features     string    `json:"features"`
		Variant      any       `json:"variant"`
		Digest       string    `json:"digest"`
		Os           string    `json:"os"`
		OsFeatures   string    `json:"os_features"`
		OsVersion    any       `json:"os_version"`
		Size         int       `json:"size"`
		Status       string    `json:"status"`
		LastPulled   time.Time `json:"last_pulled"`
		LastPushed   time.Time `json:"last_pushed"`
	} `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
	MediaType           string    `json:"media_type"`
	ContentType         string    `json:"content_type"`
	Digest              string    `json:"digest"`
}
