package main

import (
	"fmt"
	"net/http"
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
			 //full event title
			 value:=selection.FindNodes(childNodes[class]).Find("div.menutitle").Children().Nodes
			 teamNames:=selection.FindNodes(value...).Find("b span").Text()
			 fmt.Printf("class:%d text: %s \n", class, teamNames)
		 }
	 })
	 c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	 c.Visit("http://www.rojadirecta.me")
	 fmt.Println("pizza pizza")
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