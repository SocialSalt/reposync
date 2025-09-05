package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v3"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %+v", err)
	}
	if cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Debug("loaded config", "config", cfg)

	repoPath, err := CurrentRepoPath()
	if err != nil {
		log.Fatalf("Encountered an error when finding repo path: %+v", err)
	}
	projectName := filepath.Base(repoPath)
	slog.Debug(projectName)
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
	command := &cli.Command{
		Name:  "reposync",
		Usage: "Sync current git repo with remote host",
		Commands: []*cli.Command{
			{
				Name:  "push",
				Usage: "Send files to remote host",
				Flags: flags,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					host := cmd.Args().Get(0)
					dryRun := cmd.Bool("dry-run")
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					host := cmd.Args().Get(0)
					dryRun := cmd.Bool("dry-run")
					excludes := defaultExcludes()
					return executeRsync(
						strings.Join([]string{host, ":", projectName, "/"}, ""),
						repoPath,
						dryRun,
						excludes,
					)
				},
			},
		},
	}
	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func executeRsync(src string, target string, dryRun bool, excludes []string) error {
	fmt.Printf("rsync %s %s\n", src, target)
	formatted_excludes := formatExcludes(excludes)
	args := []string{
		src,
		target,
		"--verbose",
		"--archive",
		"--progress",
		"--human-readable",
		"--compress",
		"--itemize-changes",
		"--prune-empty-dirs",
		"--exclude-from=.gitignore",
	}
	args = append(args, formatted_excludes...)
	if dryRun {
		args = append(args, "--dry-run")
	}
	command := exec.Command(
		"rsync",
		args...,
	)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	err := command.Run()
	return err
}
