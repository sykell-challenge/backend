package crawl

type Tag struct {
	TagName string `json:"tagName"`
	Count   int    `json:"count"`
}

type Stats struct {
	NumberOfInternalLinks int `json:"numberOfInternalLinks"`
	NumberOfExternalLinks int `json:"numberOfExternalLinks"`
	NumberOfBrokenLinks   int `json:"numberOfBrokenLinks"`
}

type CrawlResponse struct {
	InternalLinks []string `json:"internalLinks"`
	ExternalLinks []string `json:"externalLinks"`
	BrokenLinks   []string `json:"brokenLinks"`
	Title         string   `json:"title"`
	Stats         Stats    `json:"stats"`
	Tags          []Tag    `json:"tags"`
}
