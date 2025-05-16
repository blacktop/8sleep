<p align="center">
  <a href="https://github.com/blacktop/clim8"><img alt="clim8 Logo" src="https://raw.githubusercontent.com/blacktop/clim8/main/docs/logo.png" height="400" /></a>
  <h4><p align="center">Control Eight Sleep via CLI</p></h4>
  <p align="center">
    <a href="https://github.com/blacktop/clim8/actions" alt="Actions">
          <img src="https://github.com/blacktop/clim8/actions/workflows/go.yml/badge.svg" /></a>
    <a href="https://github.com/blacktop/clim8/releases/latest" alt="Downloads">
          <img src="https://img.shields.io/github/downloads/blacktop/clim8/total.svg" /></a>
    <a href="https://github.com/blacktop/clim8/releases" alt="GitHub Release">
          <img src="https://img.shields.io/github/release/blacktop/clim8.svg" /></a>
    <a href="http://doge.mit-license.org" alt="LICENSE">
          <img src="https://img.shields.io/:license-mit-blue.svg" /></a>
</p>
<br>

## Why? 🤔

Oh, I think you know why 😏

## Getting Started

### Install

```bash
go install github.com/blacktop/clim8@latest
```

Or download the latest [release](https://github.com/blacktop/clim8/releases/latest)

Or via [homebrew](https://brew.sh)

```bash
brew install blacktop/tap/clim8
```

### Running

```bash
❱ clim8 --help
```
```bash
Eight Sleep CLI

Usage:
  clim8 [command]

Available Commands:
  help        Help about any command
  off         Turn off the Eight Sleep Pod
  on          Turn on the Eight Sleep Pod
  temp        Set the temperature of the Eight Sleep Pod

Flags:
  -e, --email string      Email address
  -h, --help              help for clim8
  -p, --password string   Password
  -V, --verbose           Enable verbose debug logging

Use "clim8 [command] --help" for more information about a command.
```

![demo](vhs.gif)

## License

MIT Copyright (c) 2025 **blacktop**