package mutex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChNil(t *testing.T) {
	var ch chan string

	c1 := make(chan string)
	c2 := make(chan string)

	stopR := make(chan bool)
	stopW := make(chan bool)

	go func() {
		select {
		case <-ch:
			c1 <- "OK"
		case <-stopR:
			c1 <- "Reader Stopped"
		}
	}()

	go func() {
		select {
		case ch <- "write":
			c2 <- "OK"
		case <-stopW:
			c2 <- "Writer Stopped"
		}
	}()

	time.Sleep(100 * time.Millisecond)

	stopR <- true
	stopW <- true

	for i := 0; i < 2; i++ {
		// Await both of these values
		// simultaneously, printing each one as it arrives.
		select {
		case msg1 := <-c1:
			assert.Equal(t, "Reader Stopped", msg1)
		case msg2 := <-c2:
			assert.Equal(t, "Writer Stopped", msg2)
		}
	}
}

func TestChNotEmpty(t *testing.T) {
	ch := make(chan string)

	c1 := make(chan string)
	c2 := make(chan string)

	stopR := make(chan bool)
	stopW := make(chan bool)

	go func() {
		select {
		case v := <-ch:
			c1 <- v
		case <-stopR:
			c1 <- "Reader Stopped"
		}
	}()

	go func() {
		select {
		case ch <- "write":
			c2 <- "write OK"
		case <-stopW:
			c2 <- "Writer Stopped"
		}
	}()

	time.Sleep(100 * time.Millisecond)

	stopR <- true
	stopW <- true

	for i := 0; i < 2; i++ {
		// Await both of these values
		// simultaneously, printing each one as it arrives.
		select {
		case msg1 := <-c1:
			assert.Equal(t, "write", msg1)
		case msg2 := <-c2:
			assert.Equal(t, "write OK", msg2)
		}
	}
}
