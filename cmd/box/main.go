package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/YouEclipse/steam-box/pkg/steambox"
	"github.com/google/go-github/github"
)

func main() {
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	steamID, _ := strconv.ParseUint(os.Getenv("STEAM_ID"), 10, 64)
	appIDs := os.Getenv("APP_ID")
	appIDList := make([]uint32, 0)

	for _, appID := range strings.Split(appIDs, ",") {
		appid, err := strconv.ParseUint(appID, 10, 32)
		if err != nil {
			continue
		}
		appIDList = append(appIDList, uint32(appid))
	}

	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")
	gistID := os.Getenv("GIST_ID")

	updateOption := os.Getenv("UPDATE_OPTION") // options for update: GIST,MARKDOWN,GIST_AND_MARKDOWN
	markdownFile := os.Getenv("MARKDOWN_FILE") // the markdown filename

	var updateGist, updateMarkdown bool
	if updateOption == "MARKDOWN" {
		updateMarkdown = true
	} else if updateOption == "GIST_AND_MARKDOWN" {
		updateGist = true
		updateMarkdown = true
	} else {
		updateGist = true
	}

	box := steambox.NewBox(steamAPIKey, ghUsername, ghToken)

	ctx := context.Background()

	lines, err := box.GetPlayTime(ctx, steamID, appIDList...)
	if err != nil {
		panic("GetPlayTime err:" + err.Error())
	}

	filename := "🎮 Steam playtime leaderboard"

	if updateGist {
		gist, err := box.GetGist(ctx, gistID)
		if err != nil {
			panic("GetGist err:" + err.Error())
		}

		f := gist.Files[github.GistFilename(filename)]

		f.Content = github.String(strings.Join(lines, "\n"))
		gist.Files[github.GistFilename(filename)] = f

		err = box.UpdateGist(ctx, gistID, gist)
		if err != nil {
			panic("UpdateGist err:" + err.Error())
		}
	}

	if updateMarkdown && markdownFile != "" {
		title := filename
		if updateGist {
			title = fmt.Sprintf(`####  <a href="https://gist.github.com/%s" target="_blank">%s</a>`, gistID, title)
		}

		content := bytes.NewBuffer(nil)
		content.WriteString("\n")
		content.WriteString(strings.Join(lines, "\n"))

		err = box.UpdateMarkdown(ctx, title, markdownFile, content.Bytes())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("updating markdown successfully on", markdownFile)
	}
}
