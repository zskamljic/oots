package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zskamljic/oots"
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

	tl := oots.NewThreadLimiter()
	for i := 1; i <= max; i++ {
		tl.WaitTurn()
		tl.Add(1)
		go fetchImage(i, tl)
	}

	tl.Wait()
	log.Println("Done")
}

func fetchImage(index int, tl *oots.ThreadLimiter) {
	defer tl.Done()

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
