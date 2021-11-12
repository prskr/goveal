# Goveal

[![Actions Status](https://github.com/baez90/goveal/workflows/Go/badge.svg)](https://github.com/baez90/goveal/actions)

Goveal is very small and very simple tool that reads Markdown from a given file, renders it into a HTML template and
serves it as local HTTP server. Right now Goveal uses Reveal.js 4.1.2 to create presentations and therefore also
includes all features of Reveal.js 4.1.2.

## Install

The easiest and fastest way to install Goveal is to use a pre-built binary from the [releases](https://github.com/baez90/goveal/releases/latest).

If you have Go in the latest version installed you can also build your own version of Goveal:

```shell
task build
```

Requirements:

- [task](https://taskfile.dev/)
- _Optional:_ [goreleaser](https://goreleaser.com/) (for `task snapshot-release` to build all binaries)

_Note: All script tasks in the [Taskfile.yml](Taskfile.yml) are meant to be executed with Linux/MacOS. Binaries for Windows are provided but not tested!_ 

## Usage

```bash
goveal serve ./slides.md
```

| Param                    | Description                                                                  | Default value           |
| ------------------------ | ---------------------------------------------------------------------------- | ----------------------- |
| `--host`                 | Hostname the binary is listening on                                          | `localhost`             |
| `--port`                 | Port the binary is listening on                                              | `2233`                  |
| `--code-theme`           | highlight.js theme to use                                                    | `monokai`               |
| `--transition`           | Transition effect to show between slides                                     | `none`                  |
| `--navigationMode`       | Navigation mode to use when using the cursor keys to navigate through slides | `default`               |
| `--config`               | Path to the config file see [config](#config)                                | `$HOME/goveal:./goveal` |
| `--horizontal-separator` | horizontal separator to split slides                                         | `===`                   |
| `--vertical-separator`   | vertical separator to split slides                                           | `---`                   |
| `--theme`                | reveal.js theme to use                                                       | `white`                 |
| `-h` / `--help`          | shows help                                                                   |                         |

## Config

Goveal supports multiple configuration mechanisms. It tries to load a configuration file from `$HOME` or from `.`
i.e. `$HOME/goveal.yaml` or `$HOME/goveal.yml` or `./goveal.yaml` and so on.

Most options that can be set via commandline flags can also be set via configuration file (actually all but
the `--config` switch does not make sense in the configuration file, does it? :wink:). It is more a convenience feature
to be able to set a theme and so on and so forth for the presentation without having to pass it every time as parameter.

Furthermore goveal supports configuration hot reloading i.e. you can play around with different themes and the rendered
theme will be changed whenever you hit the save button!

See also an example at [`./examples/goveal.yaml`](./examples/goveal.yaml).

### Custom CSS

To add custom CSS as theme overrides use a config file and add the `stylesheets` property. It takes a list of relative (
mandatory!) paths to CSS files that are included automatacally after the page was loaded so that they really overload
everything added by Reveal and plugins.

the sample configuration file [`./examples/goveal.yaml`](./examples/goveal.yaml) also contains a sample how to add
custom CSS.
