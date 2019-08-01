package gossh

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

// CreateClient with user(optional), host, port and password(optional)
// https://stackoverflow.com/questions/35450430/how-to-increase-golang-org-x-crypto-ssh-verbosity
// As a quick hack you can open $GOPATH/golang.org/x/crypto/ssh/mux.go file,
// change const debugMux = false to const debugMux = true and recompile your program.
func CreateClient(host string, port int, user, password string) (*ssh.Client, error) {
	auth, err := createAuth(password)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		Timeout:         3 * time.Second,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := Dial("tcp", host+":"+strconv.Itoa(port), config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type Dialer func(ctx context.Context, net, addr string) (c net.Conn, err error)

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) Dialer {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, cTimeout)
		if err != nil {
			return conn, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}

// Dial starts a client connection to the given SSH server. It is a
// convenience function that connects to the given network address,
// initiates the SSH handshake, and then sets up a Client.  For access
// to incoming channels and requests, use net.Dial with NewClientConn
// instead.
func Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	dialer := TimeoutDialer(3*time.Second, 3*time.Second)
	conn, err := dialer(context.Background(), network, addr)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}
	return ssh.NewClient(c, chans, reqs), nil
}

func createAuth(password string) (ssh.AuthMethod, error) {
	var err error
	if password != "" {
		return ssh.Password(password), nil
	}

	auth, err := CreatePublicKey()
	if err != nil {
		return nil, err
	}
	if auth != nil {
		return auth, nil
	}

	return nil, errors.New("please use password or auto ssh by ~/.ssh/id_rsa")
}

// CreatePublicKey from ~/.ssh/id_rs
func CreatePublicKey() (ssh.AuthMethod, error) {
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

type singleWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

type AutoExitMode bool

const AutoExitOn AutoExitMode = true
const AutoExitOff AutoExitMode = false

func RunScripts(client *ssh.Client, scripts []string, autoExit AutoExitMode) (string, error) {
	if len(scripts) == 0 {
		return "", nil
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("vt100", 800, 400, modes); err != nil {
		return "", err
	}
	w, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	var b singleWriter
	session.Stdout = &b
	session.Stderr = &b

	if err := session.Shell(); err != nil {
		return b.b.String(), err
	}

	for _, cmd := range scripts {
		_, _ = w.Write([]byte(cmd + "\n"))
	}
	if autoExit == AutoExitOn {
		_, _ = w.Write([]byte("exit\n"))
	}

	if err := session.Wait(); err != nil {
		return b.b.String(), err
	}

	return b.b.String(), err
}

// SplitScriptLines scripts to lines, ignoring comments or blank lines and auto join lines end with "\"
func SplitScriptLines(scripts string) []string {
	scriptLines := make([]string, 0)
	lastLine := ""
	for _, line := range strings.Split(scripts, "\n") {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" || strings.HasPrefix(trimLine, "#") {
			continue
		}

		switch {
		case strings.HasSuffix(line, "\\"):
			lastLine += line[:len(line)-1]
		case lastLine != "":
			scriptLines = append(scriptLines, lastLine+line)
			lastLine = ""
		default:
			scriptLines = append(scriptLines, trimLine)
		}
	}

	if lastLine != "" {
		scriptLines = append(scriptLines, lastLine)
	}

	return scriptLines
}
