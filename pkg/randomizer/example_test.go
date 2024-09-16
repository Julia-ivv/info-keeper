package randomizer

import (
	"bytes"
	"fmt"
	"strings"
)

func ExampleGenerateRandomBytes() {
	length := 5
	b, err := GenerateRandomBytes(length)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(bytes.Equal(b, make([]byte, length)))

	// Output:
	// false
}

func ExampleGenerateRandomString() {
	length := 5
	randString, err := GenerateRandomString(length)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(strings.EqualFold(randString, ""))

	// Output:
	// false
}
