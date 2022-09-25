package rendering

import (
	"fmt"
	"html/template"
	"regexp"

	"github.com/gomarkdown/markdown/parser"
	"github.com/valyala/bytebufferpool"

	"code.icb4dc0.de/prskr/goveal/config"
)

const (
	parserExtensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
		parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs |
		parser.BackslashLineBreak | parser.DefinitionLists | parser.MathJax | parser.Titleblock |
		parser.OrderedListStart | parser.Attributes
)

func ToHTML(markdown string, renderCfg config.Rendering) (rendered []byte, err error) {
	var slides []rawSlide
	if slides, err = splitIntoRawSlides(markdown, renderCfg); err != nil {
		return nil, err
	}

	buf := bytebufferpool.Get()
	defer func() {
		buf.Reset()
		bytebufferpool.Put(buf)
	}()

	for idx := range slides {
		if rendered, err := slides[idx].ToHTML(); err != nil {
			return nil, err
		} else {
			_, _ = buf.WriteString(string(rendered))
		}
	}

	return buf.Bytes(), nil
}

type rawSlide struct {
	Content  string
	Children []rawSlide
}

func (s rawSlide) HasNotes() bool {
	return notesLineRegexp.MatchString(s.Content)
}

func (s rawSlide) ToHTML() (template.HTML, error) {
	if rendered, err := renderTemplate("slide.gohtml", s); err != nil {
		return "", err
	} else {
		//nolint:gosec // should not be sanitized
		return template.HTML(rendered), nil
	}
}

func splitIntoRawSlides(markdown string, renderCfg config.Rendering) ([]rawSlide, error) {
	var (
		verticalSplit, horizontalSplit *regexp.Regexp
		err                            error
	)
	if verticalSplit, err = regexp.Compile(fmt.Sprintf(splitFormat, renderCfg.VerticalSeparator)); err != nil {
		return nil, err
	}

	if horizontalSplit, err = regexp.Compile(fmt.Sprintf(splitFormat, renderCfg.HorizontalSeparator)); err != nil {
		return nil, err
	}

	horizontalSlides := horizontalSplit.Split(markdown, -1)
	slides := make([]rawSlide, 0, len(horizontalSlides))
	for _, hs := range horizontalSlides {
		s := rawSlide{
			Content: hs,
		}
		verticalSlides := verticalSplit.Split(hs, -1)
		s.Children = make([]rawSlide, 0, len(verticalSlides))
		for _, vs := range verticalSlides {
			s.Children = append(s.Children, rawSlide{Content: vs})
		}
		slides = append(slides, s)
	}

	return slides, nil
}
