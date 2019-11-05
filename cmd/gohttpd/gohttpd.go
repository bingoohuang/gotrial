package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

var (
	port  int
	pport int
	print bool
)

func init() {
	flag.IntVar(&port, "p", 8800, "listen port")
	flag.IntVar(&port, "pp", 0, "pprof port")
	flag.BoolVar(&print, "print", false, "print sth")

	flag.Parse()
}

func main() {
	if pport > 0 {
		go func() {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", pport), nil); err != nil {
				logrus.Fatalf("error %v", err)
			}
		}()
	}

	http.HandleFunc("/say", sayHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/ui/index.html", uiStatusHandler)
	addr := ":" + strconv.Itoa(port)
	server := &http.Server{Addr: addr}
	server.SetKeepAlivesEnabled(false)
	log.Fatal(server.ListenAndServe())
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "my own error message", http.StatusInternalServerError)
}

func sayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, say I love you! port %d, pport %d", port, pport)
	if print {
		fmt.Fprintf(os.Stdout, "Hi there, say I love you! port %d, pport %d\n", port, pport)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "Hi there, I love you!")
}

func uiStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
