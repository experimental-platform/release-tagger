package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/libgit2/git2go.v24"
)

type buildsDatum struct {
	Build       int32             `json:"build"`
	Codename    string            `json:"codename"`
	URL         string            `json:"url"`
	PublishedAt string            `json:"published_at"`
	Images      map[string]string `json:"images"`
}

type buildsData []buildsDatum

type buildsRepo struct {
	directory string
	repo      *git.Repository
}

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

func prepareRepo() (*buildsRepo, error) {
	dir, err := ioutil.TempDir("", "tagger")
	if err != nil {
		return nil, err
	}

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

	return &buildsRepo{directory: dir, repo: repo}, nil
}

func (br *buildsRepo) Close() {
	if br.repo != nil {
		br.repo = nil
	}

	if br.directory != "" {
		os.RemoveAll(br.directory)
		br.directory = ""
	}
}

func (br *buildsRepo) GetDirectory() string {
	return br.directory
}

func (br *buildsRepo) addAndCommitChannel(channelName, commitMessage string) error {
	idx, err := br.repo.Index()
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

	tree, err := br.repo.LookupTree(treeID)
	if err != nil {
		return err
	}

	branch, err := br.repo.LookupBranch("master", git.BranchLocal)
	if err != nil {
		return err
	}

	commitTarget, err := br.repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	signature := &git.Signature{
		Name:  "Platform Tagger",
		Email: "engineering@protonet.info",
		When:  time.Now(),
	}

	_, err = br.repo.CreateCommit("refs/heads/master", signature, signature, commitMessage, tree, commitTarget)
	if err != nil {
		panic(err)
	}

	return nil
}

func (br *buildsRepo) push() error {
	remote, err := br.repo.Remotes.Lookup("origin")
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

func (br *buildsRepo) loadChannel(channelName string) (buildsData, error) {
	fileName := fmt.Sprintf("%s.json", channelName)
	filePath := path.Join(br.directory, fileName)

	var builds buildsData

	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawData, &builds)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

func (br *buildsRepo) saveChannel(channelName string, data buildsData) error {
	fileName := fmt.Sprintf("%s.json", channelName)
	filePath := path.Join(br.directory, fileName)

	rawData, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, rawData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (br *buildsRepo) dumpChannel(channelName string) (string, error) {
	fileName := fmt.Sprintf("%s.json", channelName)
	filePath := path.Join(br.directory, fileName)

	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
