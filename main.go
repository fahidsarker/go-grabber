package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "dl":
		handleDownload()
	case "export":
		handleExport()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Go Grabber - File download utility")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  grabber dl [-workers=N] [-d|-debug] -from-url=<URL> -o=<output_dir>")
	fmt.Println("  grabber dl [-workers=N] [-d|-debug] -from-file=<exported_file> -o=<output_dir>")
	fmt.Println("  grabber export [-d|-debug] -from-url=<URL> -o=<output_file>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  dl      Download files from URL or from exported file")
	fmt.Println("  export  Export downloadable URLs to a file")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -workers=N      Number of concurrent workers (default: 4)")
	fmt.Println("  -from-url=URL   Source URL to parse for downloadable files")
	fmt.Println("  -from-file=FILE File containing URLs to download")
	fmt.Println("  -o=PATH         Output directory for downloads or output file for export")
	fmt.Println("  -d, -debug      Enable debug mode (saves HTML content to debug.html)")
}

func handleDownload() {
	fs := flag.NewFlagSet("dl", flag.ExitOnError)
	workers := fs.Int("workers", 4, "Number of concurrent download workers")
	fromURL := fs.String("from-url", "", "URL to download files from")
	fromFile := fs.String("from-file", "", "File containing URLs to download")
	output := fs.String("o", "", "Output directory")
	debug := fs.Bool("debug", false, "Enable debug mode (saves HTML to debug.html)")
	debugShort := fs.Bool("d", false, "Enable debug mode (saves HTML to debug.html)")

	fs.Parse(os.Args[2:])

	if *output == "" {
		fmt.Println("Error: -o flag is required")
		return
	}

	if *fromURL == "" && *fromFile == "" {
		fmt.Println("Error: Either -from-url or -from-file must be specified")
		return
	}

	if *fromURL != "" && *fromFile != "" {
		fmt.Println("Error: Cannot specify both -from-url and -from-file")
		return
	}

	var links []LinkWithDir
	var err error

	if *fromURL != "" {
		debugMode := *debug || *debugShort
		urlLinks, err := extractLinksFromURL(*fromURL, debugMode)
		if err != nil {
			fmt.Printf("Error extracting links from URL: %v\n", err)
			return
		}
		// Convert plain URLs to LinkWithDir with empty subdir
		links = make([]LinkWithDir, len(urlLinks))
		for i, link := range urlLinks {
			links[i] = LinkWithDir{URL: link, SubDir: ""}
		}
	} else {
		links, err = readLinksFromFile(*fromFile)
		if err != nil {
			fmt.Printf("Error reading links from file: %v\n", err)
			return
		}
	}

	if len(links) == 0 {
		fmt.Println("No downloadable files found")
		return
	}

	fmt.Printf("Found %d downloadable files\n", len(links))

	// Prepare output directory
	err = os.MkdirAll(*output, 0755)
	if err != nil {
		fmt.Printf("Error creating output dir: %v\n", err)
		return
	}

	downloadFiles(links, *output, *workers)
}

func handleExport() {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	fromURL := fs.String("from-url", "", "URL to extract downloadable files from")
	output := fs.String("o", "", "Output file to save URLs")
	debug := fs.Bool("debug", false, "Enable debug mode (saves HTML to debug.html)")
	debugShort := fs.Bool("d", false, "Enable debug mode (saves HTML to debug.html)")

	fs.Parse(os.Args[2:])

	if *fromURL == "" {
		fmt.Println("Error: -from-url flag is required")
		return
	}

	if *output == "" {
		fmt.Println("Error: -o flag is required")
		return
	}

	links, err := extractLinksFromURL(*fromURL, *debug || *debugShort)
	if err != nil {
		fmt.Printf("Error extracting links from URL: %v\n", err)
		return
	}

	if len(links) == 0 {
		fmt.Println("No downloadable files found")
		return
	}

	fmt.Printf("Found %d downloadable files\n", len(links))

	err = writeLinksToFile(links, *output)
	if err != nil {
		fmt.Printf("Error writing links to file: %v\n", err)
		return
	}

	fmt.Printf("URLs exported to: %s\n", *output)
}

