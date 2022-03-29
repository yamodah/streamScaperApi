package main

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"os"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)
 
 type StreamLinks struct{
	 EventName string `json:"event"`
	 Streams []string `json:"streams"`
 }
//Big function will work as a switch board housing unique stream scrapers to keep main function clean
 func getStreams(w http.ResponseWriter, r *http.Request){

	c:=colly.NewCollector()
	desiredEvent:= r.URL.Query().Get("stream")

	 switch r.URL.Path {
	 case "/":
		 fmt.Fprint(w," Welcome to Yamodah's Streams ! \n to receive links for a stream, query in url by adding '/all?stream=' followed by the name of a live event. \n for example looking for a fight enter a fighters last name, a football match enter a team name for the query. \n this is a work in progress but please enjoy!")
	 case "/all":
		rojaScrape(w,c,desiredEvent)
	// 	liveTvScrape(w)
	// 	mamaHDScrape(w)
	// 	streamEastScrape(w)
	//  case "/roja":
	// 	rojaScrape(w,c,desiredEvent)
	//  case "/livetv":
	// 	liveTvScrape(w)
	//  case "/mamahd":
	// 	mamaHDScrape(w)
	//  case "/east":
	// 	streamEastScrape(w)
	 default:
		fmt.Fprint(w,"Big fat Error")
	 }
 }
 func rojaScrape(w http.ResponseWriter, c *colly.Collector, desiredEvent string){

	 fmt.Println("roja scraping ...")
	 
	 c.OnHTML("#agendadiv span.list", func(h *colly.HTMLElement){
		 //full table
		 selection := h.DOM
		 //individual events
		 childNodes:= selection.Children().Nodes
		 for class:=0;class<34;class++{

			 titles:=selection.FindNodes(childNodes[class]).Find("div.menutitle").Children().Nodes
			 teamNames:=strings.ToLower(selection.FindNodes(titles...).Find("b span").Text())
			 if strings.Contains(teamNames,desiredEvent) {
				 table:=selection.FindNodes(childNodes[class]).Find("tbody").Children().Nodes
				 rows:=selection.FindNodes(table...).Children().Nodes
				 links:=selection.FindNodes(rows...).Find("td a")
				 var streams []string
				 links.Each(func(i int, s *goquery.Selection) {
					 link,_ := s.Attr("href")
					 fmt.Printf("game: %s \nlink: %s\n", teamNames, link)
					 streams = append(streams, link)
				 })
				streamPack := StreamLinks{teamNames,streams}
				json.NewEncoder(w).Encode(streamPack)
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
//  func liveTvScrape(w http.ResponseWriter){
// 	fmt.Fprint(w, "LiveTV \n")
//  }
//  func mamaHDScrape(w http.ResponseWriter){
// 	fmt.Fprint(w, "MAMAHD \n")
//  }
//  func streamEastScrape(w http.ResponseWriter){
// 	fmt.Fprint(w, "StreamEast \n")
//  }

 func main(){
	 port:=os.Getenv("PORT")
	 http.HandleFunc("/", getStreams)
	 log.Fatal(http.ListenAndServe(":"+port, nil))
 }