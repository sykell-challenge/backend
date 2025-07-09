package crawl

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"

	"sykell-challenge/backend/utils"
)

func CrawlURL(g *gin.Context) {

	mainCollector := colly.NewCollector()
	secondCollector := colly.NewCollector()

	var data struct {
		InternalLinks []string `json:"internal_links"`
		ExternalLinks []string `json:"external_links"`
		BrokenLinks   []string `json:"broken_links"`
		Title         string   `json:"title"`
		Stats         struct {
			NumberOfInternalLinks int `json:"number_of_internal_links"`
			NumberOfExternalLinks int `json:"number_of_external_links"`
			NumberOfBrokenLinks   int `json:"number_of_broken_links"`
		} `json:"stats"`
		Tags []struct {
			TagName string `json:"tag_name"`
			Count   int    `json:"count"`
		} `json:"tags"`
	}

	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println("Crawling URL: ", req.URL)

	host := utils.GetHostFromURL(req.URL)

	fmt.Println("Host: ", host)

	secondCollector.OnError(func(r *colly.Response, err error) {

		data.BrokenLinks = append(data.BrokenLinks, r.Request.URL.String())
	})

	mainCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
			return
		}

		if strings.HasPrefix(link, "http") {
			secondaryHost := utils.GetHostFromURL(link)

			if secondaryHost != host {
				data.ExternalLinks = append(data.ExternalLinks, link)
			} else {
				data.InternalLinks = append(data.InternalLinks, link)
			}

			secondCollector.Visit(link)
		} else {
			final_link := e.Request.AbsoluteURL(link)
			data.InternalLinks = append(data.InternalLinks, final_link)

			secondCollector.Visit(final_link)

		}

	})

	mainCollector.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text

		data.Title = title

		fmt.Println("Title found: ", title)
	})

	mainCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	url := req.URL
	mainCollector.Visit(url)
	mainCollector.Wait()

	data.Stats.NumberOfInternalLinks = len(data.InternalLinks)
	data.Stats.NumberOfExternalLinks = len(data.ExternalLinks)
	data.Stats.NumberOfBrokenLinks = len(data.BrokenLinks)

	data.InternalLinks = utils.Difference(data.InternalLinks, data.BrokenLinks)
	data.ExternalLinks = utils.Difference(data.ExternalLinks, data.BrokenLinks)

	g.JSON(200, gin.H{
		"message": "Crawling started for " + url,
		"data":    data,
	})
}
