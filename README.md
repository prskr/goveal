# Goveal

[![Actions Status](https://github.com/baez90/goveal/workflows/Go/badge.svg)](https://github.com/baez90/goveal/actions)

Goveal is very small and very simple tool that reads Markdown from a given file, renders it into an HTML template and serves it as local HTTP server.
Right now Goveal uses Reveal.js 4.2.1 to create presentations and therefore also includes all features of Reveal.js 4.2.1.

In contrary to Reveal.js `goveal` ships with its own Markdown parser and renderer which is why some features are working slightly different from Reveal.js.
See [Markdown](#markdown) for further details.

Besides all features from Reveal.js `goveal` comes with first class support for [mermaid-js](https://mermaid-js.github.io/).
Just inline your diagrams as code and enjoy!

## Install

The easiest and fastest way to install Goveal is to use a pre-built binary from the [releases](https://github.com/baez90/goveal/releases/latest).

There's also a pre-built container image available you can use if you don't want to download the binary.

### Build from source

If you have Go in the latest version installed you can also build your own version of Goveal:

```shell
task build
```

Requirements:

- [task](https://taskfile.dev/)
- _Optional:_ [goreleaser](https://goreleaser.com/) (for `task snapshot-release` to build all binaries)

_Note: All script tasks in the [Taskfile.yml](Taskfile.yml) are meant to be executed with Linux/MacOS. Binaries for
Windows are provided but not tested!_

## Usage

### Local installation

```bash
goveal serve ./slides.md
```

| Param            | Description                                   | Default value           |
| ---------------- | --------------------------------------------- | ----------------------- |
| `--host`         | Hostname the binary is listening on           | `127.0.0.1`             |
| `--port`         | Port the binary is listening on               | `2233`                  |
| `--config`       | Path to the config file see [config](#config) | `$HOME/goveal:./goveal` |
| `--open-browser` | Open a browser after starting the web server  | `true`                  |
| `-h` / `--help`  | shows help                                    |                         |

### Container

Assuming your slides are in a file called `slides.md` in the current directory, you can start the presentation like
this:

```shell
podman/docker run --rm -ti -v `pwd`:/work -w /work -p 2233:2233 ghcr.io/baez90/goveal:0 serve --host 0.0.0.0 slides.md
```

By default `goveal` only listens on `127.0.0.1`. To allow traffic from outside of the container you've to change the
binding to `0.0.0.0`.

## Config

Goveal supports multiple configuration mechanisms. It tries to load a configuration file from `$HOME` or from `.`
i.e. `$HOME/goveal.yaml` or `$HOME/goveal.yml` or `./goveal.yaml` and so on.

The config allows setting a lot of options exposed by Reveal.js.
There are still a few missing, and I won't guarantee to support all options in the future.

Furthermore, `goveal` supports configuration hot reloading i.e. you can play around with different themes and the rendered
theme will be changed whenever you hit the save button!

See also an example at [`./examples/goveal.yaml`](./examples/goveal.yaml).
I try to keep the example up to date to cover **all** supported config options as kind of documentation.

## Markdown

A good point to start is the [slides.md](examples/slides.md) in the `examples` directory.
I try to showcase everyting possible with `goveal` in this presentation.

The most remarkable difference between `goveal` and the `marked` driven Reveal.js markdown is how line numbers in listings are working:

Reveal.js:

<pre lang="no-highlight"><code>```cs [1-2|3|4]
var i = 10;
for (var j = 0; j < i; j++) {
    Console.WriteLine($"{j}");
}
```</code></pre>

`goveal`:

<pre lang="no-highlight"><code>
{line-numbers="1-2|3|4"}
```cs
var i = 10;
for (var j = 0; j < i; j++) {
    Console.WriteLine($"{j}");
}
```</code></pre>

This is because the Markdown parser used in `goveal` currently does not support additional attributes for code blocks.

## Custom CSS

To add custom CSS as theme overrides use a config file and add the `stylesheets` property. It takes a list of relative (
mandatory!) paths to CSS files that are included automatacally after the page was loaded so that they really overload
everything added by Reveal and plugins.

Changes in the custom CSS files are monitored and propagated via SSE to the presentation immediately. 
No page reload necessary!

The sample configuration file [`./examples/goveal.yaml`](./examples/goveal.yaml) also contains a sample how to add
custom CSS.