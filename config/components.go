package config

const (
	defaultWidth  uint = 960
	defaultHeight uint = 700
)

var defaults = map[string]any{
	"mermaid.theme":                       "forest",
	"theme":                               "beige",
	"width":                               defaultWidth,
	"height":                              defaultHeight,
	"codeTheme":                           "monokai",
	"verticalSeparator":                   `\*\*\*`,
	"horizontalSeparator":                 `---`,
	"transition":                          TransitionNone,
	"controlsLayout":                      ControlsLayoutEdges,
	"controls":                            true,
	"progress":                            true,
	"history":                             true,
	"center":                              true,
	"slideNumber":                         true,
	"menu.numbers":                        true,
	"menu.useTextContentForMissingTitles": true,
	"menu.transitions":                    true,
	"menu.hideMissingTitles":              true,
	"menu.markers":                        true,
	"menu.openButton":                     true,
}

const (
	TransitionNone    Transition = "none"
	TransitionFade    Transition = "fade"
	TransitionSlide   Transition = "slide"
	TransitionConvex  Transition = "convex"
	TransitionConcave Transition = "concave"
	TransitionZoom    Transition = "zoom"

	ControlsLayoutBottomRight ControlsLayout = "bottom-right"
	ControlsLayoutEdges       ControlsLayout = "edges"
)

type (
	Transition     string
	ControlsLayout string
	Mermaid        struct {
		Theme string `json:"theme"`
	}
	Rendering struct {
		VerticalSeparator   string
		HorizontalSeparator string
		Stylesheets         []string
	}
	Reveal struct {
		Theme          string         `json:"theme"`
		CodeTheme      string         `json:"codeTheme"`
		Transition     Transition     `json:"transition"`
		Controls       bool           `json:"controls"`
		ControlsLayout ControlsLayout `json:"controlsLayout"`
		Progress       bool           `json:"progress"`
		History        bool           `json:"history"`
		Center         bool           `json:"center"`
		SlideNumber    bool           `json:"slideNumber"`
		Width          uint           `json:"width"`
		Height         uint           `json:"height"`
		Menu           struct {
			Numbers                        bool `json:"numbers"`
			UseTextContentForMissingTitles bool `json:"useTextContentForMissingTitles"`
			Transitions                    bool `json:"transitions"`
			HideMissingTitles              bool `json:"hideMissingTitles"`
			Markers                        bool `json:"markers"`
			OpenButton                     bool `json:"openButton"`
		} `json:"menu"`
	}
	Components struct {
		ConfigFileInUse string    `mapstructure:"-"`
		Reveal          Reveal    `mapstructure:",squash"`
		Rendering       Rendering `mapstructure:",squash"`
		Mermaid         Mermaid
	}
)

func (t Transition) String() string {
	return string(t)
}
