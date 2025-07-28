package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type downloadProgress struct {
	total    int
	current  int
	lastFile string
	mu       sync.Mutex
}

func newDownloadProgress(total int) *downloadProgress {
	return &downloadProgress{
		total: total,
	}
}

func (dp *downloadProgress) increment(filename string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.current++
	dp.lastFile = filename
	dp.printProgress()
}

func (dp *downloadProgress) printProgress() {
	percentage := float64(dp.current) / float64(dp.total) * 100
	fmt.Printf("\rDownloading Files: %d/%d (%.1f%%) | Last: %s ",
		dp.current, dp.total, percentage, dp.lastFile)

	if dp.current == dp.total {
		fmt.Println("\nDownload Complete!")
	}
}

func fetchDownloadLinks(radarURL string) ([]string, error) {
	resp, err := http.Get(radarURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("div.bdpLink a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			absoluteURL := resolveURL(radarURL, href)
			links = append(links, absoluteURL)
		}
	})

	return links, nil
}

func resolveURL(baseURL, link string) string {
	base, _ := url.Parse(baseURL)
	relative, _ := url.Parse(link)
	resolvedURL := base.ResolveReference(relative)
	return resolvedURL.String()
}

func downloadFile(url string, outputDir string, progress *downloadProgress, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("\nError downloading %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	filename := filepath.Base(url)
	join := filepath.Join(outputDir, filename)

	out, err := os.Create(join)
	if err != nil {
		fmt.Printf("\nError creating file %s: %v\n", filename, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("\nError writing file %s: %v\n", filename, err)
		return
	}

	progress.increment(filename)
}

func downloadFiles(links []string, outputDir string) []string {
	var wg sync.WaitGroup
	progress := newDownloadProgress(len(links))

	// Increased number of workers to 50
	maxConcurrent := 50
	semaphore := make(chan struct{}, maxConcurrent)
	var mu sync.Mutex
	var downloadedFiles []string

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(url string) {
			defer func() { <-semaphore }()
			downloadFile(url, outputDir, progress, &wg)

			mu.Lock()
			downloadedFiles = append(downloadedFiles, filepath.Base(url))
			mu.Unlock()
		}(link)
	}

	wg.Wait()
	fmt.Println() // New line after progress

	return downloadedFiles
}

func promptInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	radar := promptInput("Enter radar site (KHTX): ")
	month := promptInput("Enter month (03): ")
	day := promptInput("Enter day (15): ")
	year := promptInput("Enter year (2025): ")

	url := fmt.Sprintf("https://www.ncdc.noaa.gov/nexradinv/bdp-download.jsp?id=%s&yyyy=%s&mm=%s&dd=%s&product=AAL2",
		radar, year, month, day)

	outputDir := fmt.Sprintf("%s_%s_%s_%s", radar, year, month, day)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	links, err := fetchDownloadLinks(url)
	if err != nil {
		fmt.Printf("Error fetching download links: %v\n", err)
		return
	}

	downloadedFiles := downloadFiles(links, outputDir)

	fmt.Printf("Total files downloaded: %d\n", len(downloadedFiles))
	fmt.Printf("Files saved in: %s\n", outputDir)
}
