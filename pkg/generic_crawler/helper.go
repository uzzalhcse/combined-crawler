package generic_crawler

import "strings"

func (gc *GenericCrawler) GetFullUrl(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// If url is already a full Url, return it as is
		return url
	} else if strings.HasPrefix(url, "//") {
		// If url starts with "//", use the protocol from BaseUrl
		if strings.HasPrefix(gc.BaseUrl, "https://") {
			return "https:" + url
		}
		return "http:" + url
	}
	// Otherwise, concatenate with BaseUrl
	return gc.BaseUrl + url
}
