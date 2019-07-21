package scrapper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	domain "github.com/solairerove/linden-honey-go-scrapper/domain"
	"golang.org/x/text/encoding/charmap"
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
func ScrapLetov(db *sql.DB) {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domain: www.gr-oborona.ru
		colly.AllowedDomains("www.gr-oborona.ru"),

		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(1),

		// Visit only root url and urls which start with "text" on www.gr-oborona.ru
		colly.URLFilters(
			regexp.MustCompile("http://www.gr-oborona.ru/texts/"),
		),
	)

	songCollector := c.Clone()

	var song domain.Song

	// On every a element which has href attribute call callback
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// ignore self link
		if e.Text == "" {
			return
		}

		// Print link
		decodedSongTitle := decodeWindows1251([]byte(e.Text))
		log.Printf("Song title found: %q\n", decodedSongTitle)

		song.Title = string(decodedSongTitle)

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		songCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// On every a element which has `div[id=cont]` attribute call callback
	songCollector.OnHTML(`div[id=cont]`, func(e *colly.HTMLElement) {
		log.Println("Song link found", e.Request.URL)

		song.Link = e.Request.URL.String()

		e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
			decodedSmth := decodeWindows1251([]byte(elem.Text))
			log.Printf("Find smth from loop %s", decodedSmth)

			// substring after Автор:
			if strings.Contains(string(decodedSmth), "Автор") {
				rau := myRexexp{regexp.MustCompile("(?:Автор:[\\s])(?P<author>.+)")}
				song.Author = rau.findStringSubmatchMap(string(decodedSmth))["author"]
			}

			// substring after Альбом:
			if strings.Contains(string(decodedSmth), "Альбом") {
				ral := myRexexp{regexp.MustCompile("(?:Альбом:[\\s])(?P<album>.+)")}
				song.Album = ral.findStringSubmatchMap(string(decodedSmth))["album"]
			}
		})

		// Find body with lyrics
		dirtyHTML, _ := e.DOM.Html()

		// fixme
		rl := regexp.MustCompile("(</script>)(.+)(<p>)")
		lyricHTML := rl.FindString(dirtyHTML)

		// fixme
		ril := regexp.MustCompile(`<\/p><p><strong>.+<\/strong>.+<\/p>(?P<Lyrics>.+)<p>`)
		improvedLyricsHTML := ril.FindAllStringSubmatch(lyricHTML, -1)
		names := ril.SubexpNames()

		// if non match patter return
		if improvedLyricsHTML == nil {
			return
		}

		// create map with group name -> content
		md := map[string]string{}
		for i, n := range improvedLyricsHTML[0] {
			md[names[i]] = n
		}

		// split to verses group
		rlp := regexp.MustCompile(`<br/><br/>`)
		unparsedLyrics := rlp.Split(md["Lyrics"], -1)

		// split to separated verses
		dirtyVerses := make([]string, 0)
		for _, e := range unparsedLyrics {
			log.Print("\n")
			str := regexp.MustCompile(`<br/>`).Split(e, -1)
			for _, s := range str {
				result := regexp.MustCompile(`&#39;`).ReplaceAllString(s, "'")

				// &nbsp;
				// &#160;
				// &#xA0;
				// ⌥ Opt+Space
				// non suka breaking space replaced by human readble space
				trimmedResult := regexp.MustCompile(" ").ReplaceAllString(result, " ")
				decodedResult := decodeWindows1251([]byte(trimmedResult))
				dirtyVerses = append(dirtyVerses, string(decodedResult)+"\n")

				log.Printf("Lyrics found %s", string(decodedResult))
			}

			dirtyVerses = append(dirtyVerses, "\n\n")
		}

		verses := make([]domain.Verse, 0)

		for i, v := range dirtyVerses {
			verses = append(verses, domain.Verse{Ordinal: i, Data: v})
		}

		song.Verses = verses

		marshaledSong, _ := json.Marshal(song)
		log.Printf("Prepare to save next Song -> %s", string(marshaledSong))

		song.SaveSong(db)
	})

	// fixme
	c.Visit("http://www.gr-oborona.ru/texts/")
}

// decode shitty cp1251 to human readalbe utf-8
func decodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}
