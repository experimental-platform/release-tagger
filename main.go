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

func updateJSON(repo *buildsRepo, opts taggerOptions, tagTimestamp, isoTimestamp string) error {
	oldBuilds, err := repo.loadChannel(opts.SourceChannel)
	if err != nil {
		return err
	}

	newBuilds := []buildsDatum{oldBuilds[0]}
	if opts.Build != 0 {
		// if build number was given on commandline then set to it
		newBuilds[0].Build = opts.Build
	} else {
		destBuilds, err2 := repo.loadChannel(opts.TargetChannel)
		if err2 != nil {
			// if targetchannel doesn't exist, set to #1
			newBuilds[0].Build = 1
		} else {
			// otherwise increment
			newBuilds[0].Build = destBuilds[0].Build + 1
		}
	}
	newBuilds[0].PublishedAt = isoTimestamp
	if opts.URL != "" {
		newBuilds[0].URL = opts.URL
	}
	if opts.Codename != "" {
		newBuilds[0].Codename = opts.Codename
	}

	if opts.Retag {
		for k := range newBuilds[0].Images {
			newBuilds[0].Images[k] = tagTimestamp
		}
	}

	log.Printf("Old build version: %d", oldBuilds[0].Build)
	log.Printf("New build version: %d", newBuilds[0].Build)

	err = repo.saveChannel(opts.TargetChannel, newBuilds)
	if err != nil {
		return fmt.Errorf("Failed to save channel json: %s", err.Error())
	}

	if opts.Commit == true {
		commitMessage := fmt.Sprintf("release on channel '%s' at %s", opts.TargetChannel, isoTimestamp)
		err := repo.addAndCommitChannel(opts.TargetChannel, commitMessage)
		if err != nil {
			return err
		}

		err = repo.push()
		if err != nil {
			return err
		}
		log.Println("Push successful")
	} else {
		dump, _ := repo.dumpChannel(opts.TargetChannel)
		log.Printf("New JSON:\n%s\n", dump)
	}

	return nil
}

type taggerOptions struct {
	Commit        bool   `short:"c" long:"commit" description:"Commit the changes. Will make a dry run without this flag."`
	Build         int32  `short:"b" long:"build" required:"false" default:"0" description:"Specify the build number to be placed inside the JSON."`
	SourceChannel string `short:"s" long:"source-channel" default:"development" description:"Release channel to be retagging/copying from."`
	TargetChannel string `short:"t" long:"target-channel" default:"soul3" description:"Release channel to be retagging to."`
	URL           string `short:"u" long:"url" description:"Release notes URL"`
	Codename      string `short:"n" long:"codename" description:"Release codename"`
	Copy          bool   `long:"copy" description:"Copy from one release channel to another"`
	Retag         bool   `long:"retag" description:"Retag the source channel's images with a timestamp (e.g. for tagging from 'development')"`
}

func retaggingStep(images map[string]string, opts *taggerOptions, tagTimestamp string) {
	if opts.Commit == true {

		checkIfTokensPresent()

		err := retagAll(images, opts.SourceChannel, tagTimestamp)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Printf("Dry run. Would otherwise retag following images from '%s' to '%s' and update channel '%s':\n", opts.SourceChannel, tagTimestamp, opts.TargetChannel)
		for k := range images {
			log.Printf(" * %s\n", k)
		}
	}
}

func parseOptions(opts *taggerOptions) {
	parser := flags.NewParser(opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if opts.Copy && opts.Retag {
		fmt.Fprintln(os.Stderr, "--copy and --retag are mutually exclusive.")
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if !opts.Copy && !opts.Retag {
		fmt.Fprintln(os.Stderr, "You must select either --copy or --retag option.")
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
}

func main() {
	var opts taggerOptions

	parseOptions(&opts)

	currentTime := time.Now().UTC()
	tagTimestamp := currentTime.Format("2006-01-02-1504")
	isoTimestamp := currentTime.Format("2006-01-02T15:04:05Z")
	fmt.Printf("Tag timestamp: %s\n", tagTimestamp)
	fmt.Printf("ISO timestamp: %s\n", isoTimestamp)

	repo, err := prepareRepo()
	if err != nil {
		log.Fatalf("Failed to clone the builds repo: %s", err.Error())
	}
	defer repo.Close()

	builds, err := repo.loadChannel(opts.SourceChannel)
	if err != nil {
		log.Fatalf("Failed to load build data from channel '%s'", opts.SourceChannel)
	}

	// skip this step if merely copying a channel over
	if opts.Retag {
		retaggingStep(builds[0].Images, &opts, tagTimestamp)
	}

	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	if err != nil {
		log.Fatal(err)
	}
}
