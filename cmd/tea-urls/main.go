package main

import (
	"fmt"

	"github.com/ZelenyMK/tea-urls/internal/shortener"
)

func main() {
	result, _ := shortener.ShortenURL("https://www.google.com/")
	fmt.Println(result)
}
