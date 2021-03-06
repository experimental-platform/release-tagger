package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/experimental-platform/release-tagger/git"
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

func updateJSON(repo *git.BuildsRepo, opts taggerOptions, tagTimestamp, isoTimestamp string) error {
	var (
		Retag bool
	)

	switch opts.Args.Action {
	case "copy":
		break
	case "create":
		Retag = true
		break
	default:
		fmt.Fprintf(os.Stderr, "The only allowed actions are 'copy' and 'create'\n")
		os.Exit(1)
		break
	}

	oldBuilds, err := repo.LoadChannel(opts.Args.SourceChannel)
	if err != nil {
		return err
	}

	newBuilds := []git.BuildsDatum{oldBuilds[0]}
	if opts.Build != 0 {
		// if build number was given on commandline then set to it
		newBuilds[0].Build = opts.Build
	} else {
		destBuilds, err2 := repo.LoadChannel(opts.Args.TargetChannel)
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

	if Retag {
		for k := range newBuilds[0].Images {
			newBuilds[0].Images[k] = tagTimestamp
		}
	}

	log.Printf("Old build version: %d", oldBuilds[0].Build)
	log.Printf("New build version: %d", newBuilds[0].Build)

	err = repo.SaveChannel(opts.Args.TargetChannel, newBuilds)
	if err != nil {
		return fmt.Errorf("Failed to save channel json: %s", err.Error())
	}

	if opts.Commit == true {
		commitMessage := fmt.Sprintf("release on channel '%s' at %s", opts.Args.TargetChannel, isoTimestamp)
		err := repo.AddAndCommitChannel(opts.Args.TargetChannel, commitMessage)
		if err != nil {
			return err
		}

		err = repo.Push()
		if err != nil {
			return err
		}
		log.Println("Push successful")
	} else {
		dump, _ := repo.DumpChannel(opts.Args.TargetChannel)
		log.Printf("New JSON:\n%s\n", dump)
	}

	return nil
}

type taggerOptionsArgs struct {
	Action        string `description:"either 'copy' or 'create'"`
	SourceChannel string `description:"Release channel to be creating/copying from."`
	TargetChannel string `description:"Release channel to be creating/copying to."`
}

type taggerOptions struct {
	Args taggerOptionsArgs `positional-args:"true" required:"true"`

	Commit    bool   `short:"c" long:"commit" description:"Commit the changes. Will make a dry run without this flag."`
	Build     int32  `short:"b" long:"build" required:"false" default:"0" description:"Specify the build number to be placed inside the JSON."`
	URL       string `short:"u" long:"url" description:"Release notes URL"`
	Codename  string `short:"n" long:"codename" description:"Release codename"`
	GitClient string `long:"git-client" default:"libgit" description:"Git client. Either 'libgit' or 'command'"`
}

func retaggingStep(images map[string]string, opts *taggerOptions, tagTimestamp string) {
	if opts.Commit == true {

		checkIfTokensPresent()

		err := retagAll(images, opts.Args.SourceChannel, tagTimestamp)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		log.Printf("Dry run. Would otherwise create following tags from '%s' to '%s' and update channel '%s':\n", opts.Args.SourceChannel, tagTimestamp, opts.Args.TargetChannel)
		for k := range images {
			log.Printf(" * %s\n", k)
		}
	}
}

func parseOptions(opts *taggerOptions) {
	parser := flags.NewParser(opts, flags.Default)
	_, err := parser.Parse()

	if err != nil {
		// this condition prevents the help from being printed twice when specifically requested by the -h|--help parameter
		if flagserr, ok := err.(*flags.Error); !ok || flagserr.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
		}
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

	repo, err := git.PrepareRepo(opts.GitClient)
	if err != nil {
		log.Fatalf("Failed to clone the builds repo: %s", err.Error())
	}
	defer repo.Close()

	builds, err := repo.LoadChannel(opts.Args.SourceChannel)
	if err != nil {
		log.Fatalf("Failed to load build data from channel '%s'", opts.Args.SourceChannel)
	}

	// skip this step if merely copying a channel over
	if opts.Args.Action == "create" {
		retaggingStep(builds[0].Images, &opts, tagTimestamp)
	}

	err = updateJSON(repo, opts, tagTimestamp, isoTimestamp)
	if err != nil {
		log.Fatal(err)
	}
}
