package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/bingoohuang/golang-trial/gossh"
	"github.com/mitchellh/go-homedir"
)

type App struct {
	usr      string
	host     string
	port     int
	scripts  string
	password string
	prompt   string
}

func createApp() App {
	var scriptsFile string
	var app App

	flag.StringVar(&app.usr, "u", "", "user, default to current user")
	flag.StringVar(&app.host, "h", "", "host")
	flag.StringVar(&app.scripts, "s", "", "scripts string")
	flag.StringVar(&app.password, "P", "", "password")
	flag.StringVar(&app.prompt, "t", ">", "prompt tip")
	flag.StringVar(&scriptsFile, "f", "", "scripts file")
	flag.IntVar(&app.port, "p", 22, "port")

	flag.Parse()

	if app.usr == "" {
		current, _ := user.Current()
		app.usr = current.Username
	}

	if app.host == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "host argument missed\n")
		os.Exit(1)
	}

	if scriptsFile != "" {
		f, _ := homedir.Expand(scriptsFile)
		s, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("unable to read scripts file: %v", err)
		}
		app.scripts = string(s)
	}

	if app.scripts == "" {
		app.scripts = "uname -a"
	}

	return app
}

func main() {
	app := createApp()
	client, err := gossh.CreateClient(app.host, app.port, app.usr, app.password)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	fmt.Println("connected")
	start := time.Now()

	scriptLines := gossh.SplitScriptLines(app.scripts)
	out, err := gossh.RunScripts(client, scriptLines, gossh.AutoExitOn)
	fmt.Println("cost:", time.Since(start).String())
	fmt.Println(out)
	if err != nil {
		fmt.Println("error:", err)
	}

	//for _, scriptLine := range scriptLines {
	//	fmt.Println(prompt, scriptLine)
	//	out, err := gossh.RunScript(client, scriptLine)
	//	if err != nil {
	//		log.Fatalf("RunScript %s error %v", scriptLine, err)
	//	}
	//	fmt.Println(prompt, out)
	//}
}
