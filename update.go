package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	q "github.com/PuerkitoBio/goquery"
)

const startingURL = "https://www.iskultur.com.tr/kitap/modern-klasikler?orderby=title&order=ASC"

type book struct {
	img       string
	name      string
	url       string
	author    string
	authorURL string
	len       int
}

func (b book) String() string {
	return fmt.Sprintf("* [ ] <img src=\"%s\" width=\"60\" height=\"98\"> [%s](%s) - [%s](%s) - %d\n",
		b.img, b.name, b.url, b.author, b.authorURL, b.len)
}

func fillDetails(b *book) {
	doc, err := q.NewDocument(b.url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".tabDiv tr").Each(func(i int, s *q.Selection) {
		tds := s.Find("td")
		if strings.TrimSpace(tds.First().Text()) == "Sayfa Sayısı" {
			var err error
			if b.len, err = strconv.Atoi(tds.Last().Text()); err != nil {
				log.Fatal("Couldn't parse len for ", b.url, err)
			}
		}
	})
}

func parsePage(doc *q.Document) ([]book, string) {
	var books []book
	doc.Find(".productList").Each(func(i int, s *q.Selection) {
		books = append(books, book{
			url:       s.Find("a").First().AttrOr("href", ""),
			name:      s.Find("a").First().AttrOr("title", ""),
			img:       s.Find(".resIMG").First().AttrOr("src", ""),
			author:    s.Find("a.text3").Text(),
			authorURL: s.Find("a.text3").AttrOr("href", ""),
		})
	})
	for i := range books {
		fillDetails(&books[i])
	}
	return books, doc.Find(".paging .emm-next").First().AttrOr("href", "")
}

func getBooks(url string) []book {
	log.Println("Parsing ", url)
	doc, err := q.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	books, next := parsePage(doc)
	log.Printf("Parsed %d\n", len(books))
	if next != "" {
		books = append(books, getBooks(next)...)
	}
	return books
}

func main() {
	f, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("## [Turkiye Is Bankasi Yayinlari - Modern Klasikler Serisi](%s)", startingURL))
	f.WriteString("\n\n")
	for _, b := range getBooks(startingURL) {
		f.WriteString(b.String())
	}
}
