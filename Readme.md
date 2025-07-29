# Go Grabber

Go Grabber is a command-line tool for downloading files from URLs. It can parse HTML content from web pages to extract downloadable file links, export those URLs to files, and download them with concurrent workers for better performance.

## Features

- **Multiple operation modes**: Download directly from URLs or from exported URL files
- **Export functionality**: Save discovered downloadable URLs to a file for later use
- **Concurrent downloads**: Configurable number of worker threads for parallel downloads (default: 4)
- **Universal file support**: Automatically detects and downloads files with any extension
- **Directory organization**: Support for subdirectory structure in URL files using comments
- **Robust error handling**: Clear error messages and graceful failure handling
- **Cross-platform**: Works on Windows, macOS, and Linux
- **Dependency minimal**: Only requires `golang.org/x/net` for HTML parsing

## Installation

### Option 1: Download Pre-built Binary
Download the pre-built binary from the [GitHub Releases](https://github.com/fahidsarker/go-grabber/releases) page.

1. Visit the [releases page](https://github.com/fahidsarker/go-grabber/releases).
2. Download the appropriate binary for your operating system.
3. Extract the binary to a location included in your system's PATH.

### Option 2: Build from Source

**Prerequisites:**
- Go 1.20 or higher

**Steps:**
```bash
git clone https://github.com/fahidsarker/go-grabber.git
cd go-grabber
go mod tidy
go build -o grabber main.go
```

**For system-wide installation:**
```bash
# Build and install to $GOPATH/bin
go install

# Or build and move to /usr/local/bin (requires sudo)
go build -o grabber main.go
sudo mv grabber /usr/local/bin/
```

## Usage

Go Grabber supports two main commands: `dl` (download) and `export`.

### Download Files from URL
Download all files found on a webpage:
```bash
./grabber dl -from-url=<URL> -o=<output_directory> [-workers=N]
```

### Download Files from Exported File
Download files from a previously exported URL list:
```bash
./grabber dl -from-file=<url_file> -o=<output_directory> [-workers=N]
```

### Export URLs to File
Extract and save downloadable URLs from a webpage:
```bash
./grabber export -from-url=<URL> -o=<output_file>
```

### Examples

**Download files from a webpage:**
```bash
./grabber dl -from-url=https://example.com/files -o=./downloads -workers=8
```

**Export URLs for later use:**
```bash
./grabber export -from-url=https://example.com/files -o=urls.txt
```

**Download from exported URL file:**
```bash
./grabber dl -from-file=urls.txt -o=./downloads -workers=4
```

## Command Reference

### Commands
- `dl` - Download files from URL or from exported file
- `export` - Export downloadable URLs to a file

### Flags
- `-workers=N` - Number of concurrent workers (default: 4)
- `-from-url=URL` - Source URL to parse for downloadable files
- `-from-file=FILE` - File containing URLs to download (one per line)
- `-o=PATH` - Output directory for downloads or output file for export

## URL File Format

When using `-from-file`, the file should contain one URL per line:
```
https://example.com/file1.pdf
https://example.com/file2.zip
# This is a comment and will be ignored
https://example.com/file3.mp4
```

Lines starting with `#` are treated as comments and ignored.

### Advanced: Directory Organization

You can organize downloads into subdirectories by using comments to specify directory names:

```
# documents
https://example.com/manual.pdf
https://example.com/guide.pdf

# videos
https://example.com/intro.mp4
https://example.com/tutorial.mp4

# archives
https://example.com/source.zip
https://example.com/backup.tar.gz
```

Files listed after a `# directory_name` comment will be downloaded to `output_dir/directory_name/`.

## Requirements

- Go 1.20 or higher (for building from source)
- Internet connection for downloading files
- Write permissions for the output directory

## Technical Details

- **HTML Parsing**: Uses `golang.org/x/net/html` for robust HTML parsing
- **Concurrent Downloads**: Worker pool pattern for efficient parallel downloads
- **File Detection**: Automatically detects downloadable files by checking for file extensions
- **URL Resolution**: Properly resolves relative URLs against the base URL
- **Error Recovery**: Individual download failures don't stop the entire process

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source. Please check the repository for license details.

## Troubleshooting

### Common Issues

**"No downloadable files found"**
- The webpage might not contain direct file links
- Files might be served through JavaScript or require authentication
- Check if the URL is accessible and returns HTML content

**"Permission denied" errors**
- Ensure you have write permissions to the output directory
- Try running with appropriate permissions or choose a different output directory

**Download failures**
- Some servers may block automated requests
- Files might be behind authentication or have access restrictions
- Network connectivity issues

### Getting Help

If you encounter issues:
1. Check that the URL is accessible in your browser
2. Verify you have write permissions to the output directory
3. Try reducing the number of workers (`-workers=1`) for debugging
4. Open an issue on GitHub with details about the error