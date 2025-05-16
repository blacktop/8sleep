<p align="center">
  <a href="https://github.com/blacktop/8sleep"><img alt="8sleep Logo" src="https://raw.githubusercontent.com/blacktop/8sleep/main/docs/logo.webp" /></a>
  <h1 align="center">8sleep</h1>
  <h4><p align="center">Control 8sleep via CLI</p></h4>
  <p align="center">
    <a href="https://github.com/blacktop/8sleep/actions" alt="Actions">
          <img src="https://github.com/blacktop/8sleep/actions/workflows/go.yml/badge.svg" /></a>
    <a href="https://github.com/blacktop/8sleep/releases/latest" alt="Downloads">
          <img src="https://img.shields.io/github/downloads/blacktop/8sleep/total.svg" /></a>
    <a href="https://github.com/blacktop/8sleep/releases" alt="GitHub Release">
          <img src="https://img.shields.io/github/release/blacktop/8sleep.svg" /></a>
    <a href="http://doge.mit-license.org" alt="LICENSE">
          <img src="https://img.shields.io/:license-mit-blue.svg" /></a>
</p>
<br>

## Why? 🤔

Oh, I think you know why 😏

## Getting Started

### Install

```bash
go install github.com/blacktop/8sleep@latest
```

Or download the latest [release](https://github.com/blacktop/8sleep/releases/latest)

Or via [homebrew](https://brew.sh)

```bash
brew install blacktop/tap/8sleep
```

### Running

```bash
8sleep CLI

Usage:
  8sleep [command]

Available Commands:
  help        Help about any command
  off         Turn off the 8Sleep Pod
  on          Turn on the 8Sleep Pod
  temp        Set the temperature of the 8Sleep Pod

Flags:
  -e, --email string      Email address
  -h, --help              help for 8sleep
  -p, --password string   Password
  -V, --verbose           Enable verbose debug logging

Use "8sleep [command] --help" for more information about a command.
```

![demo](vhs.gif)

## License

MIT Copyright (c) 2025 **blacktop**