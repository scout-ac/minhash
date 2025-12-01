# Minhash

Forked from the original: <https://github.com/dgryski/go-minhash>

## Usage

You need to supply a hash function that takes `[]byte` as a single argument and returns **two** `uint64` values.
There are two ways you might want to do this.

**Option one:**

Inside your hash function, use two completely different hash functions.

```go
import (
	"github.com/spaolacci/murmur3"
	"github.com/zeebo/xxh3"
	"scout.ac/go/minhash"
)

// The number of elements in the signature.
// 128 is a reasonable size for most purposes.
const size = 128

func myHashFunc(b []byte) (uint64, uint64) {
	hashOne := xxh3.Hash(b)

	h := murmur3.New64()
	h.Reset()
	h.Write(b)
	hashTwo := h.Sum64()

	return hashOne, hashTwo
}

mh := minhash.NewMinhash(myHashFunc, size)
```

**Option two:**

Inside your hash function, use a single hash function but return two independent values.
A common approach is to take 128 bits, and split them into two 64-bit values.

```go
import (
	"github.com/zeebo/xxh3"
	"scout.ac/go/minhash"
)

const size = 128

func myHashFunc(b []byte) (uint64, uint64) {
	hash := xxh3.Hash128(b)
	return hash.Lo, hash.Hi
}

mh := minhash.NewMinhash(myHashFunc, size)
```

You can then add data to the `mh` hasher, and check the signature like so:

```go
mh := minhash.NewMinhash(myHashFunc, 4)
mh.PushStrings([]string{"some", "string", "data"})
fmt.Println(mh.Signature())
```

Output will be something like so:

```
[3747055036463229037 13280352286690612864 10737123926858574984 3467178332513026824]
```

Add byte data like so:

```go
mh.Push([]byte{"some string data"})

// or...
mh.PushAll([][]byte{
	[]byte("something with"),
	[]byte("several chunks"),
	[]byte("of data"),
})
```

Signatures can be compared like so:

```go
s1 := []string{"Excuse", "me", "while", "I", "kiss", "the", "sky"}
s2 := []string{"Excuse", "me", "while", "I", "kiss", "this", "guy"}

size := len(s1)

m1 := NewMinhash(hashFunc, size)
m2 := NewMinhash(hashFunc, size)

m1.PushStrings(s1)
m2.PushStrings(s2)

fmt.Println(m1.Similarity(m2))
```

The output should be `0.7142857142857143`, corresponding to 71.4% of the words (5 of the 7) in `s1` and `s2` being common to both.


## Maintainer's covenant

This project came about due to a personal need, and will be maintained as long as that need exists.
Consequently, bug fixes are important, but feature requests are unlikely to be accepted.
If you find a bug feel free to raise an issue, and preferably a pull-request with a solution.
Bugs will be prioritised according to whether they affect my ongoing requirements for this project.

Anyone is very welcome to fork this codebase if they wish.


## Thanks

Many thanks to [@dgryski](https://github.com/dgryski) for making their original code open-source and permissively licensed.
Maybe go [buy them a coffee](https://buymeacoffee.com/dgryski).
