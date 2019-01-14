package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var (
	usr      string
	host     string
	port     int
	scripts  string
	password string
	prompt   string
)

func init() {
	var scriptsFile string

	flag.StringVar(&usr, "u", "", "user, default to current user")
	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&scripts, "s", "", "scripts string")
	flag.StringVar(&password, "P", "", "password")
	flag.StringVar(&prompt, "t", ">", "prompt tip")
	flag.StringVar(&scriptsFile, "f", "", "scripts file")
	flag.IntVar(&port, "p", 22, "port")

	flag.Parse()

	if usr == "" {
		current, _ := user.Current()
		usr = current.Username
	}

	if host == "" {
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
		scripts = string(s)
	}

	if scripts == "" {
		scripts = "uname -a"
	}
}

func main() {
	client, err := gossh.CreateClient(usr, host, port, password)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer client.Close()

	scriptLines := gossh.SplitScriptLines(scripts)
	for _, scriptLine := range scriptLines {
		fmt.Println(prompt, scriptLine)
		out, err := gossh.RunScript(client, scriptLine)
		if err != nil {
			log.Fatalf("RunScript %s error %v", scriptLine, err)
		}
		fmt.Println(prompt, out)
	}
}
