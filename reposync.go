package reposync

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	xdg "github.com/adrg/xdg"
	cli "github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Repo struct {
	Name     string   `yaml:"name,omitempty"`
	Excludes []string `yaml:"excludes,omitempty"`
}

type Config struct {
	GlobalExcludes []string `yaml:"global_excludes,omitempty"`
	Repos          []Repo   `yaml:"repos,omitempty"`
}

func defaultExcludes() []string {
	return []string{".git", "target", "__pyenv__", ".DS_Store"}
}

func loadConfig() (Config, error) {
	configFile := filepath.Join(xdg.ConfigHome, "reposync/reposync.yaml")

	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		cfg := Config{
			GlobalExcludes: defaultExcludes(),
		}
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return cfg, err
		}
		err = os.Mkdir(filepath.Dir(configFile), 0755)
		if err != nil {
			return cfg, err
		}
		err = os.WriteFile(configFile, data, 0755)
		if err != nil {
			return cfg, err
		}
		return cfg, nil

	} else {
		yamlFile, err := os.ReadFile(configFile)
		if err != nil {
			return Config{}, err
		}
		cfg := Config{}
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
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
