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
	"github.com/bmatcuk/doublestar/v2"
	"github.com/spf13/viper"
)

type RevealParams struct {
	Theme                 string   `mapstructure:"theme"`
	CodeTheme             string   `mapstructure:"code-theme"`
	Transition            string   `mapstructure:"transition"`
	NavigationMode        string   `mapstructure:"navigationMode"`
	HorizontalSeparator   string   `mapstructure:"horizontal-separator"`
	VerticalSeparator     string   `mapstructure:"vertical-separator"`
	SlideNumberVisibility string   `mapstructure:"slide-number-visibility"`
	SlideNumberFormat     string   `mapstructure:"slide-number-format"`
	StyleSheets           []string `mapstructure:"stylesheets"`
	FilesToMonitor        []string `mapstructure:"filesToMonitor"`
}

func (params *RevealParams) Load() {
	_ = viper.Unmarshal(params)
	expandGlobs(params)
}

func expandGlobs(params *RevealParams) {
	var allFiles []string

	for _, f := range params.FilesToMonitor {
		var err error

		var matches []string
		if matches, err = doublestar.Glob(f); err != nil {
			continue
		}
		allFiles = append(allFiles, matches...)
	}
	params.FilesToMonitor = allFiles
}
