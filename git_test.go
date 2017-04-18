package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestPrepareRepo(t *testing.T) {
	repo, err := prepareRepo("", "", "")
	assert.Nil(t, err)
	defer repo.Close()

	dir := repo.GetDirectory()
	info, err := os.Stat(dir)
	assert.Nil(t, err, "Directory '%s' does not exist!", dir)
	assert.True(t, info.IsDir(), "'%s' is not a directory!", dir)

	gitDir := dir + "/.git"

	info2, err := os.Stat(gitDir)
	assert.Nil(t, err, "Directory '%s' does not exist!", gitDir)
	assert.True(t, info2.IsDir(), "'%s' is not a directory!", gitDir)
}

func TestAddAndCommit(t *testing.T) {
	repo, err := prepareRepo("", "", "")
	assert.Nil(t, err)
	defer repo.Close()

	dir := repo.GetDirectory()
	cmd1 := exec.Command("git", "-C", dir, "show-ref", "refs/heads/master")
	out1, err := cmd1.CombinedOutput()
	assert.Nil(t, err)
	ref1 := strings.Split(string(out1), " ")[0]

	r := make([]byte, 64)
	_, err = rand.Read(r)
	assert.Nil(t, err)
	randomData := fmt.Sprintf("%x", r)
	filePath := fmt.Sprintf("%s/%s.json", dir, randomData)

	ioutil.WriteFile(filePath, []byte(randomData), 0644)

	err = repo.addAndCommitChannel(randomData, "foobar commit")

	cmd2 := exec.Command("git", "-C", dir, "show-ref", "refs/heads/master")
	out2, err := cmd2.CombinedOutput()
	assert.Nil(t, err)
	ref2 := strings.Split(string(out2), " ")[0]

	assert.NotEqual(t, ref1, ref2, "Repository's master ref should have changed")

	_, err = os.Stat(filePath)
	assert.Nil(t, err, "File '%s' does not exist", filePath)
}
