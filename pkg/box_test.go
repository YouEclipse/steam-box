package steambox

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestBox_GetPlayTime(t *testing.T) {
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

	box := NewBox(steamAPIKey, ghUsername, ghToken)
	lines, err := box.GetPlayTime(context.Background(), steamID, appIDList...)
	if err != nil {
		t.Error(err)
	}
	t.Log(strings.Join(lines, "\n"))

}
