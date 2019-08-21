package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "p", 8800, "listen port")

	flag.Parse()
}

func main() {
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			logrus.Fatalf("error %v", err)
		}
	}()

	http.HandleFunc("/say", sayHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/ui/index.html", uiStatusHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "my own error message", http.StatusInternalServerError)
}

func sayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, say I love you!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "Hi there, I love you!")
}

func uiStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
