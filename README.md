# Goveal

Goveal is very small an very simple tool that reads Markdown from a given file, renders  it into a HTML template and serves it as local HTTP server.
Right now Goveal uses Reveal.js 3.8.0 to create presentations and therefore also includes all features of Reveal.js 3.8.0.

## Usage

The easiest way to use `goveal` is to download a release, and run it from your `$PATH`.

```bash
goveal serve ./slides.md
```

| Param                    | Description                                   | Default value           |
| ------------------------ | --------------------------------------------- | ----------------------- |
| `--host`                 | Hostname the binary is listening on           | `localhost`             |
| `--port`                 | Port the binary is listening on               | `2233`                  |
| `--code-theme`           | highlight.js theme to use                     | `monokai`               |
| `--config`               | Path to the config file see [config](#config) | `$HOME/goveal:./goveal` |
| `-h` / `--help`          | shows help                                    |                         |
| `--horizontal-separator` | horizontal separator to split slides          | `===`                   |
| `--vertical-separator`   | vertical separator to split slides            | `---`                   |
| `--theme`                | reveal.js theme to use                        | `white`                 |

## Config

Goveal supports multiple configuration mechanisms.
It tries to load a configuration file from `$HOME` or from `.` i.e. `$HOME/goveal.yaml` or `$HOME/goveal.yml` or `./goveal.yaml` and so on.

Most options that can be set via commandline flags can also be set via configuration file (actually all but the `--config` switch does not make sense in the configuration file, does it? :wink:).
It is more a convenience feature to be able to set a theme and so on and so forth for the presentation without having to pass it every time as parameter.

Furthermore goveal supports configuration hot reloading i.e. you can play around with different themes and the rendered theme will be changed whenever you hit the save button!

See also an example at [`./examples/goveal.yaml`](./examples/goveal.yaml).