package main

import (
	"fmt"      // For printing formatted output
	"log"      // For logging errors
	"net/http" // For making HTTP requests
	"strings"  // For string manipulation like TrimSpace

	"github.com/PuerkitoBio/goquery" // Our HTML parsing library
)

// scrapeWebsite fetches a URL and extracts its title.
func scrapeWebsite(url string) (string, error) {
	// Make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get URL %s: %w", url, err)
	}
	defer res.Body.Close() // Ensure the response body is closed

	// Check for successful status code
	if res.StatusCode != 200 {
		return "", fmt.Errorf("request failed with status code %d: %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML from %s: %w", url, err)
	}

	// Find the title element and get its text
	// For example.com, the title is within the <h1> tag
	// For a more general approach, you'd use "title" tag within <head>
	// title := doc.Find("head title").First().Text()
	title := doc.Find("h1").First().Text() // example.com's main heading

	return strings.TrimSpace(title), nil
}

func main() {
	targetURL := "http://example.com"
	fmt.Printf("Attempting to scrape title from: %s\n", targetURL)

	title, err := scrapeWebsite(targetURL)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err) // log.Fatalf will print the error and exit
	}

	fmt.Printf("Successfully scraped title: '%s'\n", title)
}