# Reposync
reposync is a (mostly unnecessary) wrapper around [rsync](https://linux.die.net/man/1/rsync) that will sync the current git repo you are in to a remote machine.
The program will find the git repo that contains the current working directory and send it to a remote machine.
If it fails to find a git repo then it will send the current working directory.

## Usage
Usage: `reposync [push|pull] HOST [OPTIONS]`

The `HOST` can be a host alias as defined in [your `~/.ssh/config`.](https://www.digitalocean.com/community/tutorials/how-to-configure-custom-connection-options-for-your-ssh-client) or it can be a `User@Hostname`.

### Options
`--dry-run, -d, -n`  Add --dry-run, -n flag to rsync command (default: false)

## Config

The config is stored in a `yaml` file and the structure of the config object is as follows:
```go
type Repo struct {
	Name     string   `yaml:"name,omitempty"`
	Excludes []string `yaml:"excludes,omitempty"`
}

type Config struct {
	GlobalExcludes []string `yaml:"global_excludes,omitempty"`
	Repos          []Repo   `yaml:"repos,omitempty"`
}
```
reposync will always ignore patterns that appear in the `global_excludes` list, but it also allows the user to exclude certain patterns per repo.

Here is an example config file.

```yaml
global_excludes:
- .git
- target
- __pycache__
- .DS_Store

repos:
  - name: reposync
    exclues: ["*.yaml", "LICENSE"]
  - name: todo-project
    excludes: [".venv"]
```
With this config file reposync will always ignore the `.git` dir, and dirs named `target` or `__pycache__` or `.DS_Store`.

## Behavior

Running `reposync push hostname` from a git repo named `project` will create the directory in `$HOME/project` on the remote host and send all files not ignored by `rsync`, comparable to `rsync -az /path/to/project/ hostname:project/`.
The command `reposync pull hostname` will reverse the `SRC` and `DEST` to get files from the remote host, comparable `rsync -az hostname:project/ /path/to/project/`).
The only available flag is `-d, --dry-run` which will add the `-n, --dry-run` flag to the `rsync` command.

If there is no config file found, then `reposync` will exclude the following directories `.git, target, .DS_Store, __pycache__`.

`reposync` uses the following flags on `rsync`

```console
-az
--verbose
--prune-empty-dirs
--exclude-from=.gitignore
--human-readable
--progress
--itemize-changes
```
Noe that we use the `--exclude-from=.gitignore` so reposync will also ignore any pattern in the `.gitignore` file.


## Why should you use it

`reposync` automatticaly finds the root of your current project.
If for example you are in `/path/to/project/some/sub/dir` and you run `reposync push host` it will sync from `/path/to/project` as long as `/path/to/project/.git` exists.

`reposync` allows for easy management of different `rsync` exclude settings between different repos, so you don't need to copy and paste different `rsync` commands or edit the same command based on which repo you want to sync.
