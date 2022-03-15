package emoji

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gomarkdown/markdown/ast"
)

func NewEmojiParser() *EmojiParser {
	return &EmojiParser{
		seen: make(map[string]bool),
	}
}

// Node is a node containing an emoji
type Node struct {
	ast.Leaf
}

type EmojiParser struct {
	seen map[string]bool
}

func (p *EmojiParser) EmojiParser(data []byte) (parsedNode ast.Node, result []byte, newLength int) {
	if p.seen[string(data)] {
		// Already processed
		return nil, nil, 0
	}
	if bytes.Contains(data, []byte("class=\"emoji\"")) {
		// Already processed
		return nil, nil, 0
	}
	dataLen := len(data)
	if dataLen <= 1 {
		// Not long enough to be an emoji
		return nil, nil, 0
	}
	if bytes.IndexByte(data, ':') == -1 {
		// No emoji delimiters
		return nil, nil, 0
	}
	// Translate emojis to HTML
	resData := make([]byte, 0)
	startIndex := bytes.IndexByte(data, ':')
	resData = append(resData, data[0:startIndex]...)
	for {
		if startIndex >= len(data) {
			// Done
			break
		}
		endIndex := bytes.IndexByte(data[startIndex+1:], ':') + startIndex + 1
		if endIndex > startIndex {
			name := string(data[startIndex+1 : endIndex])
			if isValidEmoji([]byte(name)) {
				startIndex = endIndex + 1
				url := fmt.Sprintf(`<img class="emoji" src=%q alt=":%s:"></img>`, GenerateEmojiURL(name), name)
				resData = append(resData, []byte(url)...)
			} else {
				resData = append(resData, data[startIndex:endIndex]...)
				startIndex = endIndex
			}
			if startIndex == dataLen {
				break
			}
		} else {
			break
		}
	}
	if startIndex < dataLen {
		resData = append(resData, data[startIndex:]...)
	}

	if !bytes.Contains(resData, []byte("class=\"emoji\"")) {
		// Processed with no changes
		p.seen[string(resData)] = true
	}

	return &ast.Softbreak{}, resData, dataLen
}

func GenerateEmojiURL(emoji string) string {
	code, exists := emojiMap[emoji]
	if !exists {
		return ""
	}
	res := ""
	chars := utf8.RuneCountInString(code)
	curChar := 1
	for _, c := range code {
		tmp := strings.Trim(strings.ToLower(strconv.QuoteRuneToASCII(c)), "'")
		if tmp != `\ufe0f` || (chars > 2 && curChar == chars) {
			// Valid character to add
			if curChar > 1 {
				res += "-"
			}
			if len(tmp) == 1 {
				res += fmt.Sprintf("%x", []byte(tmp))
			} else {
				res += strings.TrimLeft(tmp[2:], "0")
			}
		}
		curChar++
	}
	return fmt.Sprintf("https://twemoji.maxcdn.com/2/svg/%s.svg", res)
}
