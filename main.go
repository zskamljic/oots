package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"log"

	"strconv"

	"sync"

	"net/http"

	"runtime"

	"github.com/PuerkitoBio/goquery"
)

// URL is the format for all images
const (
	BASE = "http://www.giantitp.com"
	URL  = BASE + "/comics/oots%04d.html"
)

func main() {
	doc, err := goquery.NewDocument("http://www.giantitp.com/comics/oots.html")
	if err != nil {
		log.Fatal(err)
	}

	sel := doc.Find("p.ComicList").First()
	sel.Children().Remove()

	maxStr, err := sel.Html()
	if err != nil {
		log.Fatal(err)
	}

	max, err := strconv.Atoi(strings.TrimSpace(maxStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("There's a total of ", max, "comics")

	os.Mkdir("comics", 0755)

	var wg sync.WaitGroup
	for i := 1; i <= max; i++ {
		if i%runtime.NumCPU() == 0 {
			wg.Wait()
		}
		wg.Add(1)
		go fetchImage(i, &wg)
	}

	wg.Wait()
	log.Println("Done")
}

func fetchImage(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := goquery.NewDocument(fmt.Sprintf(URL, index))
	if err != nil {
		log.Println(err)
		return
	}

	link, _ := doc.Find("td > img").Eq(3).Attr("src")

	dotIndex := strings.LastIndex(link, ".")

	ext := link[dotIndex:]
	img, _ := os.Create(fmt.Sprintf("comics/%04d%s", index, ext))
	defer img.Close()

	resp, _ := http.Get(BASE + link)
	defer resp.Body.Close()

	io.Copy(img, resp.Body)

	log.Printf("Done with %d\n", index)
}
