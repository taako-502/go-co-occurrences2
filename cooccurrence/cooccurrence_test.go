package cooccurrence

import (
	"reflect"
	"sort"
	"testing"
)

func TestCalculateCooccurrences(t *testing.T) {
	testCases := []struct {
		name       string
		content    string
		targetWord string
		expected   []Cooccurrence
	}{
		{
			name:       "日本語の場合",
			content:    "<p>侍は日本のファイターです。彼らは武士、もののふ、忍者などと呼ばれます。日本刀の達人であり、居合と呼ばれる剣技で相対する敵を切りつけます。</p>",
			targetWord: "侍",
			expected: []Cooccurrence{
				{"日本", 3},
				{"ファイター", 2},
				{"彼ら", 1},
			},
		},
		{
			name:       "２つ連なった単語の場合",
			content:    "<p>薩英戦争は薩摩藩と大英帝国の間で発生した戦闘です。薩摩藩は一部被害を受けましたが、大英帝国の軍艦を撃破しました。</p>",
			targetWord: "薩英戦争",
			expected:   []Cooccurrence{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateCooccurrences(tc.content, tc.targetWord)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}

func TestMergeCooccurrences(t *testing.T) {
	type args struct {
		toCooccurrences   []Cooccurrence
		fromCooccurrences []Cooccurrence
	}
	tests := []struct {
		name string
		args args
		want []Cooccurrence
	}{
		{
			name: "",
			args: args{
				toCooccurrences: []Cooccurrence{
					{Word: "あ", Point: 1},
					{Word: "い", Point: 2},
					{Word: "う", Point: 3},
					{Word: "え", Point: 4},
					{Word: "お", Point: 5},
				},
				fromCooccurrences: []Cooccurrence{
					{Word: "う", Point: 6},
					{Word: "え", Point: 7},
					{Word: "お", Point: 8},
					{Word: "か", Point: 9},
					{Word: "き", Point: 10},
				},
			},
			want: []Cooccurrence{
				{Word: "あ", Point: 1},
				{Word: "い", Point: 2},
				{Word: "う", Point: 9},
				{Word: "え", Point: 11},
				{Word: "お", Point: 13},
				{Word: "か", Point: 9},
				{Word: "き", Point: 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeCooccurrences(tt.args.toCooccurrences, tt.args.fromCooccurrences); !reflect.DeepEqual(SortByWordAndPoint(got), SortByWordAndPoint(tt.want)) {
				t.Errorf("MergeCooccurrences() = %v, want %v", got, tt.want)
			}
		})

	}
}

// 結果表示用にソートする
func SortByWordAndPoint(coOccurrences []Cooccurrence) []Cooccurrence {
	sort.Slice(coOccurrences, func(i, j int) bool {
		// Word で昇順ソート
		if coOccurrences[i].Word < coOccurrences[j].Word {
			return true
		} else if coOccurrences[i].Word > coOccurrences[j].Word {
			return false
		}

		// Word が同じ場合は Point で降順ソート
		return coOccurrences[i].Point > coOccurrences[j].Point
	})
	return coOccurrences
}
