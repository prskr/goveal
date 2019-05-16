#!/usr/bin/env bash

mkdir -p ./assets/reveal
curl -sL "https://github.com/hakimel/reveal.js/archive/${1:-3.8.0}.tar.gz" | tar -xvz --strip-components=1 -C ./assets/reveal --wildcards *.js --wildcards *.css --wildcards *.html --exclude test --exclude gruntfile.js
mkdir -p ./assets/reveal/plugin/
git clone https://github.com/denehyg/reveal.js-menu.git ./assets/reveal/plugin/menu

rm -f ./assets/reveal/plugin/menu/bower.json
rm -f ./assets/reveal/plugin/menu/CONTRIBUTING.md
rm -f ./assets/reveal/plugin/menu/LICENSE
rm -f ./assets/reveal/plugin/menu/package.json
rm -f ./assets/reveal/plugin/menu/README.md
