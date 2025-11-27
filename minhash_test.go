package minhash

import (
	"testing"
)

var s = []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetuer", "adipiscing", "elit", "nulla", "posuere"}

func TestNewPermutations(t *testing.T) {
	size := 64
	p, err := NewPermutations(size)
	if err != nil {
		t.Fatal(err)
	}
	if p.Size != size {
		t.Fatalf("Size expected %d but got %d", size, p.Size)
	}
	for _, value := range p.Values {
		if value.Hi < 1 || value.Hi > mersennePrime {
			t.Fatalf("Random %d out of bounds", value.Hi)
		}
		if value.Lo < 1 || value.Lo > mersennePrime {
			t.Fatalf("Random %d out of bounds", value.Lo)
		}
	}
}

func TestNewMinhash(t *testing.T) {
	perms, err := NewPermutations(64)
	if err != nil {
		t.Fatal(err)
	}
	m := NewMinhash(perms)
	if len(m.HashValues) != perms.Size {
		t.Fatalf("Hashvalues expected size %d but got %d", 64, len(m.HashValues))
	}
	if len(m.Permutations.Values) != perms.Size {
		t.Fatalf("Permutations expected size %d but got %d", 64, len(m.Permutations.Values))
	}

	for _, value := range m.HashValues {
		if value != maxUint64 {
			t.Fatalf("Expected infinity but got %d", value)
		}
	}

}

func TestUpdate(t *testing.T) {
	perms, err := NewPermutations(64)
	if err != nil {
		t.Fatal(err)
	}
	m := NewMinhash(perms)

	s := "lorem"
	m.Update([]byte(s))
}

func TestJaccardSame(t *testing.T) {
	perms, err := NewPermutations(64)
	if err != nil {
		t.Fatal(err)
	}
	m1 := NewMinhash(perms)
	m2 := NewMinhash(perms)

	for _, word := range s {
		m1.Update([]byte(word))
	}
	for _, word := range s {
		m2.Update([]byte(word))
	}

	ans, err := m1.Jaccard(m2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if ans != 1 {
		t.Fatalf("We should get similarity of 1")
	}
}

func TestJaccardDifferent(t *testing.T) {
	perms, err := NewPermutations(64)
	if err != nil {
		t.Fatal(err)
	}
	m1 := NewMinhash(perms)
	m2 := NewMinhash(perms)
	s1 := []string{"a", "b", "c", "d"}
	for _, s := range s1 {
		m1.Update([]byte(s))
	}
	s2 := []string{"e", "f", "g", "f"}
	for _, s := range s2 {
		m2.Update([]byte(s))
	}

	ans, err := m1.Jaccard(m2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if ans != 0 {
		t.Fatalf("We should get similarity of 0 but got %f", ans)
	}
}

func TestJaccardHalfEqual(t *testing.T) {
	perms, err := NewPermutations(60) // FIXME Should this be 64?
	if err != nil {
		t.Fatal(err)
	}
	m1 := NewMinhash(perms)
	m2 := NewMinhash(perms)
	s1 := []string{"a", "b", "c", "d"}
	for _, s := range s1 {
		m1.Update([]byte(s))
	}
	s2 := []string{"e", "f", "a", "b"}
	for _, s := range s2 {
		m2.Update([]byte(s))
	}

	ans, err := m1.Jaccard(m2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if ans <= 0.3 {
		t.Fatalf("We should get similarity of at least 0.3 but got %f", ans)
	}

}

func TestRandom(t *testing.T) {
	value, err := random(uint64(0), uint64(10))
	if err != nil {
		t.Fatal(err)
	}
	if value > 10 {
		t.Fatal("Expected less that 10")
	}
	value, err = random(uint64(0), uint64(0))
	if err != nil {
		t.Fatal(err)
	}
	if value != 0 {
		t.Fatal("Expected zero")
	}
	value, err = random(uint64(100), uint64(100))
	if err != nil {
		t.Fatal(err)
	}
	if value != 100 {
		t.Fatal("Expected hundred")
	}
}

func BenchmarkNew(b *testing.B) {
	perms, err := NewPermutations(64)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMinhash(perms)
	}
}

func BenchmarkUpdate(b *testing.B) {
	perms, err := NewPermutations(64)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m1 := NewMinhash(perms)
		for _, word := range s {
			m1.Update([]byte(word))
		}
	}
}
