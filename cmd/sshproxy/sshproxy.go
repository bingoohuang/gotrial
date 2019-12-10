package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	sshlib "github.com/blacknon/go-sshlib"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// Proxy ssh server
	host1     = "192.168.27.3"
	port1     = "60022"
	user1     = "huangjinbing"
	password1 = "Df5941B81A85#"

	// Target ssh server
	host2     = "192.168.29.11"
	port2     = "22"
	user2     = "footstone"
	password2 = ""

	termlog = "./test_term.log"
)

// https://godoc.org/github.com/blacknon/go-sshlib
func main() {
	// ==========
	// proxy connect
	// ==========

	// Create proxy sshlib.Connect
	proxyCon := &sshlib.Connect{}

	// Create proxy ssh.AuthMethod
	proxyAuthMethod := sshlib.CreateAuthMethodPassword(password1)

	// Connect proxy server
	err := proxyCon.CreateClient(host1, port1, user1, []ssh.AuthMethod{proxyAuthMethod})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// ==========
	// target connect
	// ==========

	// Create target sshlib.Connect
	targetCon := &sshlib.Connect{
		ProxyDialer: proxyCon.Client,
	}

	// Create target ssh.AuthMethod
	targetAuthMethod := sshlib.CreateAuthMethodPassword(password2)

	// Connect target server
	err = targetCon.CreateClient(host2, port2, user2, []ssh.AuthMethod{targetAuthMethod})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set terminal log
	targetCon.SetLog(termlog, false)

	session, err := targetCon.CreateSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// disable echoing input/output speed = 14.4kbaud
	modes := ssh.TerminalModes{ssh.ECHO: 0, ssh.TTY_OP_ISPEED: 14400, ssh.TTY_OP_OSPEED: 14400}
	if err := session.RequestPty("vt100", 800, 400, modes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := session.Shell(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux([]string{"date"}, w, r)
}

func mux(cmd []string, w io.Writer, r io.Reader) {
	var buf [65 * 1024]byte
	firstCmd := true
	last := ""

	for {
		t, err := r.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}

		sbuf, lastTwo := parseBuf(t, buf[:])
		switch lastTwo {
		case "$ ", "# ":
			if firstCmd {
				a := GetLastLine(last + sbuf)
				fmt.Print(a)
			} else {
				fmt.Print(last)
				fmt.Print(sbuf)
			}
			last = ""

			if len(cmd) > 0 {
				w.Write([]byte(cmd[0] + "\n"))
				cmd = cmd[1:]
				firstCmd = false
			} else {
				return
			}
		default:
			last += sbuf
		}
	}

}

// GetLastLine gets the last line of s.
func GetLastLine(s string) string {
	pos := strings.LastIndex(s, "\n")
	if pos < 0 || pos == len(s)-1 {
		return s
	}

	return s[pos+1:]
}

func parseBuf(t int, buf []byte) (sbuf, lastTwo string) {
	if t > 0 {
		sbuf = string(buf[:t])
	}
	if len(sbuf) > 2 {
		lastTwo = sbuf[t-2:]
	}

	return
}

// Advanced Unicode normalization and filtering,
// see http://blog.golang.org/normalization and
// http://godoc.org/golang.org/x/text/unicode/norm for more
// details.
func stripCtlAndExtFromUnicode(str string) string {
	isOk := func(r rune) bool {
		return r < 32 || r >= 127
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	str, _, _ = transform.String(t, str)
	return str
}
