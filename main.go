package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

const (
	downloadDir   = "Scientific America 60s"
	transcriptDir = downloadDir + "/transcripts"
	audioDir      = downloadDir + "/audios"
)

type page struct {
	no    int
	items []*item
}

func createPage(wg *sync.WaitGroup, no int, pageURL string) *page {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	items := make([]*item, 0)

	// big playing area
	c.OnHTML(`div[class="podcasts container"]`, func(e *colly.HTMLElement) {
		item := createItemFromBig(e)
		items = append(items, item)
		wg.Add(1)
		go item.downloadTranscript(wg)
		wg.Add(1)
		go item.downloadMP3(wg)
	})

	// grid item
	c.OnHTML(`div[class="grid__col large-1-2 xlarge-1-2 medium-1-2 small-no-pad"]`, func(e *colly.HTMLElement) {
		item := createItemFromGrid(e)
		items = append(items, item)
		wg.Add(1)
		go item.downloadTranscript(wg)
		wg.Add(1)
		go item.downloadMP3(wg)
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

func (it *item) downloadMP3(wg *sync.WaitGroup) {
	defer wg.Done()
	if resp, err := http.Get(it.urlMp3); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		if file, err := os.Create(fmt.Sprintf("%s/%s.mp3", audioDir, fuckWindowsFS(it.title))); err != nil {
			panic(err)
		} else {
			io.Copy(file, resp.Body)
			file.Sync()
			file.Close()
		}
	}
}

func (it *item) downloadTranscript(wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	// full transcript
	c.OnHTML(`div[id="transcripts-body"]`, func(e *colly.HTMLElement) {
		content := ""
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			content += el.Text + "\n\n"
		})
		it.transcript = content
		// save transcript
		if file, err := os.Create(fmt.Sprintf("%s/%s.txt", transcriptDir, fuckWindowsFS(it.title))); err != nil {
			panic(err)
		} else {
			file.WriteString(it.transcript)
			file.Sync()
			file.Close()
		}
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

var replacer = strings.NewReplacer(
	"<", "＜",
	">", "＞",
	":", "：",
	"\"", "＂",
	"/", "⁄",
	"\\", "﹨",
	"|", "│",
	"?", "﹖",
	"*", "﹡",
)

func fuckWindowsFS(s string) string {
	return replacer.Replace(s)
}

func checkDownloadDir() {
	if err := os.MkdirAll(transcriptDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(audioDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func main() {
	checkDownloadDir()
	count := 2
	flag.Parse()
	if len(flag.Args()) > 0 {
		if countTmp, err := strconv.Atoi(flag.Arg(0)); err != nil {
			panic(err)
		} else {
			count = countTmp + 1
		}
	}

	pageNo := 1
	pages := make([]*page, 0)
	waitGroup := sync.WaitGroup{}

	start := time.Now()
	for ; pageNo < count; pageNo++ {
		waitGroup.Add(1)
		go func(pageNo int) {
			defer waitGroup.Done()
			page := createPage(&waitGroup, pageNo, fmt.Sprintf("https://www.scientificamerican.com/podcast/60-second-science/?page=%d", pageNo))
			pages = append(pages, page)
		}(pageNo)
	}
	waitGroup.Wait()
	end := time.Now()
	fmt.Printf("Totally used %d s. Parsed %d pages (%d articles)\n", end.Sub(start)/time.Second, len(pages), len(pages)*len(pages[0].items))
}
