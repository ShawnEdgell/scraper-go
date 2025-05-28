package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url" // For joining URL paths
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type QuoteItem struct {
	Text   string   `json:"text"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
}

// scrapeOneQuotesPage scrapes a single page and returns quotes and the next page URL (if any)
func scrapeOneQuotesPage(pageURL string) ([]QuoteItem, string, error) {
	var quotes []QuoteItem
	var nextPageURL string

	fmt.Printf("Fetching URL: %s\n", pageURL)
	res, err := http.Get(pageURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get URL %s: %w", pageURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, "", fmt.Errorf("request to %s failed with status code %d: %s", pageURL, res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse HTML from %s: %w", pageURL, err)
	}

	doc.Find("div.quote").Each(func(i int, s *goquery.Selection) {
		var currentQuote QuoteItem
		currentQuote.Text = strings.TrimSpace(s.Find("span.text").First().Text())
		currentQuote.Author = strings.TrimSpace(s.Find("small.author").First().Text())
		s.Find("div.tags a.tag").Each(func(j int, tagSelection *goquery.Selection) {
			currentQuote.Tags = append(currentQuote.Tags, strings.TrimSpace(tagSelection.Text()))
		})
		quotes = append(quotes, currentQuote)
	})
	fmt.Printf("Found %d quotes on this page.\n", len(quotes))

	// Find the "Next" page link
	// Common selector for the next page link: "li.next a"
	nextPageRelativePath, exists := doc.Find("li.next a").First().Attr("href")
	if exists {
		// Resolve the relative path to an absolute URL
		base, _ := url.Parse(pageURL) // Use the current pageURL as the base
		nextPageRef, _ := url.Parse(nextPageRelativePath)
		nextPageURL = base.ResolveReference(nextPageRef).String()
		fmt.Printf("Found next page link: %s\n", nextPageURL)
	} else {
		fmt.Println("No next page link found.")
	}

	return quotes, nextPageURL, nil
}

func main() {
	startURL := "http://quotes.toscrape.com/"
	outputQuotesJsonPath := "all_quotes.json"
	var allQuotes []QuoteItem
	currentURL := startURL
	maxPages := 5 // Let's limit to 5 pages for now to be polite

	fmt.Printf("Starting to scrape quotes, max %d pages...\n", maxPages)

	for i := 0; i < maxPages; i++ {
		if currentURL == "" {
			fmt.Println("No more pages to scrape.")
			break
		}
		fmt.Printf("\nScraping page %d: %s\n", i+1, currentURL)
		
		quotesFromPage, nextPageURL, err := scrapeOneQuotesPage(currentURL)
		if err != nil {
			log.Printf("Failed to scrape page %s: %v. Stopping.", currentURL, err)
			break // Stop if a page fails
		}

		allQuotes = append(allQuotes, quotesFromPage...)
		currentURL = nextPageURL // Prepare for the next iteration

		if currentURL == "" && i < maxPages-1 { // No next page, but we haven't hit maxPages
		    fmt.Println("Reached the last page of quotes.")
		    break
		}
	}
	
	fmt.Printf("\nTotal quotes scraped: %d\n", len(allQuotes))

	if len(allQuotes) > 0 {
		jsonData, err := json.MarshalIndent(allQuotes, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal all quotes to JSON: %v", err)
		}

		err = os.WriteFile(outputQuotesJsonPath, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write JSON data to %s: %v", outputQuotesJsonPath, err)
		}
		fmt.Printf("Successfully scraped quote data saved to %s\n", outputQuotesJsonPath)
	} else {
		fmt.Println("No quotes found to save.")
	}

	fmt.Println("All quotes scraping process complete.")
}