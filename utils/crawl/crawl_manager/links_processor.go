package crawl_manager

import (
	"slices"
	"strings"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils"
)

func (cm *CrawlManager) processLinks() {
	slices.Sort(cm.linksFound)

	cm.linksFound = slices.Compact(cm.linksFound)

	for _, link := range cm.linksFound {
		if strings.HasPrefix(link, "/") {
			link = cm.currentHost + link
		}

		if isAvailable, statusCode := utils.IsURLAvailable(link); isAvailable {
			cm.data.Links = append(cm.data.Links, models.Link{Link: link, Type: cm.determineLinkType(link), StatusCode: statusCode})
		} else {
			cm.data.Links = append(cm.data.Links, models.Link{Link: link, Type: "inaccessible", StatusCode: statusCode})
		}
	}
}

func (cm *CrawlManager) determineLinkType(link string) string {
	if strings.HasPrefix(link, "http") {
		linkHost := utils.GetHostFromURL(link)
		if linkHost != cm.currentHost {
			return "external"
		}
	}
	return "internal"
}
