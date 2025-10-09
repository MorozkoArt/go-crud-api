package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	
	secret := base64.StdEncoding.EncodeToString(bytes)
	fmt.Println("Generated JWT Secret:")
	fmt.Println(secret)
	
	fmt.Println("\nAdd to config.yaml:")
	fmt.Printf("auth:\n  jwt_secret: \"%s\"\n  token_expiry: 24h\n", secret)
}