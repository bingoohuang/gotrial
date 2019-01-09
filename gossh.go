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
	"time"
)

var (
	usr     string
	host    string
	port    int
	scripts string
)

func init() {
	flag.StringVar(&usr, "u", "", "user, default to current user")
	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&scripts, "s", "uname -a", "scripts")
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
}

func main() {
	file, _ := homedir.Expand("~/.ssh/id_rsa")
	fmt.Println("file", file)

	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	fmt.Println("key read")

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	fmt.Println("PrivateKey parsed")

	config := &ssh.ClientConfig{
		User:    usr,
		Timeout: 30 * time.Second,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	fmt.Println("Connected!")
	defer client.Close()

	// create session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("unable to NewSession: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run(scripts)
}
