package wikipedia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const wikipediaAPIURL = "https://ja.wikipedia.org/w/api.php"

type WikipediaResponse struct {
	Batchcomplete bool `json:"batchcomplete"`
	Warnings      struct {
		Main struct {
			Warnings string `json:"warnings"`
		} `json:"main"`
		Revisions struct {
			Warnings string `json:"warnings"`
		} `json:"revisions"`
	} `json:"warnings"`
	Query struct {
		Pages []struct {
			Pageid    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Revisions []struct {
				Contentformat string `json:"contentformat"`
				Contentmodel  string `json:"contentmodel"`
				Content       string `json:"content"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

// retrieveWikipediaContent: Wikipedia API呼び出し
func RetrieveWikipediaContent(word string) (WikipediaResponse, error) {
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "revisions")
	params.Add("rvprop", "content")
	params.Add("titles", word)
	params.Add("formatversion", "2")

	url := fmt.Sprintf("%s?%s", wikipediaAPIURL, params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	var wikiResp WikipediaResponse
	if err := json.Unmarshal(body, &wikiResp); err != nil {
		log.Fatalf("Error: %v", err)
	}

	return wikiResp, nil
}

func FetchWikipediaContent(wikiResp WikipediaResponse) (string, error) {
	for _, pageInfo := range wikiResp.Query.Pages {
		if len(pageInfo.Revisions) > 0 {
			return pageInfo.Revisions[0].Content, nil
		}
	}
	return "", errors.New("can not fetch")
}
