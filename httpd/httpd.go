package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "p", 8800, "listen port")

	flag.Parse()
}

func main() {
	http.HandleFunc("/", uiStatusHandler)
	http.HandleFunc("/ui/index.html", uiStatusHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func uiStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
