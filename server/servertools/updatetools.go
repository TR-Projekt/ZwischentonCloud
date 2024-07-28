package servertools

import (
	"context"
	"errors"
	"os/exec"

	"github.com/google/go-github/v56/github"
)

func RunUpdate(currentVersion string, organisation string, repository string, updateScriptPath string) (string, error) {

	if currentVersion == "development" {
		return "", errors.New("this is a development server please update manually")
	}

	newestVersion, err := LatestVersion(organisation, repository)
	if err != nil {
		return "", errors.New("failed to retrieve latest release version with error: " + err.Error())
	}

	if newestVersion == currentVersion {
		return "", errors.New("this server is already up-to-date")
	}

	cmd := exec.Command("/bin/bash", "-c", "sudo "+updateScriptPath+" &")
	err = cmd.Start()
	if err != nil {
		return "", errors.New("Failed to run update script with error: " + err.Error())
	}

	return currentVersion + " => " + newestVersion, nil
}

func LatestVersion(organisation string, repository string) (string, error) {

	client := github.NewClient(nil)

	tags, _, err := client.Repositories.ListTags(context.Background(), organisation, repository, nil)
	if err != nil {
		return "", errors.New("Failed to fetch tags for repository: " + err.Error())
	}

	if len(tags) > 0 {
		latestTag := tags[0]
		return *latestTag.Name, nil
	}

	return "", errors.New("Failed to fetch tags for repository '" + repository + "' in organisation '" + organisation + "' with error: " + err.Error())
}
