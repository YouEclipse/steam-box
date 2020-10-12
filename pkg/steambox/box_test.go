package steambox

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestBox_GetPlayTime(t *testing.T) {
	var err error
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	steamID, _ := strconv.ParseUint(os.Getenv("STEAM_ID"), 10, 64)

	multiLined := false // boolean for whether hours should have their own line
	if os.Getenv("MULTILINE") != "" {
		multiLined, err = strconv.ParseBool(os.Getenv("MULTILINE"))
		if err != nil {
			panic("multiLined option error: "+ err.Error())
		}
	}

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

	box := NewBox(steamAPIKey, ghUsername, ghToken)
	lines, err := box.GetPlayTime(context.Background(), steamID, multiLined, appIDList...)
	if err != nil {
		t.Error(err)
	}
	t.Log(strings.Join(lines, "\n"))
}

func TestBox_GetRecentGames(t *testing.T) {
	var err error
	steamAPIKey := os.Getenv("STEAM_API_KEY")
	steamID, _ := strconv.ParseUint(os.Getenv("STEAM_ID"), 10, 64)

	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")

	multiLined := false // boolean for whether hours should have their own line - YES = true, NO = false
	if os.Getenv("MULTILINE") != "" {
		lineOption := os.Getenv("MULTILINE")
		if lineOption == "YES" {
			multiLined = true
		}
	}

	box := NewBox(steamAPIKey, ghUsername, ghToken)
	lines, err := box.GetRecentGames(context.Background(), steamID, multiLined)
	if err != nil {
		t.Error(err)
	}
	t.Log(strings.Join(lines, "\n"))
}

func TestBox_Readme(t *testing.T) {

	ghToken := os.Getenv("GH_TOKEN")
	ghUsername := os.Getenv("GH_USER")

	box := NewBox("", ghUsername, ghToken)

	ctx := context.Background()

	filename := "test.md"
	title := `####  <a href="https://gist.github.com/YouEclipse/9bc7025496e478f439b9cd43eba989a4" target="_blank">🎮 Steam playtime leaderboard</a>`
	content := []byte(`🔫 Counter-Strike: Global Offensive  🕘 1546 hrs 25 mins
🚓 Grand Theft Auto V                🕘 52 hrs 15 mins
💻 Wallpaper Engine                  🕘 39 hrs 59 mins
🍳 PLAYERUNKNOWN'S BATTLEGROUNDS     🕘 34 hrs 40 mins
🌏 Sid Meier's Civilization V        🕘 11 hrs 9 mins`)

	err := box.UpdateMarkdown(ctx, title, filename, content)
	if err != nil {
		t.Error(err)
	}
	c, _ := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", c)
}
