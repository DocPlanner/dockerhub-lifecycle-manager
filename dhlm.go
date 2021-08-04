package main

import (
	"dhlm/dockerhub"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Usage = "DockerHub lifecycle manager"
	app.UsageText = "dhlm [global options] [organization name] [repository name]"
	app.HideVersion = true

	var dhUsername string
	var dhPassword string
	var dhOrg string
	var dhRepo string
	var days string
	var dryRun bool

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username",
			Destination: &dhUsername,
		},
		cli.StringFlag{
			Name:        "password",
			Destination: &dhPassword,
		},
		cli.StringFlag{
			Name:        "days",
			Destination: &days,
			Value:       "30",
		},
		cli.BoolFlag{
			Name:        "dry-run",
			Destination: &dryRun,
		},
	}

	app.Action = func(c *cli.Context) error {
		dhOrg = c.Args().Get(0)
		dhRepo = c.Args().Get(1)

		if len(dhOrg) == 0 || len(dhRepo) == 0 {
			cli.ShowAppHelp(c)

			return nil
		}

		dh := dockerhub.NewClient(dockerhub.Auth{
			Username: dhUsername,
			Password: dhPassword,
		})

		daysInt, _ := strconv.Atoi(days)
		timeBefore := time.Now().Add(-time.Hour * 24 * time.Duration(daysInt))

		pageNumber := 1
		for tagsList := dh.GetImages(dhOrg, dhRepo, pageNumber, timeBefore); len(tagsList.Next) > 0; pageNumber++ {
			fmt.Println("Checking page:", pageNumber)
			var digests []string
			var ignoreList []*dockerhub.IgnoreWarnings
			for _, tag := range tagsList.Results {
				if tag.LastPulled.Unix() < timeBefore.Unix() {
					fmt.Println("Removing " + dhOrg + "/" + dhRepo + ":" + tag.Digest + " | " + tag.LastPulled.Format(time.RFC3339) + " | " + tag.LastPushed.Format(time.RFC3339))

					digests = append(digests, tag.Digest)
					var ignTags []string

					for _, t := range tag.Tags {
						if t.IsCurrent == true {
							ignTags = append(ignTags, t.Tag)
						}
					}

					ignoreList = append(ignoreList, &dockerhub.IgnoreWarnings{
						Repository: dhRepo,
						Digest:     tag.Digest,
						Warning:    "current_tag",
						Tags:       ignTags,
					})
				}
			}
			deletedImages := dh.DeleteImages(dhOrg, dhRepo, digests, timeBefore, dryRun, ignoreList)
			fmt.Printf("Summary of deleted images ➡ manifest_deletes: %d, manifest_errors: %d, tag_deletes: %d, tag_errors: %d \n",
				deletedImages.Metrics.ManifestDeletes, deletedImages.Metrics.ManifestErrors, deletedImages.Metrics.TagDeletes, deletedImages.Metrics.TagDeletes)

			tagsList = dh.GetImages(dhOrg, dhRepo, pageNumber, timeBefore)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
