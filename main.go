package main

import (
	"fmt"
	"log"
	"os"

	"github.com/taako-502/go-co-occurrence/cooccurrence"
	"github.com/taako-502/go-co-occurrence/scrape"
	"github.com/taako-502/go-co-occurrence/wikipedia"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <word>", os.Args[0])
	}

	word := os.Args[1]
	wikiResp, err := wikipedia.RetrieveWikipediaContent(word)
	if err != nil {
		log.Fatalf("Wikipedia API Error: %v", err)
	}

	var contents []string
	content, err := wikipedia.FetchWikipediaContent(wikiResp)
	if err != nil {
		log.Fatalf("Error fetching Wikipedia content: %v", err)
	}
	contents = append(contents, content)

	// 実際はGoogleの検索結果の上位10件をスクレイピングする
	urls := []string{
		"https://www.ozmall.co.jp/cosme/shampoo/",
		"https://www.cosme.net/categories/item/920/ranking/",
		"https://lalahair.co.jp/magazine/hair/shampoo/really-good/",
		"https://osusume.mynavi.jp/articles/5448/",
		"https://www.neuve-a.net/shop/pages/s_RM_shampoo.aspx",
		"https://www.cosme.com/products/ranking.php?category_id=721",
		"https://shampoo.kazukikishi.com/",
		"https://kakakumag.com/houseware/?id=19306",
	}

	// スクレイピング
	for _, url := range urls {
		content, err := scrape.FetchHtml(url)
		if err != nil {
			// エラーなら次
			continue
		}
		contents = append(contents, content)
	}

	var cooccurrences []cooccurrence.Cooccurrence
	for _, content := range contents {
		c := cooccurrence.CalculateCooccurrences(content, word)
		cooccurrences = cooccurrence.MergeCooccurrences(cooccurrences, c)
	}
	fmt.Printf("Co-occurring words with '%s':\n", word)
	for i, c := range cooccurrences {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, c.Word)
	}
}
