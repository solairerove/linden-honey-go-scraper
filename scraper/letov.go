package scrapper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"golang.org/x/text/encoding/charmap"
)

const (
	allowedDomain = "www.gr-oborona.ru"
	textsPage     = "http://www.gr-oborona.ru/texts/"
	maxDepth      = 1
)

// Song is a entity to describe grabber results.
type Song struct {
	Title  string   `json:"title,omitempty"`
	Link   string   `json:"link,omitempty"`
	Author string   `json:"author,omitempty"`
	Album  string   `json:"album,omitempty"`
	Verses []string `json:"verses,omitempty"`
}

// ScrapLetov poor Letov
func ScrapLetov() []Song {

	// var song res.Song
	var songs []Song

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domain: www.gr-oborona.ru
		colly.AllowedDomains(allowedDomain),

		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)

	extensions.RandomUserAgent(c)

	songCollector := c.Clone()

	// Limit the maximum parallelism to cpu num
	err := c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: runtime.NumCPU()})
	if err != nil {
		log.Fatalf("Something wrong with main collector limit %v, %v", runtime.NumCPU(), err)
	}

	// On each existing text find and visit link
	c.OnHTML(`ul[id=abc_list]`, func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
			elem.ForEach("a", func(_ int, link *colly.HTMLElement) {
				fullLink := fmt.Sprintf("http://%s/text_print.php?area=go_texts&id=%s",
					allowedDomain, link.Attr("href")[7:17])

				err := songCollector.Visit(fullLink)
				if err != nil {
					log.Fatalf("Something wrong with song collector visiting %s, %v", fullLink, err)
				}
			})
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Limit the maximum parallelism to cpu num
	err = songCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: runtime.NumCPU()})
	if err != nil {
		log.Fatalf("Something wrong with song collector limit %v, %v", runtime.NumCPU(), err)
	}

	// Before making a request print "Visiting ..."
	songCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	songCollector.OnHTML(`body`, func(e *colly.HTMLElement) {
		var currentSong Song
		currentSong.Link = e.Request.URL.String()

		title := e.ChildText("h2")
		decodedTitle := decodeWindows1251([]byte(title))
		currentSong.Title = string(decodedTitle)

		e.ForEach("p", func(i int, elem *colly.HTMLElement) {
			if i == 0 {
				decodedAuthor := decodeWindows1251([]byte(elem.Text))
				currentSong.Author = string(decodedAuthor)
			}
			if i == 1 {
				decodedAlbum := decodeWindows1251([]byte(elem.Text))
				if strings.Contains(string(decodedAlbum), "Альбом") {
					currentSong.Album = string(decodedAlbum)
				} else {
					processLyrics(elem, &currentSong)
				}
			}
			if i == 2 {
				processLyrics(elem, &currentSong)
			}
		})

		songs = append(songs, currentSong)
		log.Println(len(songs))
	})

	// Start scraping on http://www.gr-oborona.ru/texts/
	err = c.Visit(textsPage)
	if err != nil {
		log.Fatalf("Something wrong with main collector visiting %s, %v", textsPage, err)
	}

	// Wait until threads are finished
	c.Wait()
	songCollector.Wait()

	return songs
}

func processLyrics(elem *colly.HTMLElement, s *Song) {
	elem.DOM.Each(func(_ int, lyrics *goquery.Selection) {
		lyrics.Contents().Each(func(i int, lyric *goquery.Selection) {
			processLyric(lyric, s)
		})
	})
}

func processLyric(lyric *goquery.Selection, s *Song) {
	// handle spaces
	if !lyric.Is("br") {
		//trimmedLyric := strings.TrimSpace(lyric.Text())
		commaLyric := regexp.MustCompile(`&#39;`).ReplaceAllString(lyric.Text(), "'")
		// ⌥ Opt+Space
		nbspLyric := regexp.MustCompile(` `).ReplaceAllString(commaLyric, " ")
		decodedLyric := decodeWindows1251([]byte(nbspLyric))
		s.Verses = append(s.Verses, string(decodedLyric))
		//log.Println(string(decodedLyric), "-", i)
	} else {
		// TODO: handle new lines?
	}
}

// decode shitty cp1251 to human readalbe utf-8
func decodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}
