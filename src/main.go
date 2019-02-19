package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// PasswordJSON is a json object containing a password string as well as the word count of the dictionary used
// https://golangcode.com/json-encode-an-array-of-objects/ try this cleaner implementation
type PasswordJSON struct {
	Password            string `json:"password"`
	DictionaryWordCount int    `json:"dictionaryWordCount"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func getSecurePassword(words []string) ([]string, int) {
	var securePassowrd []string
	dictionaryWordCount := len(words)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	for len(securePassowrd) < 4 {
		securePassowrd = append(securePassowrd, words[r.Intn(dictionaryWordCount)])
	}
	return securePassowrd, dictionaryWordCount
}

func getSecurePasswordOpts(words []string) {}

func getDictionaryWordCount(words []string) int {
	return len(words)
}

func main() {
	words, err := readWordDictionary("../slurp/1984-list.txt")
	check(err)
	var specialChars = make([]string, 10)
	specialChars[0] = "*"
	specialChars[1] = "!"
	specialChars[2] = "@"
	specialChars[3] = "#"
	specialChars[4] = "$"
	specialChars[5] = "%"
	specialChars[6] = "~"
	specialChars[7] = "_"
	specialChars[8] = "?"
	specialChars[9] = "+"

	http.HandleFunc("/get/password", func(w http.ResponseWriter, r *http.Request) {
		securePassowrd, dictionaryWordCount := getSecurePassword(words)
		securePasswordString := strings.Join(securePassowrd, "")
		response := PasswordJSON{Password: securePasswordString, DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/password")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/get/dictionary/length", func(w http.ResponseWriter, r *http.Request) {
		dictionaryWordCount := getDictionaryWordCount(words)
		response := PasswordJSON{DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/dictionary/length")
		json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
