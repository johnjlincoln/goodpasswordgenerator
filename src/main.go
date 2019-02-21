package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Password defines the structure of password data that this endpoint will serve
type Password struct {
	Password            string
	DictionaryWordCount int
}

// Config defines the structure of configuration files for this application
type Config struct {
	WordDictionaryPath string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadConfigurations(configJSONPath string) Config {
	configJSON, err := os.Open(configJSONPath)
	check(err)
	decoder := json.NewDecoder(configJSON)
	config := Config{}
	err = decoder.Decode(&config)
	check(err)
	return config
}

func readWordDictionary(path string) ([]string, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words, scanner.Err()
}

func getSecurePassword(words []string, chars []string) ([]string, int) {
	var securePassowrd []string
	dictionaryWordCount := len(words)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	for len(securePassowrd) < 4 {
		securePassowrd = append(securePassowrd, words[r.Intn(dictionaryWordCount)])
	}
	randomChar := chars[r.Intn(len(chars))]
	randomNum := strconv.Itoa(r.Intn(999))
	securePassowrd = append(securePassowrd, randomChar, randomNum)
	return securePassowrd, dictionaryWordCount
}

func getSecurePasswordOpts(words []string) {}

func getDictionaryWordCount(words []string) int {
	return len(words)
}

func main() {
	config := loadConfigurations("config/sample.conf.json")
	words, err := readWordDictionary(config.WordDictionaryPath)
	check(err)
	specialChars := []string{"*", "!", "@", "#", "$", "%", "~", "_", "?", "+"}

	http.HandleFunc("/get/password", func(w http.ResponseWriter, r *http.Request) {
		securePassowrd, dictionaryWordCount := getSecurePassword(words, specialChars)
		securePasswordString := strings.Join(securePassowrd, "")
		response := Password{Password: securePasswordString, DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/password")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/get/dictionary/length", func(w http.ResponseWriter, r *http.Request) {
		dictionaryWordCount := getDictionaryWordCount(words)
		response := Password{DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/dictionary/length")
		json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
