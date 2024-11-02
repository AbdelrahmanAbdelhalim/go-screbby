package engine

import (
	"fmt"
	"testing"
)

func TestRandGeneration(t *testing.T) {
	prng := Init_prng(1234)
	n1 := prng.Rand64()
	n2 := prng.Rand64()
	n3 := prng.Rand64()

	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)

	if n1 == n2 || n1 == n3 || n3 == n2 {
		t.Errorf("Numbers generated are the same n1: %v, n2: %v, n3: %v", n1, n2, n3)
	}
}
