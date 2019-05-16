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

package config

import (
	"github.com/spf13/viper"
)

type RevealParams struct {
	Theme               string
	CodeTheme           string
	Transition          string
	NavigationMode      string
	HorizontalSeparator string
	VerticalSeparator   string
	StyleSheets         []string
}

func (params *RevealParams) Load() {
	params.Theme = viper.GetString("theme")
	params.CodeTheme = viper.GetString("code-theme")
	params.Transition = viper.GetString("transition")
	params.NavigationMode = viper.GetString("navigationMode")
	params.HorizontalSeparator = viper.GetString("horizontal-separator")
	params.VerticalSeparator = viper.GetString("vertical-separator")
	params.StyleSheets = viper.GetStringSlice("stylesheets")
}