func extractLinksFromURL(startURL string, debug bool) ([]string, error) {
	// Fetch and parse HTML
	resp, err := http.Get(startURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Save HTML to debug file if debug mode is enabled
	if debug {
		err := saveHTMLToDebugFile(body)
		if err != nil {
			fmt.Printf("Warning: Could not save debug HTML file: %v\n", err)
		} else {
			fmt.Println("HTML content saved to debug.html")
		}
	}

	// Parse HTML from the body
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	baseURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	return extractFileLinks(doc, baseURL), nil
}

func saveHTMLToDebugFile(htmlContent []byte) error {
	file, err := os.Create("debug.html")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(htmlContent)
	return err
}

// LinkWithDir represents a URL with its target subdirectory
type LinkWithDir struct {
	URL    string
	SubDir string
}

func readLinksFromFile(filename string) ([]LinkWithDir, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []LinkWithDir
	var currentSubDir string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			// Extract subdirectory from comment line
			// Remove leading '#' and trim whitespace
			subDir := strings.TrimSpace(line[1:])
			// Remove leading './' if present
			subDir = strings.TrimPrefix(subDir, "./")
			currentSubDir = subDir
		} else {
			// It's a URL, add it with current subdirectory
			links = append(links, LinkWithDir{
				URL:    line,
				SubDir: currentSubDir,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func writeLinksToFile(links []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, link := range links {
		_, err := writer.WriteString(link + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFiles(links []LinkWithDir, outputDir string, workers int) {
	// Start concurrent downloads
	var wg sync.WaitGroup
	jobs := make(chan LinkWithDir)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for linkInfo := range jobs {
				fmt.Printf("[Worker %d] Downloading: %s\n", workerID, linkInfo.URL)
				err := downloadFileToSubDir(linkInfo.URL, outputDir, linkInfo.SubDir)
				if err != nil {
					fmt.Printf("[Worker %d] Error: %v\n", workerID, err)
				}
			}
		}(i + 1)
	}

	for _, link := range links {
		jobs <- link
	}
	close(jobs)
	wg.Wait()
}

func extractFileLinks(n *html.Node, base *url.URL) []string {
	var links []string
	totalHref := 0

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if u, err := base.Parse(href); err == nil {
						totalHref++
						if hasAllowedExt(u.String()) {
							links = append(links, u.String())
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	fmt.Printf("Total href attributes found: %d\n", totalHref)

	return links
}

func hasAllowedExt(fullURL string) bool {
	// Define allowed file extensions
	allowedExts := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true, "bmp": true,
		"pdf": true, "doc": true, "docx": true, "txt": true, "rtf": true,
		"mp4": true, "avi": true, "mkv": true, "mov": true, "wmv": true, "flv": true, "webm": true,
		"mp3": true, "wav": true, "flac": true, "aac": true, "ogg": true,
		"zip": true, "rar": true, "7z": true, "tar": true, "gz": true, "bz2": true,
		"exe": true, "msi": true, "dmg": true, "pkg": true, "deb": true, "rpm": true,
		"iso": true, "img": true, "bin": true,
		"apk": true, "ipa": true,
		"xls": true, "xlsx": true, "ppt": true, "pptx": true,
		"css": true, "js": true, "json": true, "xml": true, "csv": true,
	}

	// Parse the URL to check both path and query parameters
	u, err := url.Parse(fullURL)
	if err != nil {
		return false
	}

	// Function to extract extension from a filename
	getExt := func(filename string) string {
		parts := strings.Split(filename, ".")
		if len(parts) < 2 {
			return ""
		}
		return strings.ToLower(parts[len(parts)-1])
	}

	// Check the URL path first
	if ext := getExt(u.Path); ext != "" && allowedExts[ext] {
		return true
	}

	// Check query parameters for filenames (like f=filename.ext)
	for _, value := range u.Query() {
		for _, v := range value {
			if ext := getExt(v); ext != "" && allowedExts[ext] {
				return true
			}
		}
	}

	return false
}

func downloadFileToSubDir(fileURL, outputDir, subDir string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fileName := path.Base(resp.Request.URL.Path)

	// Create the target directory (base output + subdirectory)
	var targetDir string
	if subDir != "" {
		targetDir = filepath.Join(outputDir, subDir)
	} else {
		targetDir = outputDir
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %v", targetDir, err)
	}

	outPath := filepath.Join(targetDir, fileName)

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

func downloadFile(fileURL, outputDir string) error {
	return downloadFileToSubDir(fileURL, outputDir, "")
}
