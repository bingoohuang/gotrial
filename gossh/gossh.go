package gossh

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// CreateClient with user(optional), host, port and password(optional)
// https://stackoverflow.com/questions/35450430/how-to-increase-golang-org-x-crypto-ssh-verbosity
// As a quick hack you can open $GOPATH/golang.org/x/crypto/ssh/mux.go file,
// change const debugMux = false to const debugMux = true and recompile your program.
func CreateClient(user, host string, port int, password string) (*ssh.Client, error) {
	auth, err := createAuth(password)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		Timeout:         3 * time.Second,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createAuth(password string) (ssh.AuthMethod, error) {
	var err error
	if password != "" {
		return ssh.Password(password), nil
	}

	auth, err := CreatePublickKey()
	if err != nil {
		return nil, err
	}
	if auth != nil {
		return auth, nil
	}

	return nil, errors.New("Please use password or auto sshed by ~/.ssh/id_rsa")
}

// CreatePublickKey from ~/.ssh/id_rs
func CreatePublickKey() (ssh.AuthMethod, error) {
	file, _ := homedir.Expand("~/.ssh/id_rsa")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, nil
	}

	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	// Use the PublicKeys method for remote authentication
	return ssh.PublicKeys(signer), nil
}

// RunScript in shell session from client
func RunScript(client *ssh.Client, scriptLine string) (string, error) {
	// create session
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(scriptLine)
	return string(out), err
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
