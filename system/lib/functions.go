/*
###################################
# Name: ZeroTwoGo                 #
# Version: Beta                   #
# Developer: VihangaYT            #
# Library: waSocket               #
# Contact: xxxxxxxxx              #
###################################
*/
package lib

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func GenerateRandomString(n int) string {
	// Karakter yang mungkin digunakan dalam string acak
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Inisialisasi generator angka acak dengan seed waktu
	rand.Seed(time.Now().UnixNano())

	// Membangun string acak
	result := make([]byte, n)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func IsValidImageURL(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		// Gagal membuat permintaan ke URL
		return false
	}
	defer resp.Body.Close()

	// Periksa status kode
	if resp.StatusCode != http.StatusOK {
		// Tanggapan tidak sukses
		return false
	}

	// Periksa header Content-Type
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image/") {
		// URL adalah URL gambar yang valid
		return true
	}

	// Content-Type bukan gambar
	return false
}

func ReqGet(url string, data interface{}) (err error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error making the request:", err)
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}
	return nil

}
