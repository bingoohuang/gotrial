package inline

import "testing"

func BenchmarkOriginGift(b *testing.B) {
	var nut = &OriginGift{
		dryFruit: &Chestnut{name: "栗子"},
	}

	for i := 0; i < b.N; i++ {
		nut.Access()
	}
}
func BenchmarkImprovedGift(b *testing.B) {
	var nut = &ImprovedGift{
		dryFruit: &Chestnut{name: "栗子"},
	}

	for i := 0; i < b.N; i++ {
		nut.Access()
	}
}
func BenchmarkOriginGiftParallel(b *testing.B) {
	var nut = &OriginGift{
		dryFruit: &Chestnut{name: "栗子"},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nut.mu.Lock()
			nut.Access()
			nut.mu.Unlock()
		}
	})
}
func BenchmarkImprovedGiftParallel(b *testing.B) {
	var nut = &ImprovedGift{
		dryFruit: &Chestnut{name: "栗子"},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nut.mu.Lock()
			nut.Access()
			nut.mu.Unlock()
		}
	})
}
