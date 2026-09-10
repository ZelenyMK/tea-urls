package shortener

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net/url"

	"crypto/md5"
	"io"

	"github.com/ZelenyMK/tea-urls/internal/config"
	urlverifier "github.com/davidmytton/url-verifier"
)

var (
	alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
)

func ShortenURL(currentURL string) (string, error) {
	verifier := urlverifier.NewVerifier()
	ret, errV := verifier.Verify(currentURL)
	if errV != nil {
		fmt.Errorf("Error: %s", errV)
		return "", errors.New("urlverifier error")
	}

	if !ret.IsURL {
		return "", errors.New("Not a valid URL")
	}

	parsedURL, errP := url.Parse(currentURL)
	if errP != nil {
		fmt.Errorf("Error: %s", errP)
		return "", errors.New("Could not parse the URL")
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
	if config.Host == "localhost" {
		shortendURL = parsedURL.Scheme + "://" + config.Host + ":" + config.Port + "/" + string(str)
	} else {
		shortendURL = parsedURL.Scheme + "://" + config.Host + "/" + string(str)
	}

	return shortendURL, nil
}
