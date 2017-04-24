package git

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type gitCommandClient struct {
	dir string
}

var _ RepoClient = &gitCommandClient{}

func newFromCommand(dir string) (*gitCommandClient, error) {
	url := "git@github.com:protonet/builds.git"
	cmd := exec.Command("git", "clone", url, dir)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return &gitCommandClient{dir: dir}, cmd.Run()
}

func (c *gitCommandClient) Close() {
	if c.dir != "" {
		c.dir = ""
	}
}

func (c *gitCommandClient) AddAndCommitChannel(channelName, commitMessage string) error {
	fileName := fmt.Sprintf("%s.json", channelName)
	addParams := []string{"--git-dir", path.Join(c.dir, ".git"), "--work-tree", c.dir, "add", fileName}
	addCmd := exec.Command("git", addParams...)
	addCmd.Stderr = os.Stderr
	addCmd.Stdin = os.Stdin
	addCmd.Stdout = os.Stdout
	err := addCmd.Run()
	if err != nil {
		return err
	}

	commitParams := []string{"--git-dir", path.Join(c.dir, ".git"), "--work-tree", c.dir, "commit", "-m", commitMessage}
	commitCmd := exec.Command("git", commitParams...)
	commitCmd.Stderr = os.Stderr
	commitCmd.Stdin = os.Stdin
	commitCmd.Stdout = os.Stdout
	err = commitCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *gitCommandClient) Push() error {
	params := []string{"--git-dir", path.Join(c.dir, ".git"), "--work-tree", c.dir, "push"}
	cmd := exec.Command("git", params...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
