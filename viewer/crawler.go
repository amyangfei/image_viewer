// Simple crawler

package viewer

import (
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type DataType uint32

const (
	Image DataType = iota
	Href
)

type CrawData struct {
	Name string
	Url  string
	Type DataType
	Data []byte
}

func getHtmlData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("get url with error")
		return nil, err
	}
	defer resp.Body.Close()
	body := resp.Body
	data, err := ioutil.ReadAll(body)
	return data, err
}

var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["'](https?://[^"']+)["']`)
var hrefRE = regexp.MustCompile(`<a[^>]+\bhref=["'](https?://[^"']+)["']`)

func findElems(htm string, re *regexp.Regexp) []string {
	elems := re.FindAllStringSubmatch(htm, -1)
	out := make([]string, len(elems))
	for i := range out {
		out[i] = elems[i][1]
	}
	return out
}

func findImages(htm string) []string {
	return findElems(htm, imgRE)
}

func findLinks(htm string) []string {
	return findElems(htm, hrefRE)
}

func crawSublink(htm string, c chan<- CrawData) {
	var wg sync.WaitGroup
	for _, link := range findLinks(htm) {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			u, err := url.Parse(src)
			if err != nil {
				log.Printf("invalid url path: %s", src)
				return
			}
			// we cannot use / in a filename
			name := strings.Replace(
				strings.TrimRight(
					strings.TrimPrefix(src, u.Scheme+"://"),
					"/"),
				"/", "_", -1)
			c <- CrawData{name, src, Href, nil}
		}(link)
	}
	wg.Wait()
	close(c)
}

func crawlImg(htm string, c chan<- CrawData) {
	var wg sync.WaitGroup
	for _, imgUrl := range findImages(htm) {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			u, err := url.Parse(src)
			if err != nil {
				log.Printf("invalid url path: %s", src)
				return
			}
			subPath := strings.Split(u.Path, "/")
			filename := subPath[len(subPath)-1]
			if len(filename) == 0 {
				uid, err := uuid.NewV4()
				if err != nil {
					log.Printf("generage image name error: %s", err)
					return
				}
				filename = uid.String()
			}

			resp, err := http.Get(src)
			if err != nil {
				log.Printf("fetch url with error: %s", err)
				return
			}
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("read url data with error: %s", err)
				return
			}
			c <- CrawData{filename, imgUrl, Image, raw}
		}(imgUrl)
	}
	wg.Wait()
	close(c)
}

func Crawl(link string) ([]CrawData, error) {
	data, err := getHtmlData(link)
	if err != nil {
		return nil, err
	}
	result := make([]CrawData, 0)
	imgCh := make(chan CrawData)
	urlCh := make(chan CrawData)
	html := string(data)
	var wg sync.WaitGroup
	wg.Add(2)
	go crawlImg(html, imgCh)
	go crawSublink(html, urlCh)
	go func() {
		defer wg.Done()
		for value := range imgCh {
			result = append(result, value)
		}
	}()
	go func() {
		defer wg.Done()
		for value := range urlCh {
			result = append(result, value)
		}
	}()
	wg.Wait()
	return result, nil
}