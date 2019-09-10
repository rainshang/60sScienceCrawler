package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	page := 1
	c := colly.NewCollector()

	// big playing area
	c.OnHTML(`div[class="podcasts container"]`, func(e *colly.HTMLElement) {
		title := e.ChildAttr(`div[class="podcasts-header podcasts-header--feature tooltip-outer"]`, "data-podcast-title")
		transcriptURL := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
		mp3 := e.ChildAttr(`a[class="podcast-download"]`, "href")
		fmt.Printf("BIG Title: %s, transcript url: %s, mp3 url: %s\n", title, transcriptURL, mp3)
	})

	// grid item
	c.OnHTML(`div[class="grid__col large-1-2 xlarge-1-2 medium-1-2 small-no-pad"]`, func(e *colly.HTMLElement) {
		title := e.Attr("data-podcast-title")
		mp3 := e.ChildAttr(`a[data-tooltip-bounds-id="podcast-group"]`, "href")
		transcriptURL := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
		fmt.Printf("title: %s, transcript url: %s, mp3 url: %s\n", title, transcriptURL, mp3)
	})

	c.Visit(fmt.Sprintf("https://www.scientificamerican.com/podcast/60-second-science/?page=%d", page))
}
