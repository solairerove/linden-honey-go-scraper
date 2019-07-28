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
	res "github.com/solairerove/linden-honey-go-scraper/response"
	"golang.org/x/text/encoding/charmap"
)

const (
	allowedDomain = "www.gr-oborona.ru"
	textsPage     = "http://www.gr-oborona.ru/texts/"
	maxDepth      = 1
)

// TODO move to separate file
type myRexexp struct {
	*regexp.Regexp
}

// TODO move to separate file as well
func (r *myRexexp) findStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}

	return captures
}

// ScrapLetov poor Letov
func ScrapLetov() []res.Song {

	// var song res.Song
	var songs []res.Song

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domain: www.gr-oborona.ru
		colly.AllowedDomains(allowedDomain),

		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(maxDepth),
		colly.Async(true),

		// Visit only root url and urls which start with "text" on www.gr-oborona.ru
		// colly.URLFilters(
		// regexp.MustCompile(textsPage),
		// ),
	)

	extensions.RandomUserAgent(c)

	songCollector := c.Clone()

	// Limit the maximum parallelism to cpu num
	//c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: runtime.NumCPU()})

	// On each existing text find and visit link
	//c.OnHTML(`ul[id=abc_list]`, func(e *colly.HTMLElement) {
	//	e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
	//		elem.ForEach("a", func(_ int, link *colly.HTMLElement) {
	//			// fmt.Println(link.Attr("href")[7:17])
	//			fullLink := fmt.Sprintf("http://%s/text_print.php?area=go_texts&id=%s", allowedDomain, link.Attr("href")[7:17])
	//			// log.Println(fullLink)
	//
	//			songCollector.Visit(fullLink)
	//		})
	//	})
	//})

	fullLink := fmt.Sprintf("http://%s/text_print.php?area=go_texts&id=%s", allowedDomain, "1056965230")
	songCollector.Visit(fullLink)

	// On every a element which has href attribute call callback
	// c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
	// link := e.Attr("href")

	// ignore self link
	// if e.Text == "" {
	// 	return
	// }

	// Print link
	// decodedSongTitle := decodeWindows1251([]byte(e.Text))
	// log.Printf("Song title found: %q\n", decodedSongTitle)

	// Visit link found on page
	// Only those links are visited which are in AllowedDomains
	// 	songCollector.Visit(e.Request.AbsoluteURL(link))
	// })

	// Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// })

	// Limit the maximum parallelism to cpu num
	songCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: runtime.NumCPU()})

	// Before making a request print "Visiting ..."
	// songCollector.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// })

	songCollector.OnHTML(`body`, func(e *colly.HTMLElement) {
		var song res.Song
		song.Link = e.Request.URL.String()

		title := e.ChildText("h2")
		decodedTitle := decodeWindows1251([]byte(title))
		song.Title = string(decodedTitle)
		// log.Println(string(decodedTitle))

		e.ForEach("p", func(i int, elem *colly.HTMLElement) {
			if i == 0 {
				decodedAuthor := decodeWindows1251([]byte(elem.Text))
				song.Author = string(decodedAuthor)
			}
			if i == 1 {
				decodedAlbum := decodeWindows1251([]byte(elem.Text))
				if strings.Contains(string(decodedAlbum), "Альбом") {
					song.Album = string(decodedAlbum)
				} else {

					// decodedVerses := decodeWindows1251([]byte(elem.Text))
					// log.Println(string(decodedVerses))
					// str := regexp.MustCompile(`<br/>`).Split(elem.Text, -1)
					// decodedVerses := decodeWindows1251([]byte(str))
					// log.Println(str)
					// for _, s := range []string{elem.Text} {
					// newLines := regexp.MustCompile(`<br/>`).ReplaceAllString(elem.Text, "\n")

					//dom := elem.DOM.Text()
					//log.Println(strings.Contains(dom, `<br/>`))
					//decodedVerses := decodeWindows1251([]byte(dom))
					//log.Println(string(decodedVerses))

					elem.DOM.Each(func(_ int, lyrics *goquery.Selection) {
						lyrics.Contents().Each(func(i int, lyric *goquery.Selection) {

							// handle spaces
							if !lyric.Is("br") {
								//trimmedLyric := strings.TrimSpace(lyric.Text())
								commaLyric := regexp.MustCompile(`&#39;`).ReplaceAllString(lyric.Text(), "'")
								nbspLyric := regexp.MustCompile(` `).ReplaceAllString(commaLyric, " ")
								decodedLyric := decodeWindows1251([]byte(nbspLyric))
								verse := res.Verse{
									Data: string(decodedLyric),
								}
								song.Verses = append(song.Verses, verse)
								log.Println(string(decodedLyric), "-", i)
							} else {
								// TODO: handle new lines?
							}
						})
					})

					// log.Println(dom.Children())
					// lines := strings.Split(elem.Text, `<br/>`)
					// brContains := strings.Contains(elem.Text, `<br/>`)
					// log.Println(brContains)
					// for _, l := range lines {
					// 	result := regexp.MustCompile(`&#39;`).ReplaceAllString(l, "'")

					// 	trimmedResult := regexp.MustCompile(" ").ReplaceAllString(result, " ")

					// 	decodedResult := decodeWindows1251([]byte(trimmedResult))

					// 	log.Println(string(decodedResult))
					// }
					// }

					// dirtyVerses := make([]string, 0)
					// for _, e := range elem.Text {
					// str := regexp.MustCompile(`<br/>`).Split(e, -1)
					// for _, s := range str {
					// result := regexp.MustCompile(`&#39;`).ReplaceAllString(s, "'")

					// &nbsp;
					// &#160;
					// &#xA0;
					// ⌥ Opt+Space
					// non suka breaking space replaced by human readble space
					// trimmedResult := regexp.MustCompile(" ").ReplaceAllString(result, " ")
					// decodedResult := decodeWindows1251([]byte(trimmedResult))
					// log.Println(string(decodedVerses))
					// dirtyVerses = append(dirtyVerses, string(decodedResult)+"\n")
					// }

					// dirtyVerses = append(dirtyVerses, "\n\n")
					// }
				}
			}
			if i == 2 {
				// lyrics
			}
		})

		songs = append(songs, song)
		// log.Println(len(songs))
	})

	// On every a element which has `div[id=headers]` attribute call callback
	// to fetch song title
	// songCollector.OnHTML(`div[id=headers]`, func(e *colly.HTMLElement) {
	// 	title := e.ChildText("h3")
	// 	decodedTitle := decodeWindows1251([]byte(title))
	// 	song.Title = string(decodedTitle)
	// 	// fmt.Println(e.Request.URL)
	// 	// log.Println(string(decodedTitle))
	// })

	// On every a element which has `div[id=cont]` attribute call callback
	// songCollector.OnHTML(`div[id=cont]`, func(e *colly.HTMLElement) {
	// 	// log.Println("Song link found", e.Request.URL)

	// 	song.Link = e.Request.URL.String()

	// 	// for each header element
	// 	e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
	// 		decodedSmth := decodeWindows1251([]byte(elem.Text))

	// 		// substring after Автор:
	// 		rau := myRexexp{regexp.MustCompile("(?:Автор:[\\s])(?P<author>.+)")}
	// 		song.Author = rau.findStringSubmatchMap(string(decodedSmth))["author"]

	// 		// substring after Альбом:
	// 		// -4
	// 		ral := myRexexp{regexp.MustCompile("(?:Альбом:[\\s])(?P<album>.+)")}
	// 		song.Album = ral.findStringSubmatchMap(string(decodedSmth))["album"]
	// 	})
	// 	// fmt.Println(e.Request.URL.String())

	// 	// Find body with lyrics
	// 	dirtyHTML, _ := e.DOM.Html()

	// 	// fixme
	// 	rl := regexp.MustCompile("(</script>)(.+)(<p>)")
	// 	lyricHTML := rl.FindString(dirtyHTML)

	// 	// fixme
	// 	ril := regexp.MustCompile(`<\/p><p><strong>.+<\/strong>.+<\/p>(?P<Lyrics>.+)<p>`)
	// 	improvedLyricsHTML := ril.FindAllStringSubmatch(lyricHTML, -1)
	// 	names := ril.SubexpNames()

	// 	// if non match patter return
	// 	if improvedLyricsHTML == nil {
	// 		return
	// 	}

	// 	songs = append(songs, song)
	// 	if len(songs) > 580 {
	// 		log.Println("songs:", len(songs))
	// 	}

	// 	// create map with group name -> content
	// 	md := map[string]string{}
	// 	for i, n := range improvedLyricsHTML[0] {
	// 		md[names[i]] = n
	// 	}

	// 	// split to verses group
	// 	rlp := regexp.MustCompile(`<br/><br/>`)
	// 	unparsedLyrics := rlp.Split(md["Lyrics"], -1)

	// 	// split to separated verses
	// 	dirtyVerses := make([]string, 0)
	// 	for _, e := range unparsedLyrics {
	// 		str := regexp.MustCompile(`<br/>`).Split(e, -1)
	// 		for _, s := range str {
	// 			result := regexp.MustCompile(`&#39;`).ReplaceAllString(s, "'")

	// 			// &nbsp;
	// 			// &#160;
	// 			// &#xA0;
	// 			// ⌥ Opt+Space
	// 			// non suka breaking space replaced by human readble space
	// 			trimmedResult := regexp.MustCompile(" ").ReplaceAllString(result, " ")
	// 			decodedResult := decodeWindows1251([]byte(trimmedResult))
	// 			dirtyVerses = append(dirtyVerses, string(decodedResult)+"\n")
	// 		}

	// 		dirtyVerses = append(dirtyVerses, "\n\n")
	// 	}

	// 	verses := make([]res.Verse, 0)

	// 	for i, v := range dirtyVerses {
	// 		verses = append(verses, res.Verse{Ordinal: i, Data: v})
	// 	}

	// 	song.Verses = verses
	// 	// songs = append(songs, song)

	// 	// log.Println("songs:", len(songs))
	// })

	// Start scraping on http://www.gr-oborona.ru/texts/
	//c.Visit(textsPage)

	// Wait until threads are finished
	//c.Wait()
	songCollector.Wait()

	return songs
}

// decode shitty cp1251 to human readalbe utf-8
func decodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}
