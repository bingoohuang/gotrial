package synk_test

import (
	"testing"

	. "github.com/bingoohuang/golang-trial/synk"
)

type one2 int

func (o *one2) Increment() {
	*o++
}

func run2(t *testing.T, once2 *Once2, o2 *one2, c chan bool) {
	once2.Do(func() { o2.Increment() })
	if v := *o2; v != 1 {
		t.Errorf("once2 failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce2(t *testing.T) {
	o2 := new(one2)
	once2 := new(Once2)
	c := make(chan bool)
	const N = 10

	for i := 0; i < N; i++ {
		go run2(t, once2, o2, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o2 != 1 {
		t.Errorf("once2 failed outside run: %d is not 1", *o2)
	}
}

func TestOncePanic2(t *testing.T) {
	var once2 Once2
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Once2.Do did not panic")
			}
		}()
		once2.Do(func() {
			panic("failed")
		})
	}()

	once2.Do(func() {
		t.Fatalf("Once2.Do called twice")
	})
}

func BenchmarkOnce2(b *testing.B) {
	var once2 Once2
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once2.Do(f)
		}
	})
}
