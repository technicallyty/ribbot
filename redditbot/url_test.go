package redditbot

import (
	"fmt"
	"strings"
	"testing"
)

func TestDeriveJSONURL(t *testing.T) {
	url := "https://www.reddit.com/r/redditdev/comments/ihgmv5/getting_audio_from_reddit_video"
	jsonURL := deriveJSONURL(url)
	if !strings.Contains(jsonURL, ".json") {
		t.Fatal("string should have .json at the end")
	}
}

func TestValidateRedditLink(t *testing.T) {
	url := "https://www.facebook.com/r/redditdev/comments/ihgmv5/getting_audio_from_reddit_video/"
	matched := redditRegExp.MatchString(url)
	if matched {
		t.Fatal("should not match an invalid string")
	}

	goodURL := "https://www.reddit.com/r/iamatotalpieceofshit/comments/oi7n6r/and_the_match_hasnt_even_started_yet/"
	matched = redditRegExp.MatchString(goodURL)
	if !matched {
		t.Fatal("valid string should be ok")
	}
}

func TestUrlToDirName(t *testing.T) {
	url := "https://www.reddit.com/r/ContagiousLaughter/comments/oipjmp/some_people_decide_to_use_this_guys_car_for_a/"
	anotherURL := "https://www.reddit.com/r/ContagiousLaughter/comments/oipjmp/some_people_decide_to_use_this_guys_car_for_a"
	f, _ := urlToResourceName(url)
	z, _ := urlToResourceName(anotherURL)
	if f != "ContagiousLaughter_oipjmp_some_people_decide_to_use_this_guys_car_for_a" || f != z {
		t.Fatal("should've made a valid dirname")
	}
}

func TestMarshalURL(t *testing.T) {
	url := "https://www.reddit.com/r/aww/comments/oj5n80/this_is_truffles_the_cat_a_stray_found_by_a/"
	url = deriveJSONURL(url)
	url, err := fetchMediaURL(url)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
}

func TestValidateURL(t *testing.T) {

}
