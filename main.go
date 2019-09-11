package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type item struct {
	Title         string
	URLTranscript string
	transcript    string
	URLMp3        string
}

func (it *item) downloadMp3() {

}

func (it *item) downloadTranscript() {
	// full transcript
	// c.OnHTML(`div[id="transcripts-body"]`, func(e *colly.HTMLElement) {
	// 	content := ""
	// 	e.ForEach("p", func(_ int, el *colly.HTMLElement) {
	// 		content += el.Text + "\n\n"
	// 	})
	// })
}

func createItemFromBig(e *colly.HTMLElement) item {
	title := e.ChildAttr(`div[class="podcasts-header podcasts-header--feature tooltip-outer"]`, "data-podcast-title")
	transcript := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
	mp3 := e.ChildAttr(`a[class="podcast-download"]`, "href")
	return item{
		Title:         title,
		URLTranscript: transcript,
		URLMp3:        mp3,
	}
}

func createItemFromGrid(e *colly.HTMLElement) item {
	title := e.Attr("data-podcast-title")
	mp3 := e.ChildAttr(`a[data-tooltip-bounds-id="podcast-group"]`, "href")
	transcript := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
	return item{
		Title:         title,
		URLTranscript: transcript,
		URLMp3:        mp3,
	}
}

func main() {
	page := 1
	items := make([]item, 0)
	c := colly.NewCollector()

	// big playing area
	c.OnHTML(`div[class="podcasts container"]`, func(e *colly.HTMLElement) {
		items = append(items, createItemFromBig(e))
	})

	// grid item
	c.OnHTML(`div[class="grid__col large-1-2 xlarge-1-2 medium-1-2 small-no-pad"]`, func(e *colly.HTMLElement) {
		items = append(items, createItemFromGrid(e))
	})

	start := time.Now()
	for ; page < 3; page++ {
		c.Visit(fmt.Sprintf("https://www.scientificamerican.com/podcast/60-second-science/?page=%d", page))
	}
	end := time.Now()
	fmt.Printf("totally used %d s", end.Sub(start)/time.Second)
}
