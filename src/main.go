package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, words)
		check(err)
	})

	http.HandleFunc("/get/password", func(w http.ResponseWriter, r *http.Request) {
		securePassowrd, dictionaryWordCount := getSecurePassword(words)
		_, err := fmt.Fprintln(w, securePassowrd, dictionaryWordCount)
		check(err)
	})

	http.HandleFunc("/get/dictionary/length", func(w http.ResponseWriter, r *http.Request) {
		dictionaryLength := getDictionaryWordCount(words)
		_, err := fmt.Fprintln(w, dictionaryLength)
		check(err)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
