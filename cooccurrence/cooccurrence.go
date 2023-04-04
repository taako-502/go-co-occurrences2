package cooccurrence

import (
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/taako-502/go-co-occurrence/cooccurrence/internal/utils"
	"golang.org/x/exp/slices"
)

type Cooccurrence struct {
	Word  string
	Point int
}

// filterString: 文字列を正規表現でフィルタリングし、フィルタにマッチしたらTrueを返却
func filterString(str string, regexArr []string) bool {
	for _, regexStr := range regexArr {
		regex := regexp.MustCompile(regexStr)
		if regex.MatchString(str) {
			return true
		}
	}
	return false
}

func CalculateCooccurrences(content, targetWord string) []Cooccurrence {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalf("Error parsing Wikipedia content: %v", err)
	}

	words := make(map[string]int)
	t := tokenizer.New() // 形態素解析器を作成
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		tokens := t.Tokenize(text) // 形態素解析を実行
		for i := 0; i < len(tokens)-1; i++ {
			token := tokens[i].Surface
			if token == targetWord { // TODO: 分割されるタイプの単語もあるので注意（例: 「東京都」→「東京」「都」）
				// キーワードから最も近い言葉に高いポイントを与える
				// キーワードから2番目に近い言葉に2番目に高いポイントを与える
				// ・・・
				// キーワードからn番目に近い言葉にn番目に高いポイントを与える
				point := 3 // 最大ポイント
				for j := i + 1; j < len(tokens); j++ {
					str := tokens[j].Surface // トークン
					pos := tokens[j].Pos()   // 品詞
					if str == token ||
						slices.Contains(utils.RemovePosFilter(), pos) ||
						slices.Contains(utils.NotCooccurrencesFilter(), strings.TrimSpace(str)) ||
						filterString(str, utils.RemoveRgx()) {
						continue
					}
					nextToken := tokens[j].Surface
					words[nextToken] = words[nextToken] + point
					point-- // ポイントを与え終わったらデクリメントする
					if point == 0 {
						// 与えるポイントが無くなったら終了
						break
					}
				}
			}
		}
	})

	cooccurrences := make([]Cooccurrence, 0, len(words))
	for word, count := range words {
		cooccurrences = append(cooccurrences, Cooccurrence{Word: word, Point: count})
	}

	sort.Slice(cooccurrences, func(i, j int) bool {
		return cooccurrences[i].Point > cooccurrences[j].Point
	})

	return cooccurrences
}

func MergeCooccurrences(list1 []Cooccurrence, list2 []Cooccurrence) []Cooccurrence {
	merged := make(map[string]int)
	for _, cooc := range list1 {
		merged[cooc.Word] += cooc.Point
	}
	for _, cooc := range list2 {
		merged[cooc.Word] += cooc.Point
	}
	result := make([]Cooccurrence, 0, len(merged))
	for word, point := range merged {
		result = append(result, Cooccurrence{Word: word, Point: point})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Point > result[j].Point
	})
	return result
}
