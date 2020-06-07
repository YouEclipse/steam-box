package steambox

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	steam "github.com/YouEclipse/steam-go/pkg"
	"github.com/google/go-github/github"
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

// GetPlayTime gets the paytime form steam web API.
func (b *Box) GetPlayTime(ctx context.Context, steamID uint64, appID ...uint32) ([]string, error) {
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

		line := pad(game.Name, " ", 33) + " " +
			pad(fmt.Sprintf("%d hrs %d mins", hours, mins), "", 16)
		lines = append(lines, line)
		max++
	}
	return lines, nil
}

func pad(s, pad string, targetLength int) string {
	padding := targetLength - utf8.RuneCountInString(s)
	if padding <= 0 {
		return s
	}

	return s + strings.Repeat(pad, padding)
}
