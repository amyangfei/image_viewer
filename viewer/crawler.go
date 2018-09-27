// Simple crawler

package viewer

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/teris-io/shortid"
)

type DataType uint32

const (
	Image DataType = iota
	Href
)

type ImageInfo struct {
	Src   string
	Class string
	Alt   string
}

type CrawData struct {
	Name string
	Url  string
	Type DataType
	Data []byte
}

// getHtmlData visits url and returns page source, if headless is true,
// javascript will also be executed
func getHtmlData(url string, headless bool, driver selenium.WebDriver) ([]byte, error) {
	if headless {
		err := driver.Get(url)
		if err != nil {
			log.Printf("headless get url with error: %s", err)
			return nil, err
		}
		data, err := driver.PageSource()
		return []byte(data), err
	} else {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get url with error: %s\n", err)
			return nil, err
		}
		defer resp.Body.Close()
		body := resp.Body
		data, err := ioutil.ReadAll(body)
		return data, err
	}
}

var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"'><]*?)["']`)
var hrefRE = regexp.MustCompile(`<a[^>]+\bhref=["']([^"'><]*?)["']`)

func findElems(htm string, re *regexp.Regexp) []string {
	elems := re.FindAllStringSubmatch(htm, -1)
	out := make([]string, len(elems))
	for i := range out {
		out[i] = elems[i][1]
	}
	return out
}

func findLinks(htm string) []string {
	return findElems(htm, hrefRE)
}

func findImages(htm string) []*ImageInfo {
	urls := findElems(htm, imgRE)
	result := make([]*ImageInfo, 0)
	for _, u := range urls {
		result = append(result, &ImageInfo{Src: u, Class: "", Alt: ""})
	}
	return result
}

func findLinks2(htm string) []string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htm)))
	if err != nil {
		log.Printf("go query parse error: %s", err)
		return findLinks(htm)
	}
	result := make([]string, 0)
	doc.Find("html a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			result = append(result, href)
		}
	})
	return result
}

func findImages2(htm string) []*ImageInfo {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htm)))
	if err != nil {
		log.Printf("go query parse error: %s", err)
		return findImages(htm)
	}
	result := make([]*ImageInfo, 0)
	doc.Find("html img").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		alt, _ := s.Attr("alt")
		src, exists := s.Attr("src")
		if exists {
			result = append(result, &ImageInfo{Src: src, Class: class, Alt: alt})
		}
	})
	return result
}

func crawSublink(baseUrl string, htm string, c chan<- CrawData, notifyWG *sync.WaitGroup) {
	var wg sync.WaitGroup
	baseU, _ := url.Parse(baseUrl)
	for _, link := range findLinks2(htm) {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			u, err := url.Parse(src)
			if err != nil {
				log.Printf("invalid url path: %s", src)
				return
			}

			// igore none url such as "javascript:void(0)"
			if u.Scheme == "javascript" {
				return
			}

			// if href field is absolute path, join location host
			if u.Scheme == "" {
				src = baseU.Scheme + "://" + baseU.Host + "/" + u.Path
				u.Scheme = baseU.Scheme
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
	notifyWG.Done()
}

func crawlImg(baseUrl string, htm string, c chan<- CrawData, notifyWG *sync.WaitGroup) {
	sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
	var wg sync.WaitGroup
	baseU, _ := url.Parse(baseUrl)
	for _, imgInfo := range findImages2(htm) {
		wg.Add(1)
		go func(info *ImageInfo) {
			defer wg.Done()
			src := info.Src
			u, err := url.Parse(src)
			if err != nil {
				log.Printf("invalid url path: %s", src)
				return
			}

			// igore none url such as "javascript:void(0)"
			if u.Scheme == "javascript" {
				return
			}

			// ignore base64 image
			if u.Scheme == "data" && strings.HasPrefix(src, "data:image") {
				i := strings.Index(src, ",")
				if i < 0 {
					log.Printf("invalid base64 image\n")
				}
				reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(src[i+1:]))
				buffer := bytes.Buffer{}
				_, err := buffer.ReadFrom(reader)
				if err != nil {
					log.Printf("read from base64 buffer error: %s", err)
					return
				}
				fm, err := DetectImageType(buffer.Bytes())
				if err != nil {
					log.Printf("read image config error: %s", err)
					return
				}
				fid := RandomId(sid)
				filename := info.Class + fid + "." + fm
				c <- CrawData{filename, src, Image, buffer.Bytes()}
				return
			}

			// If href field is absolute path, join location host.
			// Img src without "/" prefix can redirect too, so we don't check
			// prefix here.
			if u.Scheme == "" {
				src = baseU.Scheme + "://" + baseU.Host + "/" + u.Path
				u.Scheme = baseU.Scheme
			}

			filename := ""
			needExpandExt := false
			if info.Alt != "" {
				// Get filename from alt information
				fid := RandomId(sid)
				filename = info.Alt + fid
				needExpandExt = true
			} else {
				// Get filename from last path field
				subPath := strings.Split(u.Path, "/")
				filename = subPath[len(subPath)-1]
				if len(filename) == 0 {
					filename = RandomId(sid)
				}
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

			// Complete file extension if needed
			if needExpandExt {
				ext, err := DetectImageType(raw)
				if err == nil {
					filename = filename + "." + ext
				}
			}

			c <- CrawData{filename, src, Image, raw}
		}(imgInfo)
	}
	wg.Wait()
	notifyWG.Done()
}

func Crawl(link string, headless bool, driver selenium.WebDriver) ([]CrawData, error) {
	data, err := getHtmlData(link, headless, driver)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	html := string(data)
	result := make([]CrawData, 0)
	resultCh := make(chan CrawData)
	wg.Add(2)
	go crawlImg(link, html, resultCh, &wg)
	go crawSublink(link, html, resultCh, &wg)
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	for value := range resultCh {
		result = append(result, value)
	}
	return result, nil
}
