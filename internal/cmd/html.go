package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/OrkWard/wormface/internal/utils"
)

func CmdHTML(args []string) {
	htmlFlags := flag.NewFlagSet("wormface html", flag.ExitOnError)
	outputDir := htmlFlags.String("d", "./output/html", "Directory to save downloaded images")
	htmlFlags.Parse(args)

	positionArgs := htmlFlags.Args()
	if len(positionArgs) < 1 {
		fmt.Println("Usage: wormface-cli html [options] <page_url>")
		htmlFlags.PrintDefaults()
		os.Exit(1)
	}
	pageURL := positionArgs[0]

	client := http.DefaultClient

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("[ERROR] Failed to create output dir %s: %v\n", *outputDir, err)
		os.Exit(1)
	}

	fmt.Printf("[INPUT] Fetching page: %s\n", pageURL)
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		fmt.Printf("[ERROR] Failed to build request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch page: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("[ERROR] Unexpected status code: %s\n", resp.Status)
		os.Exit(1)
	}

	baseURL := resp.Request.URL

	imgURLs, err := extractImageSources(baseURL, resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse html: %v\n", err)
		os.Exit(1)
	}

	if len(imgURLs) == 0 {
		fmt.Println("[WARN] No images found on the page.")
		return
	}

	fmt.Printf("[INFO] Found %d images. Starting download...\n", len(imgURLs))
	utils.DownloadAllWithClient(imgURLs, *outputDir, client)
}

func extractImageSources(baseURL *url.URL, reader io.Reader) ([]string, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	var urls []string
	seen := make(map[string]struct{})

	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					resolved := resolveURL(baseURL, strings.TrimSpace(attr.Val))
					if resolved == "" {
						continue
					}
					if _, ok := seen[resolved]; ok {
						continue
					}
					seen[resolved] = struct{}{}
					urls = append(urls, resolved)
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}

	walker(doc)
	return urls, nil
}

func resolveURL(baseURL *url.URL, raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	if baseURL != nil && !parsed.IsAbs() {
		parsed = baseURL.ResolveReference(parsed)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	return parsed.String()
}
