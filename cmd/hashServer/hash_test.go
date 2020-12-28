package main

import (
	"crypto/sha512"
	"testing"
)

func TestEncode(t *testing.T) {
	hasher := sha512.New()

	encodedString := encode("password", hasher)
	if (len(encodedString)) == 0 {
		t.Fatal("Encoded string is null or empty")
	}
}
