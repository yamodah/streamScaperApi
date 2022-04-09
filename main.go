package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"net/http"
	"encoding/json"
	"html/template"

    "github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)
 
 type StreamLinks struct{
	 StreamingSite string `json:"site"`
	 Streams []string `json:"streams"`
 }
 type StreamsBundle struct{
	 EventName string `json:"event"`
	 Links []StreamLinks `json:"links"`
 }

//Big function will work as a switch board housing unique stream scrapers to keep main function clean
func getStreams(w http.ResponseWriter, r *http.Request){

	c:=colly.NewCollector()
	desiredEvent:= r.URL.Query().Get("stream")

		StreamsBundle:= StreamsBundle{EventName: desiredEvent}
		liveTvScrape(w,c,desiredEvent, &StreamsBundle)
		liveStream2Watch(w,c,desiredEvent, &StreamsBundle)
		rojaScrape(w,c,desiredEvent, &StreamsBundle)
		json.NewEncoder(w).Encode(StreamsBundle)

 }

 func homePage(w http.ResponseWriter, r *http.Request){
	 var filename = "search.html"
	 t,err := template.ParseFiles(filename)
	 if err != nil {
		 log.Fatal("Error: when parsing file",err)
	 }
	 err = t.ExecuteTemplate(w, filename, nil)
	 if err != nil {
		log.Fatal("Error: when executing file",err)
	}
 }

 func rojaScrape(w http.ResponseWriter, c *colly.Collector, desiredEvent string, Streamsbundle *StreamsBundle){

	//  fmt.Printf("roja scraping for %s ... \n", desiredEvent)
	 
	 c.OnHTML("#agendadiv span.list", func(h *colly.HTMLElement){
		 //full table
		 selection := h.DOM
		 //individual events
		 childNodes:= selection.Children().Nodes
		 for class:=0;class<len(childNodes)-1;class++{

			 titles:=selection.FindNodes(childNodes[class]).Find("div.menutitle").Children().Nodes
			 teamNames:=strings.ToLower(selection.FindNodes(titles...).Find("b span").Text())
			 if strings.Contains(teamNames,strings.ToLower(desiredEvent)) {
				 table:=selection.FindNodes(childNodes[class]).Find("tbody").Children().Nodes
				 rows:=selection.FindNodes(table...).Children().Nodes
				 links:=selection.FindNodes(rows...).Find("td a")
				 var streams []string
				 links.Each(func(i int, s *goquery.Selection) {
					 link,_ := s.Attr("href")
					 streams = append(streams, link)
				 })
				streamPack := StreamLinks{"ROJA",streams}
				Streamsbundle.EventName = teamNames
				Streamsbundle.Links = append(Streamsbundle.Links,streamPack)
				break
				}
		 }
	 })
	 c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	 c.Visit("http://www.rojadirecta.me")
	//  fmt.Println("... scraping complete")
 }
 func liveTvScrape(w http.ResponseWriter, c *colly.Collector, desiredEvent string, Streamsbundle *StreamsBundle){
	 liveTvURL := "http://livetv.ru/enx/megasearch/?msq="+desiredEvent
	 c.OnXML("/html/body/table/tbody/tr/td[2]/table/tbody/tr[3]/td/table/tbody/tr/td[2]/table/tbody/tr/td/table/tbody/tr[2]/td/table/tbody/tr/td/table/tbody/tr/td/table[5]/tbody/tr/td[2]/a", func(h *colly.XMLElement){
		 link:= h.Attr("href")
		 link = strings.ReplaceAll(link,"__/","")
		 var liveLink[]string
		 liveLink=append(liveLink, "http://livetv.ru"+link)
		 streamPack:= StreamLinks{"liveTV", liveLink}
		 Streamsbundle.Links = append(Streamsbundle.Links,streamPack)
		 
	 })
	 c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Body, "\nError:", err)
	})
	// fmt.Fprint(w, "LiveTV \n")
	c.Visit(liveTvURL)
 }
func liveStream2Watch(w http.ResponseWriter, c *colly.Collector, desiredEvent string, Streamsbundle *StreamsBundle){
	liveStreamURL := "https://live.xn--tream2watch-i9d.com/search?q="+desiredEvent
	c.OnHTML("div.main div.main-inner div.layouts-page-content div.layouts-search-single-item._rows._stm.stream-box div.item-body a", func(h *colly.HTMLElement) {

		link:= h.Attr("href")
		 link = strings.ReplaceAll(link,"__/","")
		 var liveLink[]string
		 liveLink=append(liveLink,link)
		 streamPack:= StreamLinks{"liveStream", liveLink}
		 Streamsbundle.Links = append(Streamsbundle.Links,streamPack)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Body, "\nError:", err)
	})
	c.Visit(liveStreamURL)
}

 func main(){
	 port:=os.Getenv("PORT")
	 r:=mux.NewRouter().StrictSlash(true)
	 r.HandleFunc("/",homePage)
	 r.HandleFunc("/all", getStreams)
	//  credentials := handlers.AllowCredentials()
    //  methods := handlers.AllowedMethods([]string{"GET"})
    //  origins := handlers.AllowedOrigins([]string{"*"})
     log.Fatal(http.ListenAndServe(":"+port, handlers.CORS()(r)))

 }