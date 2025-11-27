// Package main demonstrates how to use minhash.
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"scout.ac/go/minhash"
)

func main() {
	f, _ := os.Open("plato1.txt")
	scanner := bufio.NewScanner(f)

	var tokens1 []string
	for scanner.Scan() {
		line := scanner.Text()
		tokens1 = append(tokens1, line)
	}

	f, _ = os.Open("plato2.txt")
	scanner = bufio.NewScanner(f)
	var tokens2 []string
	for scanner.Scan() {
		line := scanner.Text()
		tokens2 = append(tokens2, line)
	}
	start := time.Now()
	perms, _ := minhash.NewPermutations(64)
	m1 := minhash.NewMinhash(perms)
	for _, t := range tokens1 {
		m1.Update([]byte(t))
	}
	m2 := minhash.NewMinhash(perms)
	for _, t := range tokens2 {
		m2.Update([]byte(t))
	}
	similarity, err := m2.Jaccard(m1)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Similar: %f and Took %s", similarity, elapsed)
}
