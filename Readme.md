# Go Grabber

Go Grabber is a command-line tool for downloading files from URLs. It can parse HTML content from web pages to extract downloadable file links, export those URLs to files, and download them with concurrent workers for better performance.

## Features

- **Multiple operation modes**: Download directly from URLs or from exported URL files
- **Export functionality**: Save discovered downloadable URLs to a file for later use
- **Concurrent downloads**: Configurable number of worker threads for parallel downloads (default: 4)
- **Google Drive support**: Automatically detects and converts Google Drive sharing URLs to direct download links
- **Universal file support**: Automatically detects and downloads files with any extension
- **Directory organization**: Support for subdirectory structure in URL files using comments
- **Debug mode**: Save HTML content to debug.html for troubleshooting
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
./grabber dl -from-url=<URL> -o=<output_directory> [-workers=N] [-allow-g-drive] [-d|-debug]
```

### Download Files from Exported File
Download files from a previously exported URL list:
```bash
./grabber dl -from-file=<url_file> -o=<output_directory> [-workers=N] [-allow-g-drive]
```

### Download Files from HTML File
Download files from a local HTML file:
```bash
./grabber dl -from-html=<html_file> -o=<output_directory> [-workers=N] [-allow-g-drive]
```

### Export URLs to File
Extract and save downloadable URLs from a webpage:
```bash
./grabber export -from-url=<URL> -o=<output_file> [-allow-g-drive] [-d|-debug]
```

### Export URLs from Multiple Sources
Extract URLs from multiple webpages listed in a file:
```bash
./grabber export -from-file=<urls_file> -o=<output_file> [-allow-g-drive] [-d|-debug]
```

### Export URLs from HTML File
Extract URLs from a local HTML file:
```bash
./grabber export -from-html=<html_file> -o=<output_file> [-allow-g-drive]
```

### Examples

**Download files from a webpage:**
```bash
./grabber dl -from-url=https://example.com/files -o=./downloads -workers=8
```

**Download files with Google Drive support:**
```bash
./grabber dl -from-url=https://example.com/files -o=./downloads -allow-g-drive
```

**Export URLs for later use:**
```bash
./grabber export -from-url=https://example.com/files -o=urls.txt
```

**Export URLs with debug mode enabled:**
```bash
./grabber export -from-url=https://example.com/files -o=urls.txt -debug
```

**Download from exported URL file:**
```bash
./grabber dl -from-file=urls.txt -o=./downloads -workers=4
```

**Process multiple URLs from a file:**
```bash
./grabber export -from-file=multiple_urls.txt -o=all_downloadable_urls.txt
```

## Command Reference

### Commands
- `dl` - Download files from URL or from exported file
- `export` - Export downloadable URLs to a file

### Flags
- `-workers=N` - Number of concurrent workers (default: 4)
- `-from-url=URL` - Source URL to parse for downloadable files
- `-from-file=FILE` - File containing URLs (for dl: exported file, for export: list of URLs)
- `-from-html=FILE` - HTML file to parse for downloadable files
- `-o=PATH` - Output directory for downloads or output file for export
- `-d, -debug` - Enable debug mode (saves HTML content to debug.html)
- `-allow-g-drive` - Allow detection and conversion of Google Drive URLs

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

## Google Drive Support

Go Grabber can detect and convert Google Drive sharing URLs to direct download links when the `-allow-g-drive` flag is used.

### Supported Google Drive URL Formats

- `https://drive.google.com/file/d/FILE_ID/view?usp=sharing`
- `https://drive.google.com/open?id=FILE_ID`
- `https://docs.google.com/document/d/FILE_ID/edit`
- `https://docs.google.com/spreadsheets/d/FILE_ID/edit`
- `https://docs.google.com/presentation/d/FILE_ID/edit`

### Example

```bash
# Enable Google Drive URL conversion
./grabber dl -from-url=https://example.com/page-with-gdrive-links -o=./downloads -allow-g-drive
```

When a Google Drive URL is detected, it will be automatically converted to `https://drive.google.com/uc?export=download&id=FILE_ID` for direct downloading.

## Debug Mode

Use the `-debug` or `-d` flag to save the HTML content of parsed web pages to `debug.html`. This is useful for troubleshooting when expected files are not being detected.

```bash
./grabber export -from-url=https://example.com/files -o=urls.txt -debug
```

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