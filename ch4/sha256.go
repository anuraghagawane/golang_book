package main

import "crypto/sha256"
import "fmt"
import "math/bits"

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)

	count := 0
	for i := 0; i < len(c1); i++ {
		mismatch := c1[i] ^ c2[i];

		count += bits.OnesCount8(mismatch)
	}
	fmt.Printf("%d\n", count)
}
