package crawl_manager

import (
	"slices"
	"strings"
	"sykell-challenge/backend/models"
	"sykell-challenge/backend/utils"
)

func (pm *CrawlManager) processLinks() {
	slices.Sort(pm.linksFound)

	pm.linksFound = slices.Compact(pm.linksFound)

	for _, link := range pm.linksFound {
		if strings.HasPrefix(link, "/") {
			link = pm.currentHost + link
		}

		if isAvailable, statusCode := utils.IsURLAvailable(link); isAvailable {
			pm.data.Links = append(pm.data.Links, models.Link{Link: link, Type: pm.determineLinkType(link), StatusCode: statusCode})
		} else {
			pm.data.Links = append(pm.data.Links, models.Link{Link: link, Type: "inaccessible", StatusCode: statusCode})
		}
	}
}

func (pm *CrawlManager) determineLinkType(link string) string {
	if strings.HasPrefix(link, "http") {
		linkHost := utils.GetHostFromURL(link)
		if linkHost != pm.currentHost {
			return "external"
		}
	}
	return "internal"
}
