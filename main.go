package main

import (
	"fmt"
	"net/http"
	"github.com/gocolly/colly"
)
//big function will work as a switch board housing unique stream scrapers to keep main function clean
 func getStreams(w http.ResponseWriter, r *http.Request){
	 switch r.URL.Path {
	 case "/":
		 fmt.Fprint(w,"Home page")
	 case "/all":
		rojaScrape(w)
		liveTvScrape(w)
		mamaHDScrape(w)
		streamEastScrape(w)
	 case "/roja":
		rojaScrape(w)
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
 func rojaScrape(w http.ResponseWriter){
	 fmt.Fprint(w, "roja \n")
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