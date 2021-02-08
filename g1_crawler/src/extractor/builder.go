package extractor

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/crawler/src/hash"
	"github.com/gocolly/colly"
	"io"
	"log"
	"strings"
	"time"
)

type news struct {
		Id 		 string `json:"id"`
		Title    string `json:"title"`
		Author   string `json:"author"`
		Date     string `json:"date"`
		Escraped string `json:"escraped"`
		Page     string `json:"page"`
		Body     string `json:"body"`
		Url      string `json:"url"`
}

type domains struct {
	Domains []string `json:"https://g1.globo.com"`
}

func Extract() string {
	var encoder *json.Encoder
	page := ""
	domains := ReadFromBucket("globo_g1", "org_site/domains.json")
	dataBuffer := new(bytes.Buffer)
	builder := colly.NewCollector(
		colly.UserAgent("Babozavuti"),
		colly.AllowedDomains("g1.globo.com"),
		colly.Async(true),
		colly.MaxDepth(12),
	)

	builder.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 8,
		RandomDelay: 5 * time.Second})

	extractor := builder.Clone()
	dataNews := make([]news, 0, 10)

	builder.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Request.AbsoluteURL(e.Attr("href"))
			if strings.Index(link, "noticia") != -1 {
				extractor.Visit(link)
			}

	})

	extractor.OnHTML("body", func(e *colly.HTMLElement) {
		currentTime := time.Now()
		timeFormat := currentTime.Format("2 Jan 2006 15:04:05")
		news := news{
			Id:       hash.HashUrl(e.Request.AbsoluteURL(e.Attr("href"))),
			Title:    e.ChildText(".content-head__title"),
			Author:   e.ChildText(".content-publication-data__from"),
			Date:     e.ChildText(".content-publication-data__updated"),
			Escraped: timeFormat,
			Page:     page,
			Body:     e.ChildText(".content-text__container"),
			Url:      e.Request.AbsoluteURL(e.Attr("href")),
		}

		dataNews = append(dataNews, news)

	})


	builder.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Failed to visit url %s failed with response %d",r.Request.URL, r.StatusCode)
		counter := 0
		defer recoverFromFailToVisit()
		for counter < 3 {
			fmt.Println("\nretrying url "+r.Request.URL.String())
			builder.Visit(r.Request.URL.String())
			counter += 1
		}
		fmt.Println("Reached max retry attempts, giving up")
	})

	builder.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	builder.OnResponse(func(r *colly.Response) {
		log.Printf("Done visiting %s", r.Request.URL.String())
		builder.AllowURLRevisit = false
	})

	for _, domains := range domains {
		for i := 1; i < 11; i++ {
			page = fmt.Sprintf("https://g1.globo.com/%s/index/feed/pagina-%d.ghtml", domains, i)
			err := builder.Visit(page)
			if err != nil {
				fmt.Println(err)
				}
				builder.Wait()
				extractor.Wait()
			}
		}

	result := make(map[string]news)

	for _, entry := range dataNews {
		result[entry.Id] = news{
			Title: entry.Title,
			Author: entry.Author,
			Date: entry.Date,
			Escraped: entry.Escraped,
			Page: entry.Page,
			Body: entry.Body,
			Url: entry.Url,
		}
	}

	encoder = json.NewEncoder(dataBuffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(result)
	if err != nil {
		panic(err)
	}

	return strings.Trim(dataBuffer.String(), "\n") + ","

}

func ReadFromBucket(bucketName string, jsonFile string) []string {
	var jsonReader bytes.Buffer
	decoder := domains{}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Panic(err)
	}

	bucket := client.Bucket(bucketName)

	obj := bucket.Object(jsonFile)
	r, err := obj.NewReader(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer r.Close()

	if _, err := io.Copy(&jsonReader, r); err != nil {
		if err != nil {
			log.Panic(err)
		}
	}

	err = json.Unmarshal(jsonReader.Bytes(), &decoder)

	return decoder.Domains
}

func recoverFromFailToVisit() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

