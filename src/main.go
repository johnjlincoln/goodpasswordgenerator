package main

import (
    "fmt"
    // "html"
    "log"
    "net/http"
    "bufio"
    "os"
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

func main() {
    words, err := readWordDictionary("../slurp/1984-list.txt")
    check(err)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        _, err := fmt.Fprintln(w, words)
        check(err)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))

}
