// Package minhash implements the Minhash algorithm.
package minhash

import (
	"fmt"
	"math/rand/v2"

	// xxh3 is well-suited as a hash implementation for Minhash.
	"github.com/zeebo/xxh3"
)

// Bitflip zeros to ones to get the max value.
const maxUint64 uint64 = ^uint64(0)

// http://en.wikipedia.org/wiki/Mersenne_prime
const mersennePrime uint64 = uint64((1 << 61) - 1)

// Permutation describes a permutation for minhash.
type Permutation struct {
	Hi uint64
	Lo uint64
}

// Permutations contains all the permutations for minhash.
type Permutations struct {
	Size   int
	Values []Permutation
}

// Minhash contains permutations and hash values.
type Minhash struct {
	Permutations *Permutations
	HashValues   []uint64
}

func random(lo uint64, hi uint64) (uint64, error) {
	if hi < lo {
		return uint64(0), fmt.Errorf("'%d' is not higher than '%d'", hi, lo)
	}
	diff := hi - lo
	// Check for maximum possible value, return full range.
	if diff == maxUint64 {
		return rand.Uint64(), nil
	}
	return rand.Uint64N(diff+1) + lo, nil
}

// NewPermutations returns new permutations of the given size.
func NewPermutations(size int) (*Permutations, error) {
	p := Permutations{}
	p.Size = size
	p.Values = make([]Permutation, size)
	for i := range p.Values {
		hi, err := random(uint64(1), mersennePrime)
		if err != nil {
			return &Permutations{}, err
		}
		lo, err := random(uint64(0), mersennePrime)
		if err != nil {
			return &Permutations{}, err
		}
		p.Values[i] = Permutation{Hi: hi, Lo: lo}
	}
	// fmt.Println(p)
	return &p, nil
}

// NewMinhash returns a new Minhash with the given permutations.
func NewMinhash(perms *Permutations) *Minhash {
	m := Minhash{}
	m.Permutations = perms
	m.initHashvalues()
	return &m
}

// Hashvalues returns the hash values.
func (m *Minhash) Hashvalues() []uint64 {
	return m.HashValues
}

// Update updates the hash values with the given bytes.
func (m *Minhash) Update(b []byte) {
	// Create a new hasher every time this func is called. This is necessary
	// because xxh3 is stateful and will accumulate previous data.
	// Alternatively, the implementation could use a single hasher on Minhash
	// and call m.Reset() here.
	hasher := xxh3.New()
	hasher.Write(b)
	val := uint64(hasher.Sum64())
	for i, hashVal := range m.HashValues {
		// Apply the linear hash function h(v) = (a*v + b) mod p
		// The .Hi and .Lo pairs simulate two independent hashes.
		newHashVal := (m.Permutations.Values[i].Hi*val + m.Permutations.Values[i].Lo) % mersennePrime
		// Select the minimum value and update the hash value.
		if newHashVal > 0 && newHashVal < hashVal {
			m.HashValues[i] = newHashVal
		}
	}
}

// Jaccard computes the Jaccard similarity with the provided minhash.
func (m *Minhash) Jaccard(other *Minhash) (float64, error) {
	if m.Permutations.Size != other.Permutations.Size {
		return float64(0), fmt.Errorf("Size mismatch")
	}
	common := 0
	for i := range m.HashValues {
		if m.HashValues[i] == other.HashValues[i] {
			common++
		}
	}
	return float64(common) / float64(m.Permutations.Size), nil
}

func (m *Minhash) initHashvalues() {
	m.HashValues = make([]uint64, m.Permutations.Size)
	for i := range m.HashValues {
		m.HashValues[i] = maxUint64
	}
}
