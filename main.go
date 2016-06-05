package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"fmt"
	"strings"

	"gopkg.in/libgit2/git2go.v24"
)

var (
	checkoutOpts = git.CheckoutOpts{
		Strategy: git.CheckoutForce,
	}
)

func main() {

	// Create the origin repo

	origin, err := ioutil.TempDir("", "origin")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(origin)
	Cmd{}.Exec("git", "init", origin)
	Cmd{origin}.Exec("bash", "-c", "echo 'test' > README")
	Cmd{origin}.Exec("git", "add", ".")
	Cmd{origin}.Exec("git", "commit", "-m", "Initial commit")

	// Create repo using command line

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("Initializing repo using command git")
	cmdPath, err := ioutil.TempDir("", "cmd")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(cmdPath)
	Cmd{}.Exec("git", "init", cmdPath)
	Cmd{cmdPath}.Exec("bash", "-c", "echo 'test' > README")
	Cmd{cmdPath}.Exec("git", "remote", "add", "origin", origin)
	Cmd{cmdPath}.Exec("git", "fetch", "origin", "+HEAD:refs/remotes/origin/HEAD", "+refs/heads/*:refs/remotes/origin/*")
	Cmd{cmdPath}.Exec("git", "show-ref")
	Cmd{cmdPath}.Exec("bash", "-c", "echo $(git rev-parse refs/remotes/origin/HEAD) > .git/HEAD")
	Cmd{cmdPath}.Exec("cat", ".git/HEAD")
	Cmd{cmdPath}.Exec("git", "status", "--short", "--branch")
	Cmd{cmdPath}.Exec("git", "checkout", "--force", "HEAD")
	Cmd{cmdPath}.Exec("git", "status", "--short", "--branch")

	// Create repo using libgit2 (git2go)

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("Initializing repo using libgit2")
	libgit2Path, err := ioutil.TempDir("", "libgit2")
	if err != nil {
		log.Fatal(err)
	}
	repo, err := git.InitRepository(libgit2Path, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("$ (initialized empty repo using libgit2 in %s)\n", libgit2Path)
	Cmd{libgit2Path}.Exec("bash", "-c", "echo 'test' > README")
	remote, err := repo.Remotes.Create("origin", origin)
	if err != nil {
		log.Fatal(err)
	}
	refspec := []string{"+HEAD:refs/remotes/origin/HEAD", "+refs/heads/*:refs/remotes/origin/*"}
	if err := remote.Fetch(refspec, &git.FetchOptions{}, ""); err != nil {
		log.Fatal(err)
	}
	Cmd{libgit2Path}.Exec("git", "show-ref")
	if err := repo.SetHead("refs/remotes/origin/HEAD"); err != nil {
		log.Fatal(err)
	}
	Cmd{libgit2Path}.Exec("cat", ".git/HEAD")
	Cmd{libgit2Path}.Exec("git", "status", "--short", "--branch")
	fmt.Println("$ (force checkout HEAD using libgit2)")
	if err := repo.CheckoutHead(&checkoutOpts); err != nil {
		log.Fatal(err)
	}
	Cmd{libgit2Path}.Exec("git", "status", "--short", "--branch")
}

type Cmd struct {
	Dir string
}

func (cmd Cmd) Exec(args ...string) {
	fmt.Println("$", strings.Join(args, " "))
	c := exec.Command(args[0], args[1:]...)
	c.Dir = cmd.Dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
