package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
)

func main() {
	repoPath, err := CurrentRepoPath()
	if err != nil {
		log.Fatalf("Encountered an error when finding repo path: %+v", err)
	}
	projectName := filepath.Base(repoPath)
	log.Println(projectName)
	if err := checkForRsync(); err != nil {
		log.Fatalf("Failed to find rsync command")
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %+v", err)
	}
	log.Printf("%#v", cfg)

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
			Name:  "pull",
			Usage: "Get files from remote host",
			Flags: flags,
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

func executeRsync(src string, target string, dryRun bool, excludes []string) error {
	log.Printf("rsync %s %s", src, target)
	formatted_excludes := formatExcludes(excludes)
	args := []string{
		src,
		target,
		"-az",
		"--verbose",
		"--prune-empty-dirs",
		"--filter=':-.gitignore'",
		"--human-readable",
		"--progress",
		"--itemize-changes",
	}
	args = append(args, formatted_excludes...)
	if dryRun {
		args = append(args, "--dry-run")
	}
	command := exec.Command(
		"rsync",
		args...,
	)
	// TODO: setup stdout and stderr pipes
	err := command.Run()
	return err
}
