package emoji

import (
	_ "embed"
	"encoding/json"
)

func isValidEmoji(input []byte) bool {
	_, exists := emojiMap[string(input)]
	return exists
}

var (
	//go:embed emoji.json
	emojiMapRaw []byte

	emojiMap map[string]string
)

func init() {
	rawMap := make(map[string][]string)
	if err := json.Unmarshal(emojiMapRaw, &rawMap); err != nil {
		panic(err)
	}

	emojiMap = make(map[string]string, len(rawMap))
	for emoji, keywords := range rawMap {
		for i := range keywords {
			emojiMap[keywords[i]] = emoji
		}
	}
}
