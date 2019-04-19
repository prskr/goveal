#!/usr/bin/env bash

mkdir -p ./assets/reveal
curl -sL "https://github.com/hakimel/reveal.js/archive/${1:-3.8.0}.tar.gz" | tar -xvz --strip-components=1 -C ./assets/reveal --wildcards *.js --wildcards *.css --exclude test --exclude gruntfile.js
