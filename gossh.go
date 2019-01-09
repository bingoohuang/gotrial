package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
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
		log.Fatalf("host argument missed")
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
	client := CreateClient(usr, host, port, scripts, password)
	defer client.Close()

	scriptLines := SplitScriptLines(scripts)
	for _, scriptLine := range scriptLines {
		fmt.Println(prompt, scriptLine)
		fmt.Print(prompt + " ")
		RunScript(client, scriptLine)
	}
}

// CreateClient from
func CreateClient(usr, host string, port int, scripts, password string) *ssh.Client {
	var auth ssh.AuthMethod
	if password != "" {
		auth = ssh.Password(password)
	} else {
		auth = CreatePublickKey()
	}

	config := &ssh.ClientConfig{
		User:            usr,
		Timeout:         10 * time.Second,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	return client
}

// CreatePublickKey from ~/.ssh/id_rs
func CreatePublickKey() ssh.AuthMethod {
	file, _ := homedir.Expand("~/.ssh/id_rsa")

	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	// Use the PublicKeys method for remote authentication
	return ssh.PublicKeys(signer)
}

// RunScript in shell session from client
func RunScript(client *ssh.Client, scriptLine string) {
	// create session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("unable to NewSession: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run(scriptLine)
}

// SplitScriptLines scripts to lines, ignoring comments or blank lines and auto join lines end with \
func SplitScriptLines(scripts string) []string {
	scriptLines := make([]string, 0)
	lastLine := ""
	for _, line := range strings.Split(scripts, "\n") {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" || strings.HasPrefix(trimLine, "#") {
			continue
		}

		if strings.HasSuffix(line, "\\") {
			lastLine += line[:len(line)-1]
		} else if lastLine != "" {
			scriptLines = append(scriptLines, lastLine+line)
			lastLine = ""
		} else {
			scriptLines = append(scriptLines, trimLine)
		}
	}

	if lastLine != "" {
		scriptLines = append(scriptLines, lastLine)
	}

	return scriptLines
}
