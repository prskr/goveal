package config

var (
	defaults = map[string]interface{}{
		"mermaid.theme":                       "forest",
		"theme":                               "beige",
		"codeTheme":                           "monokai",
		"transition":                          TransitionNone,
		"controlsLayout":                      ControlsLayoutEdges,
		"controls":                            true,
		"progress":                            true,
		"history":                             true,
		"center":                              true,
		"slideNumber":                         true,
		"menu.numbers":                        true,
		"menu.useTextContentForMissingTitles": true,
	}
)

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
		Menu           struct {
			Numbers                        bool `json:"numbers"`
			UseTextContentForMissingTitles bool `json:"useTextContentForMissingTitles"`
			Transitions                    bool
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
	switch t {
	case TransitionNone:
		return "none"
	default:
		return string(t)
	}
}
