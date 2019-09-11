package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type item struct {
	Title         string
	URLTranscript string
	URLMp3        string
}

func (it *item) downloadMp3() {

}

func (it *item) downloadTranscript() {

}

func main() {
	page := 1
	count := 1
	c := colly.NewCollector()

	// full transcript
	c.OnHTML(`div[id="transcripts-body"]`, func(e *colly.HTMLElement) {
		content := ""
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n\n"
		})
		fmt.Print(content)
	})

	// big playing area
	c.OnHTML(`div[class="podcasts container"]`, func(e *colly.HTMLElement) {
		title := e.ChildAttr(`div[class="podcasts-header podcasts-header--feature tooltip-outer"]`, "data-podcast-title")
		transcriptURL := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
		mp3 := e.ChildAttr(`a[class="podcast-download"]`, "href")
		fmt.Printf("%d.BIG Title: %s, transcript url: %s, mp3 url: %s\n\n", count, title, transcriptURL, mp3)
		count++
	})

	// grid item
	c.OnHTML(`div[class="grid__col large-1-2 xlarge-1-2 medium-1-2 small-no-pad"]`, func(e *colly.HTMLElement) {
		title := e.Attr("data-podcast-title")
		mp3 := e.ChildAttr(`a[data-tooltip-bounds-id="podcast-group"]`, "href")
		transcriptURL := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
		fmt.Printf("%d.title: %s, transcript url: %s, mp3 url: %s\n\n", count, title, transcriptURL, mp3)
		count++
	})

	start := time.Now()
	for ; page < 2; page++ {
		c.Visit(fmt.Sprintf("https://www.scientificamerican.com/podcast/60-second-science/?page=%d", page))
	}
	end := time.Now()
	fmt.Printf("totally used %d s", end.Sub(start)/time.Second)
}
