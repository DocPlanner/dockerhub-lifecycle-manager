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

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "username",
			Destination: &dhUsername,
		},
		cli.StringFlag{
			Name: "password",
			Destination: &dhPassword,
		},
		cli.StringFlag{
			Name: "days",
			Destination: &days,
			Value: "30",
		},
	}

	app.Action = func(c *cli.Context) error {
		dhOrg = c.Args().Get(0)
		dhRepo = c.Args().Get(1)

		if len(dhOrg) == 0 || len(dhRepo) == 0 {
			cli.ShowAppHelp(c)

			return nil
		}

		dh := dockerhub.NewClient(dockerhub.Auth {
			Username: dhUsername,
			Password: dhPassword,
		})

		daysInt, _ := strconv.Atoi(days)
		timeBefore := time.Now().Add(-time.Hour*24*time.Duration(daysInt))

		pageNumber := 1
		for tagsList := dh.GetTags(dhOrg, dhRepo, pageNumber); len(tagsList.Next) > 0; pageNumber++ {
			fmt.Println("Checking page:", pageNumber)
			for _, tag := range tagsList.Results {
				if tag.LastUpdated.Unix() < timeBefore.Unix() {
					fmt.Println("Removing "+dhOrg+"/"+dhRepo+":"+tag.Name+" | "+tag.LastUpdated.Format(time.RFC822))
					dh.DeleteTag(dhOrg, dhRepo, tag.Name)
				}
			}

			tagsList = dh.GetTags(dhOrg, dhRepo, pageNumber)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}