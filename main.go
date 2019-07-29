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
	WordDictionaryPath   string
	SpecialCharsListPath string
	PasswordWordCount    int
	UseSpecialChars      bool
	UseNumber            bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfigurations() Config {
	var configJSONPath string
	configJSONPath, isSet := os.LookupEnv("CONFIG_JSON_PATH")
	if !isSet {
		configJSONPath = "config/dev.conf.json"
	}
	configJSON, err := os.Open(configJSONPath)
	check(err)
	decoder := json.NewDecoder(configJSON)
	config := Config{}
	err = decoder.Decode(&config)
	check(err)
	return config
}

func readPwFile(path string) ([]string, error) {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	var elements []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		elements = append(elements, scanner.Text())
	}
	return elements, scanner.Err()
}

func getSecurePassword(words []string, chars []string, config Config) ([]string, int) {
	var securePassowrd []string
	dictionaryWordCount := len(words)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	for len(securePassowrd) < config.PasswordWordCount {
		securePassowrd = append(securePassowrd, words[r.Intn(dictionaryWordCount)])
	}
	if config.UseSpecialChars {
		randomChar := chars[r.Intn(len(chars))]
		securePassowrd = append(securePassowrd, randomChar)
	}
	if config.UseNumber {
		randomNum := strconv.Itoa(r.Intn(999))
		securePassowrd = append(securePassowrd, randomNum)
	}
	return securePassowrd, dictionaryWordCount
}

// TODO: implement fn to generate password based on options passed in from front end (post payload)
// func getSecurePasswordOpts(words []string) {}

func getDictionaryWordCount(words []string) int {
	return len(words)
}

func main() {
	config := readConfigurations()
	words, err := readPwFile(config.WordDictionaryPath)
	check(err)
	specialChars, err := readPwFile(config.SpecialCharsListPath)
	check(err)

	http.HandleFunc("/get/password", func(w http.ResponseWriter, r *http.Request) {
		securePassowrd, dictionaryWordCount := getSecurePassword(words, specialChars, config)
		securePasswordString := strings.Join(securePassowrd, "")
		response := Password{Password: securePasswordString, DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/password")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/get/dictionary/wordcount", func(w http.ResponseWriter, r *http.Request) {
		dictionaryWordCount := getDictionaryWordCount(words)
		response := Password{DictionaryWordCount: dictionaryWordCount}
		fmt.Println("hit /get/dictionary/wordcount")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
