# libgit2-checkout-head

Demonstrates the difference between `git checkout --force HEAD` and [`git_checkout_head()`](https://libgit2.github.com/libgit2/#HEAD/group/checkout/git_checkout_head). It actually uses [git2go](https://github.com/libgit2/git2go) instead of [libgit2](https://github.com/libgit2/libgit2), for convenience only.

There's a simple [main.go](main.go) file that does the following:

- initializes an origin repo with a single README file in it and master branch
- using both the command line git and libgit2 (git2go), it mirrors the origin repo created above by doing the following:
  - initializes an empty repo
  - creates the **same** README file as in the origin repo
  - adds the origin remote pointing at the origin repo created above and fetches
  - sets the `.git/HEAD` to the commit pointed to by `refs/remotes/origin/HEAD`
  - does a `git checkout --force HEAD` to sync the new repo's worktree with the origin

After the last step, the repos should be the same. However with libgit2 the README file still shows as deleted and untracked.

Note: This only happens if the README file has the same contents. If it's different, than the checkout works fine (by overwriting the one in the second repo).

## Usage

The easiest way to try it is by using [Docker](https://www.docker.com) (in the root of the repo):

```sh
docker-compose up
```

The output should be similar to this:

```
app_1  | $ git init /tmp/origin179161192
app_1  | Initialized empty Git repository in /tmp/origin179161192/.git/
app_1  | $ bash -c echo 'test' > README
app_1  | $ git add .
app_1  | $ git commit -m Initial commit
app_1  | [master (root-commit) c0f56af] Initial commit
app_1  |  1 file changed, 1 insertion(+)
app_1  |  create mode 100644 README
app_1  | --------------------------------------------------------------------------
app_1  | Initializing repo using command git
app_1  | $ git init /tmp/cmd077320871
app_1  | Initialized empty Git repository in /tmp/cmd077320871/.git/
app_1  | $ bash -c echo 'test' > README
app_1  | $ git remote add origin /tmp/origin179161192
app_1  | $ git fetch origin +HEAD:refs/remotes/origin/HEAD +refs/heads/*:refs/remotes/origin/*
app_1  | From /tmp/origin179161192
app_1  |  * [new ref]                    -> origin/HEAD
app_1  |  * [new branch]      master     -> origin/master
app_1  | $ git show-ref
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7 refs/remotes/origin/HEAD
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7 refs/remotes/origin/master
app_1  | $ bash -c echo $(git rev-parse refs/remotes/origin/HEAD) > .git/HEAD
app_1  | $ cat .git/HEAD
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7
app_1  | $ git status --short --branch
app_1  | ## HEAD (no branch)
app_1  | D  README
app_1  | ?? README
app_1  | $ git checkout --force HEAD
app_1  | $ git status --short --branch
app_1  | ## HEAD (no branch)
app_1  | --------------------------------------------------------------------------
app_1  | Initializing repo using libgit2
app_1  | $ (initialized empty repo using libgit2 in /tmp/libgit2546713562)
app_1  | $ bash -c echo 'test' > README
app_1  | $ git show-ref
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7 refs/remotes/origin/HEAD
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7 refs/remotes/origin/master
app_1  | $ cat .git/HEAD
app_1  | c0f56af237e24d93369323bb51bb85df29a621f7
app_1  | $ git status --short --branch
app_1  | ## HEAD (no branch)
app_1  | D  README
app_1  | ?? README
app_1  | $ (force checkout HEAD using libgit2)
app_1  | $ git status --short --branch
app_1  | ## HEAD (no branch)
app_1  | D  README
app_1  | ?? README
```

In the command line repo (`/tmp/cmd077320871`) the last `git status` displays a clean worktree after the checkout:

```
$ git status --short --branch
## HEAD (no branch)
```

but in the libgit2 repo (`/tmp/libgit2546713562`) the worktree remains dirty after `git_checkout_head()` with `--force`:

```
$ git status --short --branch
## HEAD (no branch)
D  README
?? README
```

## Additional Resources

- [git_checkout_head() doesn't do anything (in practice) #2864](https://github.com/libgit2/libgit2/issues/2864) on GitHub
- [git2go's CheckoutHead() not updating the index](http://stackoverflow.com/questions/34599073/git2gos-checkouthead-not-updating-the-index) on Stack Overflow