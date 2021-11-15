// Copyright © 2019 Peter Kurfer peter.kurfer@googlemail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rendering

import (
	"path"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/imdario/mergo"
	"github.com/spf13/viper"

	"github.com/baez90/goveal/internal/encoding"
)

var defaultParams = RevealParams{
	Theme:                 "white",
	CodeTheme:             "vs",
	Transition:            "None",
	NavigationMode:        "default",
	HorizontalSeparator:   "===",
	VerticalSeparator:     "---",
	SlideNumberVisibility: "all",
	SlideNumberFormat:     "h.v",
	StyleSheets:           make([]string, 0),
	FilesToMonitor:        make([]string, 0),
}

type RevealParams struct {
	Theme                 string              `mapstructure:"theme"`
	CodeTheme             string              `mapstructure:"codeTheme"`
	Transition            string              `mapstructure:"transition"`
	NavigationMode        string              `mapstructure:"navigationMode"`
	HorizontalSeparator   string              `mapstructure:"horizontalSeparator"`
	VerticalSeparator     string              `mapstructure:"verticalSeparator"`
	SlideNumberVisibility string              `mapstructure:"slideNumberVisibility"`
	SlideNumberFormat     string              `mapstructure:"slideNumberFormat"`
	StyleSheets           []string            `mapstructure:"stylesheets"`
	FilesToMonitor        []string            `mapstructure:"filesToMonitor"`
	WorkingDirectory      string              `mapstructure:"working-dir"`
	LineEnding            encoding.LineEnding `mapstructure:"-"`
}

func (params *RevealParams) Load() error {
	_ = viper.Unmarshal(params)
	expandGlobs(params)
	return mergo.Merge(params, &defaultParams)
}

func expandGlobs(params *RevealParams) {
	var allFiles []string

	for _, f := range params.FilesToMonitor {
		var err error

		f, err = filepath.Abs(f)
		if err != nil {
			continue
		}

		var matches []string
		if matches, err = doublestar.Glob(f); err != nil {
			continue
		}

		for idx := range matches {
			if relative, err := filepath.Rel(params.WorkingDirectory, matches[idx]); err != nil {
				continue
			} else {
				matches[idx] = path.Join("/", relative)
			}
		}

		allFiles = append(allFiles, matches...)
	}
	params.FilesToMonitor = allFiles
}
