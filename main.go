package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

func checkIfTokensPresent() {
	if len(os.Getenv("TOKEN_PLATFORM")) == 0 {
		log.Fatal("TOKEN_PLATFORM is not set")
	}

	if len(os.Getenv("TOKEN_PROTONET")) == 0 {
		log.Fatal("TOKEN_PROTONET is not set")
	}
}

func updateJSON(repo *buildsRepo, releaseNotesURL, tagTimestamp, isoTimestamp, targetChannel string, newBuildNumber int32, oldBuilds buildsData, commit bool) error {
	newBuilds := []buildsDatum{oldBuilds[0]}
	newBuilds[0].Build = newBuildNumber
	newBuilds[0].PublishedAt = isoTimestamp
	if releaseNotesURL != "" {
		newBuilds[0].URL = releaseNotesURL
	}

	for k := range newBuilds[0].Images {
		newBuilds[0].Images[k] = tagTimestamp
	}

	log.Printf("Old build version: %d", oldBuilds[0].Build)
	log.Printf("New build version: %d", newBuilds[0].Build)

	err := repo.saveChannel(targetChannel, newBuilds)
	if err != nil {
		return fmt.Errorf("Failed to save channel json: %s", err.Error())
	}

	if commit == true {
		commitMessage := fmt.Sprintf("release on channel '%s' at %s", targetChannel, isoTimestamp)
		err := repo.addAndCommitChannel(targetChannel, commitMessage)
		if err != nil {
			return err
		}

		err = repo.push()
		if err != nil {
			return err
		}
		log.Println("Push successful")
	} else {
		dump, _ := repo.dumpChannel(targetChannel)
		log.Printf("New JSON:\n%s\n", dump)
	}

	return nil
}

type taggerOptions struct {
	Commit    bool   `short:"c" long:"commit" description:"Commit the changes. Will make a dry run without this flag."`
	Build     int32  `short:"b" long:"build" required:"true" description:"Specify the build number to be placed inside the JSON."`
	SourceTag string `short:"s" long:"source-tag" default:"development" description:"Registry tag to be retagging from."`
	TargetTag string `short:"t" long:"target-tag" default:"soul3" description:"Registry tag to be retagging to."`
	URL       string `short:"u" long:"url" description:"Release notes URL"`
}

func main() {
	var opts taggerOptions

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	currentTime := time.Now().UTC()
	tagTimestamp := currentTime.Format("2006-01-02-1504")
	isoTimestamp := currentTime.Format("2006-01-02T15:04:05Z")

	repo, err := prepareRepo()
	if err != nil {
		log.Fatalf("Failed to clone the builds repo: %s", err.Error())
	}
	defer repo.Close()

	builds, err := repo.loadChannel(opts.SourceTag)
	if err != nil {
		log.Fatalf("Failed to load build data from channel '%s'", opts.SourceTag)
	}

	fmt.Printf("Tag timestamp: %s\n", tagTimestamp)
	fmt.Printf("ISO timestamp: %s\n", isoTimestamp)

	if opts.Commit == true {
		checkIfTokensPresent()

		err = retagAll(builds[0].Images, opts.SourceTag, tagTimestamp)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Dry run. Would otherwise retag following images from '%s' to '%s' and update channel '%s':\n", opts.SourceTag, tagTimestamp, opts.TargetTag)
		for k := range builds[0].Images {
			log.Printf(" * %s\n", k)
		}
	}

	err = updateJSON(repo, opts.URL, tagTimestamp, isoTimestamp, opts.TargetTag, opts.Build, builds, opts.Commit)
	if err != nil {
		log.Fatal(err)
	}
}
