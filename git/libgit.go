package git

import (
	"fmt"
	"log"
	"os"
	"time"

	git "gopkg.in/libgit2/git2go.v24"
)

func credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredSshKeyFromAgent("git")
	return git.ErrorCode(ret), &cred
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	// https://help.github.com/articles/what-are-github-s-ssh-key-fingerprints/
	gitHubRSAFingerprint := []byte{0x16, 0x27, 0xac, 0xa5, 0x76, 0x28, 0x2d, 0x36, 0x63, 0x1b, 0x56, 0x4d, 0xeb, 0xdf, 0xa6, 0x48}

	if hostname != "github.com" {
		log.Printf("remote hostname is '%s' instead of github.com", hostname)
		return git.ErrUser
	}

	for i := 0; i < 16; i++ {
		if cert.Hostkey.HashMD5[i] != gitHubRSAFingerprint[i] {
			log.Printf("Remote host certificate is invalid.\n")
			log.Printf("Expected: %x\n", gitHubRSAFingerprint)
			log.Printf("Got:      %x\n", cert.Hostkey.HashMD5)
			return git.ErrUser
		}
	}

	return 0
}

type libgitClient struct {
	repo *git.Repository
}

var _ RepoClient = &libgitClient{}

func newFromLibgit(dir string) (*libgitClient, error) {
	RemoteCallbacks := git.RemoteCallbacks{
		CertificateCheckCallback: certificateCheckCallback,
		CredentialsCallback:      credentialsCallback,
	}
	fetchOptions := &git.FetchOptions{RemoteCallbacks: RemoteCallbacks}
	cloneOptions := &git.CloneOptions{
		Bare:           false,
		CheckoutBranch: "master",
		FetchOptions:   fetchOptions,
	}

	repo, err := git.Clone("git@github.com:protonet/builds.git", dir, cloneOptions)
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	return &libgitClient{repo: repo}, nil
}

func (c *libgitClient) Close() {
	if c.repo != nil {
		c.repo = nil
	}
}

func (c *libgitClient) AddAndCommitChannel(channelName, commitMessage string) error {
	idx, err := c.repo.Index()
	if err != nil {
		return err
	}

	err = idx.AddByPath(fmt.Sprintf("%s.json", channelName))
	if err != nil {
		return err
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		return err
	}

	err = idx.Write()
	if err != nil {
		return err
	}

	tree, err := c.repo.LookupTree(treeID)
	if err != nil {
		return err
	}

	branch, err := c.repo.LookupBranch("master", git.BranchLocal)
	if err != nil {
		return err
	}

	commitTarget, err := c.repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	signature := &git.Signature{
		Name:  "Platform Tagger",
		Email: "engineering@protonet.info",
		When:  time.Now(),
	}

	_, err = c.repo.CreateCommit("refs/heads/master", signature, signature, commitMessage, tree, commitTarget)
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *libgitClient) Push() error {
	remote, err := c.repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	opts := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CertificateCheckCallback: certificateCheckCallback,
			CredentialsCallback:      credentialsCallback,
		}}
	err = remote.Push([]string{"refs/heads/master"}, opts)

	return err
}
