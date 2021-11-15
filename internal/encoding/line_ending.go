package encoding

import (
	"bufio"
	"io"
)

const (
	LineEndingUnknown LineEnding = ""
	LineEndingWindows LineEnding = "\r\n"
	LineEndingUnix    LineEnding = "\n"
)

type LineEnding string

func (e LineEnding) String() string {
	return string(e)
}

func (e LineEnding) Escaped() string {
	switch e {
	case LineEndingUnix:
		return "\\n"
	case LineEndingWindows:
		return "\\r\\n"
	default:
		return ""
	}
}

func Detect(reader io.Reader) (LineEnding, error) {
	bufferedReader := bufio.NewReader(reader)
	line, err := bufferedReader.ReadString(byte('\n'))
	if err != nil {
		return LineEndingUnknown, err
	}

	lineLength := len(line)
	if lineLength <= 1 || line[lineLength-2:] != LineEndingWindows.String() {
		return LineEndingUnix, nil
	}

	return LineEndingWindows, nil
}
