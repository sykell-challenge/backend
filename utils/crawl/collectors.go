package crawl

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"

	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils"
)

// CrawlData contains the main data and a callback for link processing
type CrawlData struct {
	MainData     models.URL
	LinkCount    int
	ProcessLinks func() models.URL
}

func CrawlURL(url string) models.URL {
	result := CrawlURLWithCallback(url)
	return result.ProcessLinks() // Process links and return complete data
}

func CrawlURLWithCallback(url string) CrawlData {

	mainCollector := colly.NewCollector()
	secondCollector := colly.NewCollector()

	var data models.URL
	data.URL = url
	data.Links = models.Links{}
	data.Tags = models.Tags{}

	host := utils.GetHostFromURL(url)

	// Store the main URL's HTTP status code
	var mainStatusCode int

	mainCollector.OnResponse(func(r *colly.Response) {
		mainStatusCode = r.StatusCode
		// Detect HTML version (4 or 5)
		bodyStr := string(r.Body)
		if strings.Contains(bodyStr, "<!DOCTYPE html>") {
			data.HTMLVersion = "5"
		} else if strings.Contains(strings.ToLower(bodyStr), "<!doctype html public") {
			data.HTMLVersion = "4"
		} else {
			data.HTMLVersion = "Unknown"
		}
		fmt.Printf("Main URL %s returned status: %d\n", url, r.StatusCode)
	})

	mainCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
			return
		}

		if strings.HasPrefix(link, "http") {
			secondaryHost := utils.GetHostFromURL(link)

			if secondaryHost != host {
				data.Links = append(data.Links, models.Link{
					Link:       link,
					Type:       "external",
					StatusCode: 0, // Will be updated by secondCollector
				})
			} else {
				data.Links = append(data.Links, models.Link{
					Link:       link,
					Type:       "internal",
					StatusCode: 0, // Will be updated by secondCollector
				})
			}

			secondCollector.Visit(link)
		} else {
			final_link := e.Request.AbsoluteURL(link)
			data.Links = append(data.Links, models.Link{
				Link:       final_link,
				Type:       "internal",
				StatusCode: 0, // Will be updated by secondCollector
			})

			secondCollector.Visit(final_link)
		}
	})

	mainCollector.OnHTML("form", func(e *colly.HTMLElement) {
		// Check if form contains password field (indicates login form)
		passwordField := e.DOM.Find("input[type='password']")
		if passwordField.Length() > 0 {
			data.LoginForm = true
			fmt.Println("Login form detected")
		}
	})

	mainCollector.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		data.Title = strings.TrimSpace(title) // Store the title and trim whitespace
		fmt.Println("Title found: ", title)
	})

	mainCollector.OnHTML("h1, h2, h3, h4, h5, h6, p", func(e *colly.HTMLElement) {
		currentHeaderCount := -1
		for i, tag := range data.Tags {
			if tag.TagName == e.Name {
				currentHeaderCount = i
				break
			}
		}

		if currentHeaderCount != -1 {
			data.Tags[currentHeaderCount].Count++
		} else {
			data.Tags = append(data.Tags, models.Tag{
				TagName: e.Name,
				Count:   1,
			})
		}
	})

	mainCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	mainCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error visiting URL: ", r.Request.URL.String(), " - ", err)
	})

	secondCollector.OnResponse(func(r *colly.Response) {
		// Update the status code for the corresponding link
		requestURL := r.Request.URL.String()
		fmt.Printf("Link %s returned status: %d\n", requestURL, r.StatusCode)

		// Find and update the corresponding link in our data.Links slice
		for i := range data.Links {
			if data.Links[i].Link == requestURL {
				data.Links[i].StatusCode = r.StatusCode
				break
			}
		}
	})

	secondCollector.OnError(func(r *colly.Response, err error) {
		requestURL := r.Request.URL.String()
		fmt.Printf("Error visiting link %s: %v\n", requestURL, err)

		// Update existing link with error status or add as inaccessible
		found := false
		for i := range data.Links {
			if data.Links[i].Link == requestURL {
				data.Links[i].Type = "inaccessible"
				data.Links[i].StatusCode = r.StatusCode
				found = true
				break
			}
		}

		// If not found in existing links, add as inaccessible
		if !found {
			data.Links = append(data.Links, models.Link{
				Link:       requestURL,
				Type:       "inaccessible",
				StatusCode: r.StatusCode,
			})
		}
	})

	fmt.Println("Crawling URL: ", url)
	fmt.Println("Host: ", host)

	mainCollector.Visit(url)
	// mainCollector.Wait()

	// Store the HTTP status code for the main URL
	data.StatusCode = mainStatusCode

	// Prepare the main data (without processed links)
	mainData := data
	mainData.Links = models.Links{} // Clear links for half-completed state
	linkCount := len(data.Links)

	// Return CrawlData with main data and link processing function
	return CrawlData{
		MainData:  mainData,
		LinkCount: linkCount,
		ProcessLinks: func() models.URL {
			// This function will process all the links when called
			fmt.Printf("Processing %d links...\n", linkCount)

			// Wait for secondary collector to process all links
			// Note: data.Links already contains the links from mainCollector
			// secondCollector will update their StatusCode values
			secondCollector.Wait()

			fmt.Printf("Crawl completed. Title: '%s', HTTP Status: %d, Found %d links, %d tag types, Login form: %t\n",
				data.Title, mainStatusCode, len(data.Links), len(data.Tags), data.LoginForm)

			return data
		},
	}
}
