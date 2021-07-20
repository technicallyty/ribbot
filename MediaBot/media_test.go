package MediaBot

import (
	"github.com/technicallyty/vidbot/redditbot"
	"testing"
)

func TestFullDownloadCombine(t *testing.T) {
	url := "https://www.reddit.com/r/ContagiousLaughter/comments/oipjmp/some_people_decide_to_use_this_guys_car_for_a/"
	vbot := redditbot.NewVidBot(url)
	_, _, err := vbot.Download()
	if err != nil {
		t.Fatal(err)
	}

}
