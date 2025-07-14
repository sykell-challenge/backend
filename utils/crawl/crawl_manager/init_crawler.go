package crawl_manager

import (
	"fmt"
	"strings"
	"sykell-challenge/backend/models"

	"github.com/gocolly/colly"
)

func (cm *CrawlManager) initCrawler() {

	cm.collector = colly.NewCollector()

	cm.collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	cm.collector.OnResponse(func(r *colly.Response) {
		cm.ProcessMainResponse(r)
	})

	cm.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		cm.linksFound = append(cm.linksFound, link)
	})

	cm.collector.OnHTML("form", func(e *colly.HTMLElement) {
		cm.ProcessForm(e)
	})

	cm.collector.OnHTML("title", func(e *colly.HTMLElement) {
		cm.ProcessTitle(e)
	})

	cm.collector.OnHTML("h1, h2, h3, h4, h5, h6, p", func(e *colly.HTMLElement) {
		cm.ProcessTag(e)
	})

	cm.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	cm.collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error visiting URL: ", r.Request.URL.String(), " - ", err)
	})
}

// ProcessTag processes a single HTML tag element and updates the tag count
func (cm *CrawlManager) ProcessTag(e *colly.HTMLElement) {
	tagName := e.Name
	cm.incrementTagCount(tagName)
}

// ProcessForm checks for login forms and updates the LoginForm field
func (cm *CrawlManager) ProcessForm(e *colly.HTMLElement) {
	// Check if form contains password field (indicates login form)
	passwordField := e.DOM.Find("input[type='password']")
	if passwordField.Length() > 0 {
		cm.data.LoginForm = true
		fmt.Println("Login form detected")
	}
}

// ProcessTitle extracts and stores the page title
func (cm *CrawlManager) ProcessTitle(e *colly.HTMLElement) {
	title := e.Text
	cm.data.Title = strings.TrimSpace(title) // Store the title and trim whitespace
	fmt.Println("Title found: ", title)
}

// ProcessMainResponse handles the main URL response and detects HTML version
func (cm *CrawlManager) ProcessMainResponse(r *colly.Response) {
	cm.data.StatusCode = r.StatusCode
	// Detect HTML version (4 or 5)
	bodyStr := string(r.Body)
	if strings.Contains(bodyStr, "<!DOCTYPE html>") {
		cm.data.HTMLVersion = "5"
	} else if strings.Contains(strings.ToLower(bodyStr), "<!doctype html public") {
		cm.data.HTMLVersion = "4"
	} else {
		cm.data.HTMLVersion = "Unknown"
	}
}

// Private helper methods

// incrementTagCount increments the count for a specific tag or adds it if not found
func (cm *CrawlManager) incrementTagCount(tagName string) {
	// Look for existing tag
	for i, tag := range cm.data.Tags {
		if tag.TagName == tagName {
			cm.data.Tags[i].Count++
			return
		}
	}

	// Tag not found, add new tag with count 1
	cm.data.Tags = append(cm.data.Tags, models.Tag{
		TagName: tagName,
		Count:   1,
	})
}

// shouldSkipLink checks if a link should be skipped
func (cm *CrawlManager) shouldSkipLink(link string) bool {
	return link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:")
}
