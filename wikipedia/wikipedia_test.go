package wikipedia

import "testing"

func TestFetchWikipediaContent(t *testing.T) {
	testCases := []struct {
		name     string
		wikiResp WikipediaResponse
		expected string
		err      error
	}{
		{
			name: "basic test",
			wikiResp: WikipediaResponse{
				Batchcomplete: false,
				Warnings: struct {
					Main struct {
						Warnings string "json:\"warnings\""
					} "json:\"main\""
					Revisions struct {
						Warnings string "json:\"warnings\""
					} "json:\"revisions\""
				}{},
				Query: struct {
					Pages []struct {
						Pageid    int    "json:\"pageid\""
						Ns        int    "json:\"ns\""
						Title     string "json:\"title\""
						Revisions []struct {
							Contentformat string "json:\"contentformat\""
							Contentmodel  string "json:\"contentmodel\""
							Content       string "json:\"content\""
						} "json:\"revisions\""
					} "json:\"pages\""
				}{Pages: []struct {
					Pageid    int    "json:\"pageid\""
					Ns        int    "json:\"ns\""
					Title     string "json:\"title\""
					Revisions []struct {
						Contentformat string "json:\"contentformat\""
						Contentmodel  string "json:\"contentmodel\""
						Content       string "json:\"content\""
					} "json:\"revisions\""
				}{{
					Pageid: 0,
					Ns:     0,
					Title:  "",
					Revisions: []struct {
						Contentformat string "json:\"contentformat\""
						Contentmodel  string "json:\"contentmodel\""
						Content       string "json:\"content\""
					}{
						{
							Content: "Test content", // テストデータ
						},
					},
				}}},
			},
			expected: "Test content",
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FetchWikipediaContent(tc.wikiResp)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, got: %v", tc.err, err)
			}
			if result != tc.expected {
				t.Errorf("Expected result: %v, got: %v", tc.expected, result)
			}
		})
	}
}
