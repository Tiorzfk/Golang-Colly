package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type ListLiga struct {
	ID        string `json:"id,omitempty"`
	Nama      string `json:"nama,omitempty"`
	Link      string `json:"link,omitempty"`
	LinkLiga  string `json:"link_liga,omitempty"`
}

type ListBerita struct {
	ID        string `json:"id,omitempty"`
	Judul     string `json:"judul,omitempty"`
	Link      string `json:"link,omitempty"`
	Foto      string `json:"foto,omitempty"`
	Slug      string `json:"slug,omitempty"`
}

type DetailBeritaModel struct {
	ID        string `json:"id,omitempty"`
	Isi     string `json:"isi,omitempty"`
	Foto      string `json:"foto,omitempty"`
}

var liga []ListLiga
var berita []ListBerita
var detail_berita []DetailBeritaModel

func Liga(d *gin.Context) {
	c := colly.NewCollector()

	c.OnHTML("a[class=widget-competitions-popular__competition]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Println(e)
		s := strings.Split(link,"/")

		liga = append(liga, ListLiga{ID:s[len(s)-1],LinkLiga:s[len(s)-2],Nama:e.Text,Link:link})

		// fmt.Println(s[len(s)-1])
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://www.goal.com/id/kompetisi")

	d.JSON(200, liga)
	
}

func Berita(d *gin.Context) {
	id := d.Param("id")
	c := colly.NewCollector()
	berita = berita[:0]
	c.OnHTML("a[data-side=front][itemprop=url][data-sponsorship-slot-id=front] ", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.ChildText(".title-wrapper")
		foto := e.ChildText("noscript")
		linkSplit := strings.Split(link,"/")

		listberita := ListBerita{
			ID: linkSplit[len(linkSplit) -1],
			Judul:title,
			Link:link,
			Foto:foto,
			Slug: linkSplit[len(linkSplit) -2],
		}

		berita = append(berita, listberita)
		
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	c.Visit("http://www.goal.com/id/liga/1/" + id)

	d.JSON(200, berita)
}

func DetailBerita(d *gin.Context) {
	id := d.Param("id")
	c := colly.NewCollector()

	detail_berita = detail_berita[:0]
	c.OnHTML("article", func(e *colly.HTMLElement) {
		isi := e.ChildText(".body")
		foto := e.ChildText("noscript")

		detailberita := DetailBeritaModel{
			ID : id,
			Isi : isi,
			Foto : foto,
		}

		detail_berita = append(detail_berita, detailberita)
		
	})

	c.Visit("http://www.goal.com/id/berita/detail/" + id)

	d.JSON(200, detail_berita[0])
}

func main() {

	r := gin.Default()

	v1 := r.Group("/v1")

	v1.GET("/liga", Liga)
	v1.GET("/berita/:id", Berita)
	v1.GET("/detail/:id", DetailBerita)
	
	r.Run(":3001")

}