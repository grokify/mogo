package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/grokify/mogo/crypto/aesutil/aesargon"
)

// Assume these functions exist from our previous implementation
// func AesEncryptOpenSSLCompat(plaintext, password string) (string, error)
// func AesDecryptOpenSSLCompat(ciphertext, password string) (string, error)

func main() {
	useBase58 := true
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Choose an action: [e]ncrypt or [d]ecrypt")
	action, _ := reader.ReadString('\n')
	action = action[:len(action)-1]

	// Read the secret text
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]

	// Read password securely
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	crypter := aesargon.Crypter{
		Password:  string(bytePassword),
		UseBase58: useBase58,
	}

	switch action {
	case "e":
		encrypted, err := crypter.Encrypt(text)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
		fmt.Println("Encrypted:", encrypted)
	case "d":
		decrypted, err := crypter.Decrypt(text)
		if err != nil {
			fmt.Println("Decryption error:", err)
			return
		}
		fmt.Println("Decrypted:", decrypted)
	default:
		fmt.Println("Invalid action. Use 'e' or 'd'.")
	}
}
