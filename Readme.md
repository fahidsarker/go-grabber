# Go Grabber

Go Grabber is a command-line tool for downloading files from a specified URL. It parses the HTML content of the URL, extracts downloadable file links, and downloads them to a specified directory.

## Features

- Concurrent downloads with configurable workers.
- Supports a variety of file extensions.
- Easy to use with minimal flags and arguments.

## Installation

Download the pre-built binary from the [GitHub Releases](https://github.com/fahidsarker/go-grabber/releases) page.

1. Visit the [releases page](https://github.com/fahidsarker/go-grabber/releases).
2. Download the appropriate binary for your operating system.
3. Extract the binary to a location included in your system's PATH.

## Usage

Once you have the binary, you can use the tool as follows:

```bash
./grabber [--workers=N] <URL> <output_dir>