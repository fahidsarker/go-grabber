package main

import (
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

// var allowedExt = []string{".zip", ".pdf", ".jpg", ".png", ".txt", ".mp4", ".mp3", ".doc", ".docx", ".xls", ".xlsx", ".mkv", ".avi", ".webm", ".gif", ".bmp", ".jpeg", ".svg", ".ico", ".wav", ".m4a", ".flac", ".ogg", ".wmv", ".3gp", ".mov", ".flv", ".swf", ".mkv", ".avi", ".webm", ".gif", ".bmp", ".jpeg", ".svg", ".ico", ".wav", ".m4a", ".flac", ".ogg", ".wmv", ".3gp", ".mov", ".flv", ".swf"}

func main() {
	// Flags
	workers := flag.Int("workers", 4, "Number of concurrent download workers")
	flag.Parse()

	// Args
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go [--workers=N] <URL> <output_dir>")
		return
	}
	startURL := args[0]
	outputDir := args[1]

	// Prepare output directory
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output dir: %v\n", err)
		return
	}

	// Fetch and parse HTML
	resp, err := http.Get(startURL)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status code: %d\n", resp.StatusCode)
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	baseURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Invalid base URL: %v\n", err)
		return
	}

	links := extractFileLinks(doc, baseURL)
	fmt.Printf("Found %d downloadable files\n", len(links))

	// Start concurrent downloads
	var wg sync.WaitGroup
	jobs := make(chan string)

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for link := range jobs {
				fmt.Printf("[Worker %d] Downloading: %s\n", workerID, link)
				err := downloadFile(link, outputDir)
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

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if u, err := base.Parse(href); err == nil {
						if hasAllowedExt(u.Path) {
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

	return links
}

func hasAllowedExt(p string) bool {
	ext := strings.ToLower(filepath.Ext(p))
	// for _, allowed := range allowedExt {
	// 	if ext == allowed {
	// 		return true
	// 	}
	// }
	// return false
	return ext != "" && ext[1:] != ""
}

func downloadFile(fileURL, outputDir string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fileName := path.Base(resp.Request.URL.Path)
	outPath := filepath.Join(outputDir, fileName)

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
