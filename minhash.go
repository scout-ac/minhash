// Package minhash implements the Minhash algorithm.
package minhash

import "errors"

// ErrSigSizeMismatch is returned if a comparison is attempted between
// Minhashes with signatures of different sizes.
var ErrSigSizeMismatch = errors.New("signature sizes do not match")

// Minhash contains a hash function and a signature. The signature is the
// collection of minimum hashes for the supplied data. The hash function is
// used to populate the signature.
type Minhash struct {
	sig      []uint64
	hashFunc HashFunc
}

// HashFunc describes any func that takes []byte and returns uin64.
type HashFunc func([]byte) (uint64, uint64)

// NewMinhash returns a new MinWise Hashing implementation
func NewMinhash(hashFunc HashFunc, size int) *Minhash {
	sig := make([]uint64, size)
	for i := range sig {
		// Bit-flip zeros to ones to make the max uin64.
		sig[i] = ^uint64(0)
	}
	return &Minhash{
		hashFunc: hashFunc,
		sig:      sig,
	}
}

// NewMinhashFromSigs returns a new Minhash implementation populated with a
// provided signature.
func NewMinhashFromSigs(hashFunc HashFunc, existingSig []uint64) *Minhash {
	sig := make([]uint64, len(existingSig))
	// Copy the provided signature into the new blank signature.
	copy(sig, existingSig)
	return &Minhash{
		hashFunc: hashFunc,
		sig:      sig,
	}
}

// PushAll adds all provided elements to the signature.
func (m *Minhash) PushAll(bb [][]byte) {
	for _, b := range bb {
		m.Push(b)
	}
}

// PushStrings converts strings to bytes and updates the signature with those
// bytes.
func (m *Minhash) PushStrings(ss []string) {
	for _, r := range ss {
		m.Push([]byte(r))
	}
}

// Push adds an element to the signature.
func (m *Minhash) Push(b []byte) {
	// Create two independent hashes of the input.
	v1, v2 := m.hashFunc(b)

	for i, hashVal := range m.sig {
		// Simulate a 'new' hash func for every item in the signature.
		newHashVal := v1 + uint64(i)*v2

		// Apply the 'minimum' test of Minhash :)
		if newHashVal < hashVal {
			m.sig[i] = newHashVal
		}
	}
}

// Merge combines the signatures of the second set, creating the signature of
// their union.
func (m *Minhash) Merge(other *Minhash) {
	for i, v := range other.sig {
		// Apply the 'minimum' test of Minhash :)
		if v < m.sig[i] {
			m.sig[i] = v
		}
	}
}

// Signature returns the signature.
func (m *Minhash) Signature() []uint64 {
	return m.sig
}

// Similarity computes an estimate for the similarity between two signatures,
// expressed as a float between 0.0 (0% match) and 1.0 (100% match).
func (m *Minhash) Similarity(other *Minhash) (float64, error) {
	if len(m.sig) != len(other.sig) {
		return float64(0), ErrSigSizeMismatch
	}
	intersect := 0
	for i := range m.sig {
		if m.sig[i] == other.sig[i] {
			intersect++
		}
	}
	return float64(intersect) / float64(len(m.sig)), nil
}

// =========================================================================

// // Cardinality estimates the cardinality of the set
// func (m *Minhash) Cardinality() int {
// 	// http://www.cohenwang.com/edith/Papers/tcest.pdf
// 	sum := 0.0
// 	for _, v := range m.sig {
// 		sum += -math.Log(float64(math.MaxUint64-v) / float64(math.MaxUint64))
// 	}
// 	return int(float64(len(m.sig)-1) / sum)
// }

// // SignatureBbit returns a b-bit reduction of the signature. This will result
// // in unused bits at the high-end of the words if b does not divide 64
// // evenly.
// func (m *Minhash) SignatureBbit(b uint) []uint64 {
// 	var sig []uint64 // full signature
// 	var w uint64     // current word
// 	bits := uint(64) // bits free in current word
//
// 	mask := uint64(1<<b) - 1
//
// 	for _, v := range m.sig {
// 		if bits >= b {
// 			w <<= b
// 			w |= v & mask
// 			bits -= b
// 		} else {
// 			sig = append(sig, w)
// 			w = 0
// 			bits = 64
// 		}
// 	}
//
// 	if bits != 64 {
// 		sig = append(sig, w)
// 	}
//
// 	return sig
// }

// // SimilarityBbit computes an estimate for the similarity between two b-bit
// // signatures
// func SimilarityBbit(sig1, sig2 []uint64, b uint) (float64, error) {
// 	if len(sig1) != len(sig2) {
// 		return float64(0), ErrSigSizeMismatch
// 	}
//
// 	intersect := 0
// 	count := 0
//
// 	mask := uint64(1<<b) - 1
//
// 	for i := range sig1 {
// 		w1 := sig1[i]
// 		w2 := sig2[i]
//
// 		bits := uint(64)
//
// 		for bits >= b {
// 			v1 := (w1 & mask)
// 			v2 := (w2 & mask)
//
// 			count++
// 			if v1 == v2 {
// 				intersect++
// 			}
//
// 			bits -= b
// 			w1 >>= b
// 			w2 >>= b
// 		}
// 	}
//
// 	return float64(intersect) / float64(count), nil
// }
