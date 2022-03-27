package main

import (
	"fmt"
	"net/http"
)
 func getStreams(w http.ResponseWriter, r *http.Request){
	 fmt.Fprint(w,"hakuna matata")
 }
 func main(){
	 http.HandleFunc("/", getStreams)
	 http.ListenAndServe(":5000", nil)
 }