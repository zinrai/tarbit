# tarbit

A command-line tool for automatic tar archive handling. It automatically detects whether to extract or compress an archive based on the existence of files and directories.

I made this because I can't remember the options for compressing and decompressing with tar and always have to search for them.

## Features

- Automatic operation detection (extract/compress)
- Support for multiple compression formats (gzip, bzip2, xz)
- Automatic directory creation for extraction
- Path-aware operation

## Supported Formats

- `.tar.gz`, `.tgz` (gzip)
- `.tar.bz2`, `.tbz2` (bzip2)
- `.tar.xz` (xz)

## Requirements

`tar` command with gzip, bzip2, and xz support.

## Installation

build from source:

```bash
$ go build
```

## Usage

### Extract an archive

```bash
$ tarbit archive.tar.gz
```

1. Create a directory named `archive` in the current working directory
2. Extract the contents of `archive.tar.gz` into the `archive` directory

### Compress a directory

```bash
$ tarbit output.tar.gz
```

1. Check if `output` directory exists
2. Compress the directory into `output.tar.gz`

## License

This project is licensed under the [MIT License](./LICENSE).
