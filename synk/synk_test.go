package synk_test

import "testing"

func TestSync(t *testing.T) {
	var (
		msg  string
		done bool
	)

	go func() {
		msg = "hello, world"
		done = true
	}()

	for {
		if done {
			println("msg", msg)
			break
		}
		println("retry...")
	}
}
