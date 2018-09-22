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

func findImages(htm string) []string {
	imgs := imgRE.FindAllStringSubmatch(htm, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
	}
	return out
}

func Crawl(link string) ([]CrawData, error) {
	data, err := getHtmlData(link)
	if err != nil {
		return nil, err
	}
	result := make([]CrawData, 0)
	var wg sync.WaitGroup
	for _, imgUrl := range findImages(string(data)) {
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
			result = append(result, CrawData{filename, Image, raw})
		}(imgUrl)
	}
	wg.Wait()
	return result, nil
}

/*
func main() {
	link := "http://www.lofter.com/tag/365日摄影计划"
	raw := Crawl(link)
	for _, data := range raw {
		log.Printf("name: %s, content length: %d", data.Name, len(data.Data))
	}
}
*/
