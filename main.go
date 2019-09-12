package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type page struct {
	no    int
	items []item
}

func createPage(no int, pageURL string) *page {
	c := colly.NewCollector()
	items := make([]item, 0)

	// big playing area
	c.OnHTML(`div[class="podcasts container"]`, func(e *colly.HTMLElement) {
		item := createItemFromBig(e)
		go item.downloadTranscript()
		items = append(items, *item)
	})

	// grid item
	c.OnHTML(`div[class="grid__col large-1-2 xlarge-1-2 medium-1-2 small-no-pad"]`, func(e *colly.HTMLElement) {
		item := createItemFromGrid(e)
		go item.downloadTranscript()
		items = append(items, *item)
	})

	c.Visit(pageURL)
	return &page{
		no,
		items,
	}
}

type item struct {
	title         string
	urlTranscript string
	transcript    string
	urlMp3        string
}

func (it *item) downloadMp3() {

}

func (it *item) downloadTranscript() {
	c := colly.NewCollector()

	// full transcript
	c.OnHTML(`div[id="transcripts-body"]`, func(e *colly.HTMLElement) {
		content := ""
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n\n"
		})
		it.transcript = content
	})

	c.Visit(it.urlTranscript)
}

func createItemFromBig(e *colly.HTMLElement) *item {
	title := e.ChildAttr(`div[class="podcasts-header podcasts-header--feature tooltip-outer"]`, "data-podcast-title")
	transcript := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
	mp3 := e.ChildAttr(`a[class="podcast-download"]`, "href")
	return &item{
		title:         title,
		urlTranscript: transcript,
		urlMp3:        mp3,
	}
}

func createItemFromGrid(e *colly.HTMLElement) *item {
	title := e.Attr("data-podcast-title")
	mp3 := e.ChildAttr(`a[data-tooltip-bounds-id="podcast-group"]`, "href")
	transcript := e.ChildAttr(`a[class="t_meta underlined_text t_small"]`, "href")
	return &item{
		title:         title,
		urlTranscript: transcript,
		urlMp3:        mp3,
	}
}

func main() {
	pageNo := 1
	pages := make([]page, 0)
	waitGroup := sync.WaitGroup{}

	start := time.Now()
	for ; pageNo < 3; pageNo++ {
		waitGroup.Add(1)
		go func(pageNo int) {
			page := createPage(pageNo, fmt.Sprintf("https://www.scientificamerican.com/podcast/60-second-science/?page=%d", pageNo))
			pages = append(pages, *page)
			waitGroup.Done()
		}(pageNo)
	}
	waitGroup.Wait()
	end := time.Now()
	fmt.Printf("totally used %d s", end.Sub(start)/time.Second)
}
