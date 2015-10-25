package git

import (
	"os/exec"
	"strings"
)

type Repo struct {
	path string
}

type Version struct {
	Id string
	Author string
	Time string
	Message string
}

func New(path string) *Repo {
	return &Repo{
		path: path,
	}
}

func (r *Repo) gitCmd(args ...string) *exec.Cmd {
	finalArgs := append([]string{
		"-C", r.path,
	}, args...)
	return exec.Command("git", finalArgs...)
}

func (r *Repo) Versions() []Version {
	cmd := r.gitCmd(
		"log",
		"--all",
		"--pretty=format:%h|%an|%ar|%s")
	out, _ := cmd.Output()
	return parseLog(out)
}

func parseLog(out []byte) []Version {
	var versions []Version
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		version := parseCommit(line)
		versions = append(versions, version)
	}
	return versions
}

func parseCommit(line string) Version {
	fields := strings.Split(line, "|")
	return Version {
		Id: fields[0],
		Author: fields[1],
		Time: fields[2],
		Message: fields[3],
	}
}

func (r *Repo) FileVersion(path, version string) []byte {
	cmd := r.gitCmd(
		"show",
		version + ":" + path)
	out, _ := cmd.Output()
	return out
}

func (r *Repo) Fetch() {
	cmd := r.gitCmd("fetch")
	cmd.Run()
}

func (r *Repo) DiffToPrevious(version string) []byte {
	cmd := r.gitCmd(
		"diff",
		version + "^",
		version)
	out, _ := cmd.Output()
	return out
}