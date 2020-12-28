package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
)

//Convenience method for retrieving the SHA512 Hash
func fetchSha512Hash() hash.Hash {
	return sha512.New()
}

//Convenience method for retrieving the SHA256 Hash
func fetchSha256Hash() hash.Hash {
	return sha256.New()
}

//Closure to return an initialized Hash.
//It would be nice to avvoid the overhead of initializing the Hash for every request,
//but I'm not sure that the Hash functions are thread-safe.
func fetchInitializedSha512Hash() func() hash.Hash {
	var sha512Hash hash.Hash = nil
	return func() hash.Hash {
		if sha512Hash == nil {
			sha512Hash = sha512.New()
		}
		return sha512Hash
	}
}

//Encode a given string using the provided Hash
func encode(passwd string, hasher hash.Hash) string {
	hasher.Write([]byte(passwd))

	hashedPasswd := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	fmt.Println("Hashed " + passwd + " to " + hashedPasswd)

	return hashedPasswd
}
