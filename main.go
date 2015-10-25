package main

import (
	"flag"
	"time"
	"log"
	"net/http"
	"github.com/barnacs/sgt/git"
	"github.com/barnacs/sgt/server"
)

var (
	path = flag.String("repo", ".", "repo path")
	addr = flag.String("addr", ":8001", "<host:port>")
	pollFq = flag.Duration("fetch", 0, "fetch frequency (default 0 - disabled)")
)

func main() {
	flag.Parse()
	repo := git.New(*path)
	if (*pollFq > 0) {
		go pollRepo(repo, *pollFq)
	}
	err := http.ListenAndServe(*addr, server.New(repo))
	log.Fatal(err)
}

func pollRepo(repo *git.Repo, period time.Duration) {
	ticker := time.NewTicker(period)
	for range ticker.C {
		repo.Fetch()
	}
}