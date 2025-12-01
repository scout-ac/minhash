package minhash

import (
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/zeebo/xxh3"
)

var sentenceOne = []string{"it", "was", "the", "best", "of", "times", "it", "was", "the", "worst", "of", "times"}
var sentenceTwo = []string{"twas", "the", "best", "of", "times", "twas", "the", "blurst", "of", "times"}

var expectedFromOneHashA = [][]uint64{
	[]uint64{248687575117317545, 2049998094026730544, 765690271520631217, 5166404819947994582, 1194772720553424969, 2646378093565509457, 642549248073078814, 137226539518554757},
	[]uint64{1393834836782072322, 2049998094026730544, 3368748946450498991, 6654035784550370743, 1194772720553424969, 2646378093565509457, 642549248073078814, 339547347080119394},
}

var expectedFromOneHashB = [][]uint64{
	[]uint64{260238615480509520, 1483783405500556756, 2171411139035205116, 3930872985540651228, 1952832785009616865, 6377962565580745700, 730155147010689028, 3211658695656942387},
	[]uint64{4287877928842720490, 4063838249034338438, 2171411139035205116, 884995736121811762, 1952832785009616865, 9498231353857803716, 730155147010689028, 3211658695656942387},
}

var expectedFromTwoHashes = [][]uint64{
	[]uint64{248687575117317545, 1024286465335649203, 2794534105003552036, 775288689001785297, 4218713743468478472, 41033537384652751, 51622472557045086, 311861088037554606},
	[]uint64{8766980652033501079, 1024286465335649203, 2794534105003552036, 775288689001785297, 1533108964028516867, 41033537384652751, 1499301746664882380, 46143406790770579},
}

func TestOneHashA(t *testing.T) {
	var hashFunc = func(b []byte) (uint64, uint64) {
		hasher := xxh3.Hash128(b)
		return hasher.Lo, hasher.Hi
	}
	sim := doTestAndGetSimilarity(hashFunc, expectedFromOneHashA, t)
	expected := 0.5
	if sim != expected {
		t.Fatalf("Similarity looks bad, expected %v, got %v", expected, sim)
	}
}

func TestOneHashB(t *testing.T) {
	var hashFunc = func(b []byte) (uint64, uint64) {
		hasher := murmur3.New128()
		hasher.Write(b)
		hi, lo := hasher.Sum128()
		return hi, lo
	}
	sim := doTestAndGetSimilarity(hashFunc, expectedFromOneHashB, t)
	expected := 0.5
	if sim != expected {
		t.Fatalf("Similarity looks bad, expected %v, got %v", expected, sim)
	}
}

func TestTwoHashes(t *testing.T) {
	var hashFunc = func(b []byte) (uint64, uint64) {
		hashOne := xxh3.Hash(b)
		h := murmur3.New64()
		h.Write(b)
		hashTwo := h.Sum64()
		return hashOne, hashTwo
	}
	sim := doTestAndGetSimilarity(hashFunc, expectedFromTwoHashes, t)
	expected := 0.5
	if sim != expected {
		t.Fatalf("Similarity looks bad, expected %v, got %v", expected, sim)
	}
}

func doTestAndGetSimilarity(hashFunc HashFunc, expected [][]uint64, t *testing.T) float64 {
	m1 := NewMinhash(hashFunc, len(expected[0]))
	m2 := NewMinhash(hashFunc, len(expected[1]))

	m1.PushStrings(sentenceOne)
	// // or...
	// for _, s := range test {
	// 	mh.Push([]byte(s))
	// }

	m2.PushStrings(sentenceTwo)

	// t.Log("Comparing...")
	// t.Logf("    %v", m1.Signature())
	// t.Logf("    %v", expected[0])
	for i, sigVal := range m1.sig {
		expectedVal := expected[0][i]
		if sigVal != expectedVal {
			t.Log(m1.sig)
			t.Log(expected[0])
			t.Fatalf("Sig mismatch, item %v, expected %v got %v", i, expectedVal, sigVal)
		}
	}

	// t.Log("Comparing...")
	// t.Logf("    %v", m2.Signature())
	// t.Logf("    %v", expected[1])
	for i, sigVal := range m2.sig {
		expectedVal := expected[1][i]
		if sigVal != expectedVal {
			t.Log(m2.sig)
			t.Log(expected[1])
			t.Fatalf("Sig mismatch, item %v, expected %v got %v", i, expectedVal, sigVal)
		}
	}

	result, err := m1.Similarity(m2)
	if err != nil {
		panic(err)
	}
	return result
}
