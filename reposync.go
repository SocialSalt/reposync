package reposync

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
)

type Repo struct {
	Name     string
	Excludes []string
}

type Config struct {
	GlobalExclude []string
	Repos         []Repo
}

func defaultExcludes() []string {
	return []string{".git", "target", "__pyenv__", ".DS_Store"}
}

func App() {
	repoPath, err := CurrentRepoPath()
	if err != nil {
		log.Fatalf("Encountered an error when finding repo path: %+v", err)
	}
	projectName := filepath.Base(repoPath)
	log.Println(projectName)
	if err := checkForRsync(); err != nil {
		log.Fatalf("Failed to find rsync command")
	}

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"d", "n"},
			Usage:   "Add --dry-run, -n flag to rsync command",
		},
	}
	commands := []*cli.Command{
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Send files to remote host",
			Flags:   flags,
			Action: func(ctx *cli.Context) error {
				host := ctx.Args().Get(0)
				dryRun := ctx.Bool("dry-run")
				excludes := defaultExcludes()

				return executeRsync(
					repoPath,
					strings.Join([]string{host, ":", projectName, "/"}, ""),
					dryRun,
					excludes,
				)
			},
		},
		{
			Name:    "pull",
			Aliases: []string{"y"},
			Usage:   "Get files from remote host",
			Flags:   flags,
			Action: func(ctx *cli.Context) error {
				host := ctx.Args().Get(0)
				dryRun := ctx.Bool("dry-run")
				excludes := defaultExcludes()

				return executeRsync(
					strings.Join([]string{host, ":", projectName, "/"}, ""),
					repoPath,
					dryRun,
					excludes,
				)
			},
		},
	}
	app := &cli.App{
		Name:     "reposync",
		Usage:    "Sync current git repo with remote host",
		Commands: commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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

func checkForRsync() error {
	command := exec.Command(
		"rsync",
		"--version",
	)
	err := command.Run()
	return err
}

func executeRsync(src string, target string, dryRun bool, excludes []string) error {
	command := exec.Command(
		"rsync",
		"--version",
	)
	err := command.Run()
	return err
}
