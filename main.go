package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

//big function will work as a switch board housing unique stream scrapers to keep main function clean
 func getStreams(w http.ResponseWriter, r *http.Request){

	c:=colly.NewCollector()

	 switch r.URL.Path {
	 case "/":
		 fmt.Fprint(w,"Home page")
	 case "/all":
		rojaScrape(w,c)
		liveTvScrape(w)
		mamaHDScrape(w)
		streamEastScrape(w)
	 case "/roja":
		rojaScrape(w,c)
	 case "/livetv":
		liveTvScrape(w)
	 case "/mamahd":
		mamaHDScrape(w)
	 case "/east":
		streamEastScrape(w)
	 default:
		fmt.Fprint(w,"Big fat Error")
	 }
 }
 func rojaScrape(w http.ResponseWriter, c *colly.Collector){
	 fmt.Println("roja scraping ...")

	 c.OnHTML("#agendadiv span.list", func(h *colly.HTMLElement){
		 //full table
		 selection := h.DOM
		 //individual events
		 childNodes:= selection.Children().Nodes
		 for class:=0;class<34;class++{

			 titles:=selection.FindNodes(childNodes[class]).Find("div.menutitle").Children().Nodes
			 teamNames:=strings.ToLower(selection.FindNodes(titles...).Find("b span").Text())
			 if strings.Contains(teamNames,"thunder") {
				 table:=selection.FindNodes(childNodes[class]).Find("tbody").Children().Nodes
				 rows:=selection.FindNodes(table...).Children().Nodes
				 links:=selection.FindNodes(rows...).Find("td a")
				 links.Each(func(i int, s *goquery.Selection) {
					 link,_ := s.Attr("href")
					 fmt.Println(link)
				 })
				 break
				}
		 }
	 })
	 c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	 c.Visit("http://www.rojadirecta.me")
	 fmt.Println("*** scraping complete ***")
 }
 func liveTvScrape(w http.ResponseWriter){
	fmt.Fprint(w, "LiveTV \n")
 }
 func mamaHDScrape(w http.ResponseWriter){
	fmt.Fprint(w, "MAMAHD \n")
 }
 func streamEastScrape(w http.ResponseWriter){
	fmt.Fprint(w, "StreamEast \n")
 }

 func main(){
	
	 http.HandleFunc("/", getStreams)
	 http.ListenAndServe(":5000", nil)
 }