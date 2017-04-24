package git

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type BuildsDatum struct {
	Build       int32             `json:"build"`
	Codename    string            `json:"codename"`
	URL         string            `json:"url"`
	PublishedAt string            `json:"published_at"`
	Images      map[string]string `json:"images"`
}

type BuildsData []BuildsDatum

type RepoClient interface {
	Close()
	AddAndCommitChannel(channelName, commitMessage string) error
	Push() error
}

type BuildsRepo struct {
	directory string
	client    RepoClient
}

func PrepareRepo(gitClient string) (*BuildsRepo, error) {
	dir, err := ioutil.TempDir("", "tagger")
	if err != nil {
		return nil, err
	}

	var c RepoClient

	switch gitClient {
	case "libgit":
		c, err = newFromLibgit(dir)
		if err != nil {
			return nil, err
		}
	case "command":
		c, err = newFromCommand(dir)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown git client '%s'", gitClient)
	}

	return &BuildsRepo{directory: dir, client: c}, nil
}

func (br *BuildsRepo) Close() {
	if br.client != nil {
		br.client.Close()
		br.client = nil
	}

	if br.directory != "" {
		os.RemoveAll(br.directory)
		br.directory = ""
	}
}

func (br *BuildsRepo) GetDirectory() string {
	return br.directory
}

func (br *BuildsRepo) AddAndCommitChannel(channelName, commitMessage string) error {
	return br.client.AddAndCommitChannel(channelName, commitMessage)
}

func (br *BuildsRepo) Push() error {
	return br.client.Push()
}

func (br *BuildsRepo) LoadChannel(channelName string) (BuildsData, error) {
	fileName := fmt.Sprintf("%s.json", channelName)
	filePath := path.Join(br.directory, fileName)

	var builds BuildsData

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

func (br *BuildsRepo) SaveChannel(channelName string, data BuildsData) error {
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

func (br *BuildsRepo) DumpChannel(channelName string) (string, error) {
	fileName := fmt.Sprintf("%s.json", channelName)
	filePath := path.Join(br.directory, fileName)

	rawData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
