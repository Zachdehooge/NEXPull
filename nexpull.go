package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type downloadProgress struct {
	total    int
	current  int
	lastFile string
	mu       sync.Mutex
}

type ListBucketResult struct {
	Contents []struct {
		Key string `xml:"Key"`
	} `xml:"Contents"`
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
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result ListBucketResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	var links []string
	for _, item := range result.Contents {
		url := "https://unidata-nexrad-level2.s3.amazonaws.com/" + item.Key
		links = append(links, url)
	}

	return links, nil
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

	// 50 Workers
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
	fmt.Println()

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

	/*
	 !Search for Files - https://unidata-nexrad-level2.s3.amazonaws.com/?prefix=2020/03/20/KHTX/
	 !Download URL https://unidata-nexrad-level2.s3.amazonaws.com/2025/03/15/KHTX/KHTX20250315_001350_V06
	*/

	baseURL := fmt.Sprintf("https://unidata-nexrad-level2.s3.amazonaws.com/?prefix=%s/%s/%s/%s/", year, month, day, strings.ToUpper(radar))

	outputDir := fmt.Sprintf("%s_%s_%s_%s", strings.ToUpper(radar), year, month, day)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	links, err := fetchDownloadLinks(baseURL)
	if err != nil {
		return
	}

	downloadedFiles := downloadFiles(links, outputDir)

	fmt.Printf("Total files downloaded: %d\n", len(downloadedFiles))
	fmt.Printf("Files saved in: %s\n", outputDir)
}
