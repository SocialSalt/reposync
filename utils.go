package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func checkForRsync() error {
	command := exec.Command(
		"rsync",
		"--version",
	)
	err := command.Run()
	return err
}

func CurrentRepoPath() (string, error) {
	var repoPath string
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirToSearch := cwd
	for repoPath == "" {
		_, err := os.Stat(filepath.Join(dirToSearch, ".git"))
		if err == nil {
			return dirToSearch, nil
		}
		dirToSearch = filepath.Dir(dirToSearch)
		if dirToSearch == "/" || dirToSearch == "." {
			repoPath = cwd
			log.Println("failed to find a git repo, using current working dir")
		}
	}
	return repoPath, nil
}

func formatExcludes(excludes []string) (formatted_excludes []string) {
	formatted_excludes = make([]string, len(excludes))
	for i := range excludes {
		formatted_excludes[i] = fmt.Sprintf("--exclude=%s", excludes[i])
	}
	return
}
