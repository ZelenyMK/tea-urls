package shortener

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"

	"crypto/md5"
	"io"

	"github.com/ZelenyMK/tea-urls/internal/config"
	urlverifier "github.com/davidmytton/url-verifier"
)

var (
	alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
)

func ShortenURL(currentURL string) (string, string, error) {
	verifier := urlverifier.NewVerifier()
	ret, err := verifier.Verify(currentURL)
	if err != nil {
		fmt.Errorf("Error: %s", err)
		return "", "", errors.New("urlverifier error")
	}

	if !ret.IsURL {
		return "", "", errors.New("Not a valid URL")
	}

	h := md5.New()                                        // https://stackoverflow.com/questions/48307105/how-do-i-use-a-string-as-input-to-the-rand-seed-function-in-golang
	io.WriteString(h, currentURL)                         //
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil)) //

	str := make([]rune, 16)
	r := rand.New(rand.NewSource(int64(seed)))
	for char := range str {
		str[char] = alphabet[r.Intn(62)]
	}

	shortendURL := ""
	alias := string(str)
	if config.Host == "localhost" || config.Host == "127.0.0.1" {
		shortendURL = "http" + "://" + config.Host + ":" + config.Port + "/" + alias
	} else {
		shortendURL = "http" + "://" + config.Host + "/" + alias
	}

	return shortendURL, alias, nil
}
