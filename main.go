package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Title       string
	PublishDate string
	Author      string
	Content     []string
	Links       []string
}

func scrapeLocalFile(filePath string) error {
	htmlContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read HTML file %s: %w", filePath, err)
	}
	htmlContent := string(htmlContentBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("failed to parse HTML from %s: %w", filePath, err)
	}

	pageTitle := strings.TrimSpace(doc.Find("head title").First().Text())
	fmt.Printf("Page Title: %s\n", pageTitle)

	mainHeading := strings.TrimSpace(doc.Find("h1#pageHeader").First().Text())
	fmt.Printf("Main Heading: %s\n\n", mainHeading)

	fmt.Println("Navigation Links:")
	doc.Find("header nav ul li a.nav-link").Each(func(i int, s *goquery.Selection) {
		linkText := strings.TrimSpace(s.Text())
		linkHref, _ := s.Attr("href")
		fmt.Printf("  - Text: %s, Href: %s\n", linkText, linkHref)
		if s.HasClass("active") {
			fmt.Printf("    (This link is active)\n")
		}
	})
	fmt.Println("")

	fmt.Println("Articles:")
	doc.Find("article.post").Each(func(i int, articleSelection *goquery.Selection) {
		var currentArticle Article

		currentArticle.Title = strings.TrimSpace(articleSelection.Find("h2.post-title").First().Text())
		currentArticle.PublishDate = strings.TrimSpace(articleSelection.Find("p.post-meta span.date").First().Text())
		currentArticle.Author = strings.TrimSpace(articleSelection.Find("p.post-meta span.author").First().Text())

		articleSelection.Find("div.post-content p").Each(func(j int, pSelection *goquery.Selection) {
			currentArticle.Content = append(currentArticle.Content, strings.TrimSpace(pSelection.Text()))
		})

		articleSelection.Find("div.post-content a").Each(func(k int, aSelection *goquery.Selection) {
			link, _ := aSelection.Attr("href")
			currentArticle.Links = append(currentArticle.Links, link)
		})
		
		fmt.Printf("--- Article %d ---\n", i+1)
		fmt.Printf("  Title: %s\n", currentArticle.Title)
		fmt.Printf("  Author: %s\n", currentArticle.Author)
		fmt.Printf("  Date: %s\n", currentArticle.PublishDate)
		fmt.Println("  Content Paragraphs:")
		for _, p := range currentArticle.Content {
			fmt.Printf("    - %s\n", p)
		}
		fmt.Println("  Links in Content:")
		for _, l := range currentArticle.Links {
			fmt.Printf("    - %s\n", l)
		}

		specialData, exists := articleSelection.Find("ul li span[data-id='point3-detail']").Attr("data-id")
		if exists {
			fmt.Printf("  Special Data Attribute: %s\n", specialData)
		}
		fmt.Println("------------------")
	})

	copyrightText := strings.TrimSpace(doc.Find("footer p#copyright").First().Text())
	fmt.Printf("\nFooter Copyright: %s\n", copyrightText)

	return nil
}

func main() {
	filePath := "test.html"
	fmt.Printf("Attempting to scrape local file: %s\n", filePath)

	err := scrapeLocalFile(filePath)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Println("\nSuccessfully finished scraping local file.")
}