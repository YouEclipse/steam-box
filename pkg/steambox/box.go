package steambox

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	steam "github.com/YouEclipse/steam-go/pkg"
	"github.com/google/go-github/github"
	"github.com/mattn/go-runewidth"
)

// Box defines the steam box.
type Box struct {
	steam  *steam.Client
	github *github.Client
}

// NewBox creates a new Box with the given API key.
func NewBox(apikey string, ghUsername, ghToken string) *Box {
	box := &Box{}
	box.steam = steam.NewClient(apikey, nil)
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(ghUsername),
		Password: strings.TrimSpace(ghToken),
	}

	box.github = github.NewClient(tp.Client())

	return box

}

// GetGist gets the gist from github.com.
func (b *Box) GetGist(ctx context.Context, id string) (*github.Gist, error) {
	gist, _, err := b.github.Gists.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return gist, nil
}

// UpdateGist updates the gist.
func (b *Box) UpdateGist(ctx context.Context, id string, gist *github.Gist) error {
	_, _, err := b.github.Gists.Edit(ctx, id, gist)
	return err
}

// GetPlayTime gets the top 5 Steam games played in descending order from the Steam API.
func (b *Box) GetPlayTime(ctx context.Context, steamID uint64, multiLined bool, appID ...uint32) ([]string, error) {
	params := &steam.GetOwnedGamesParams{
		SteamID:                steamID,
		IncludeAppInfo:         true,
		IncludePlayedFreeGames: true,
	}
	if len(appID) > 0 {
		params.AppIDsFilter = appID
	}

	gameRet, err := b.steam.IPlayerService.GetOwnedGames(ctx, params)
	if err != nil {
		return nil, err
	}
	var lines []string
	var max = 0
	sort.Slice(gameRet.Games, func(i, j int) bool {
		return gameRet.Games[i].PlaytimeForever > gameRet.Games[j].PlaytimeForever
	})

	for _, game := range gameRet.Games {
		if max >= 5 {
			break
		}

		hours := int(math.Floor(float64(game.PlaytimeForever / 60)))
		mins := int(math.Floor(float64(game.PlaytimeForever % 60)))

		if multiLined {
			gameLine := getNameEmoji(game.Appid, game.Name)
			lines = append(lines, gameLine)
			hoursLine := fmt.Sprintf("						    🕘 %d hrs %d mins", hours, mins)
			lines = append(lines, hoursLine)
		} else {
			line := pad(getNameEmoji(game.Appid, game.Name), " ", 35) + " " +
				pad(fmt.Sprintf("🕘 %d hrs %d mins", hours, mins), "", 16)
			lines = append(lines, line)
		}
		max++
	}
	return lines, nil
}

// GetRecentGames gets 5 recently played games from the Steam API.
func (b *Box) GetRecentGames(ctx context.Context, steamID uint64, multiLined bool) ([]string, error) {
	params := &steam.GetRecentlyPlayedGamesParams{
		SteamID: steamID,
		Count:   5,
	}

	gameRet, err := b.steam.IPlayerService.GetRecentlyPlayedGames(ctx, params)
	if err != nil {
		return nil, err
	}
	var lines []string
	var max = 0

	for _, game := range gameRet.Games {
		if max >= 5 {
			break
		}

		if game.Name == "" {
			game.Name = "Unknown Game"
		}

		hours := int(math.Floor(float64(game.PlaytimeForever / 60)))
		mins := int(math.Floor(float64(game.PlaytimeForever % 60)))

		if multiLined {
			gameLine := getNameEmoji(game.Appid, game.Name)
			lines = append(lines, gameLine)
			hoursLine := fmt.Sprintf("						    🕘 %d hrs %d mins", hours, mins)
			lines = append(lines, hoursLine)
		} else {
			line := pad(getNameEmoji(game.Appid, game.Name), " ", 35) + " " +
				pad(fmt.Sprintf("🕘 %d hrs %d mins", hours, mins), "", 16)
			lines = append(lines, line)
		}
		max++
	}
	return lines, nil
}

// UpdateMarkdown updates the content to the markdown file.
func (b *Box) UpdateMarkdown(ctx context.Context, title, filename string, content []byte) error {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("steambox.UpdateMarkdown: Error reade a file: %w", err)
	}

	start := []byte("<!-- steam-box start -->")
	before := md[:bytes.Index(md, start)+len(start)]
	end := []byte("<!-- steam-box end -->")
	after := md[bytes.Index(md, end):]

	newMd := bytes.NewBuffer(nil)
	newMd.Write(before)
	newMd.WriteString("\n" + title + "\n")
	newMd.WriteString("```text\n")
	newMd.Write(content)
	newMd.WriteString("\n")
	newMd.WriteString("```\n")
	newMd.WriteString("<!-- Powered by https://github.com/YouEclipse/steam-box . -->\n")
	newMd.Write(after)

	err = ioutil.WriteFile(filename, newMd.Bytes(), os.ModeAppend)
	if err != nil {
		return fmt.Errorf("steambox.UpdateMarkdown: Error writing a file: %w", err)
	}

	return nil
}

func pad(s, pad string, targetLength int) string {
	padding := targetLength - runewidth.StringWidth(s)

	if padding <= 0 {
		return s
	}

	return s + strings.Repeat(pad, padding)
}

func getNameEmoji(id int, name string) string {
	// hard code some game's emoji
	var nameEmojiMap = map[int]string{
		70:      "λ ",     // Half-Life
		220:     "λ² ",    // Half-Life 2
		500:     "🧟 ",     // Left 4 Dead
		550:     "🧟 ",     // Left 4 Dead 2
		570:     "⚔️ ",    // Dota 2
		730:     "🔫 ",     // CS:GO
		8930:    "🌏 ",     // Sid Meier's Civilization V
		252950:  "🚀 ",     // Rocket League
		269950:  "✈️ ",    // X-Plane 11
		271590:  "🚓 ",     // GTA 5
		359550:  "🔫 ",     // Tom Clancy's Rainbow Six Siege
		431960:  "💻 ",     // Wallpaper Engine
		578080:  "🍳 ",     // PUBG
		945360:  "🕵️‍♂️ ", // Among Us
		1250410: "🛩️ ",    // Microsoft Flight Simulator
		1091500: "🦾 ",     // Cyberpunk 2077
		594650:  "🎯 ",     // Hunt: Showdown
		230410:  "🐹 ",     // Warframe
		397540:  "🤖 ",     // Borderlands 3
		49520:   "🤖 ",     // Borderlands 2
	}

	if emoji, ok := nameEmojiMap[id]; ok {
		return emoji + name
	}

	if name == "Unknown Game" {
		return "❓ " + name
	}

	return "🎮 " + name
}
